package filing

import (
	"context"
	"fmt"
	"time"

	"github.com/complai/complai/services/go/tds-service/internal/domain"
	"github.com/complai/complai/services/go/tds-service/internal/store"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GatewayClient interface {
	SubmitForm(ctx context.Context, tenantID string, formType domain.FormType, fvuContent string) (*SubmitResult, error)
}

type SubmitResult struct {
	TokenNumber           string `json:"token_number"`
	AcknowledgementNumber string `json:"acknowledgement_number"`
	Status                string `json:"status"`
}

type FilingSagaInput struct {
	TenantID       uuid.UUID
	FormType       domain.FormType
	FinancialYear  string
	Quarter        string
	TAN            string
	Employer       *domain.EmployerDetails
	Deductor       *domain.DeductorDetails
	CountryCodes   map[string]string
	DTAAArticles   map[string]string
	DTAARates      map[string]decimal.Decimal
	CurrencyCodes  map[string]string
	ExchangeRates  map[string]decimal.Decimal
	ForeignAmounts map[string]decimal.Decimal
}

type FilingSaga struct {
	store   store.Repository
	gateway GatewayClient
}

func NewFilingSaga(s store.Repository, gw GatewayClient) *FilingSaga {
	return &FilingSaga{store: s, gateway: gw}
}

func (s *FilingSaga) Execute(ctx context.Context, input FilingSagaInput) (*domain.Filing, error) {
	idemKey := domain.FilingIdempotencyKey(input.TenantID, input.FormType, input.FinancialYear, input.Quarter)

	existing, err := s.store.GetFilingByIdempotencyKey(ctx, input.TenantID, idemKey)
	if err == nil && existing != nil {
		if existing.Status == domain.FilingFiled || existing.Status == domain.FilingSubmitted {
			return existing, nil
		}
	}

	deductees, _, err := s.store.ListDeductees(ctx, input.TenantID, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("validate: list deductees: %w", err)
	}
	if err := s.validateDeductees(deductees, input.FormType); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	entries, _, err := s.store.ListEntries(ctx, input.TenantID, input.FinancialYear, input.Quarter, 10000, 0)
	if err != nil {
		return nil, fmt.Errorf("aggregate: list entries: %w", err)
	}

	fvuContent, deducteeCount, totalTDS, err := s.generateFVU(input, deductees, entries)
	if err != nil {
		return nil, fmt.Errorf("generate FVU: %w", err)
	}

	var filing *domain.Filing
	if existing != nil && existing.Status == domain.FilingRejected {
		filing = existing
		filing.Status = domain.FilingDraft
		filing.FVUContent = fvuContent
		filing.DeducteeCount = deducteeCount
		filing.TotalTDSAmount = totalTDS.StringFixed(2)
		filing.ErrorMessage = ""
		filing.UpdatedAt = time.Now()
		s.store.UpdateFilingStatus(ctx, input.TenantID, filing.ID, domain.FilingDraft, "", "", "")
	} else {
		filing = &domain.Filing{
			ID:             uuid.New(),
			TenantID:       input.TenantID,
			FormType:       input.FormType,
			FinancialYear:  input.FinancialYear,
			Quarter:        input.Quarter,
			TAN:            input.TAN,
			Status:         domain.FilingDraft,
			DeducteeCount:  deducteeCount,
			TotalTDSAmount: totalTDS.StringFixed(2),
			FVUContent:     fvuContent,
			IdempotencyKey: idemKey,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.store.CreateFiling(ctx, input.TenantID, filing); err != nil {
			existingFiling, lookupErr := s.store.GetFilingByIdempotencyKey(ctx, input.TenantID, idemKey)
			if lookupErr == nil && existingFiling != nil {
				return existingFiling, nil
			}
			return nil, fmt.Errorf("persist filing: %w", err)
		}
	}

	result, err := s.gateway.SubmitForm(ctx, input.TenantID.String(), input.FormType, fvuContent)
	if err != nil {
		s.store.UpdateFilingStatus(ctx, input.TenantID, filing.ID, domain.FilingRejected, "", "", err.Error())
		return nil, fmt.Errorf("submit: %w", err)
	}

	filedStatus := domain.FilingFiled
	if result.Status == "SUBMITTED" {
		filedStatus = domain.FilingSubmitted
	}
	s.store.UpdateFilingStatus(ctx, input.TenantID, filing.ID, filedStatus, result.TokenNumber, result.AcknowledgementNumber, "")

	filing.Status = filedStatus
	filing.TokenNumber = result.TokenNumber
	filing.AcknowledgementNumber = result.AcknowledgementNumber

	return filing, nil
}

func (s *FilingSaga) validateDeductees(deductees []domain.Deductee, formType domain.FormType) error {
	if len(deductees) == 0 {
		return fmt.Errorf("no deductees found")
	}

	var missing int
	for _, d := range deductees {
		switch formType {
		case domain.FormType138:
			if d.PAN == "" {
				missing++
			}
		case domain.FormType140:
			if d.PAN == "" && d.ResidentStatus == domain.Resident {
				missing++
			}
		case domain.FormType144:
			if d.ResidentStatus != domain.NonResident {
				continue
			}
			if d.PAN == "" {
				missing++
			}
		}
	}

	if missing > 0 {
		return fmt.Errorf("%d deductee(s) missing PAN for Form %s filing", missing, formType)
	}
	return nil
}

func (s *FilingSaga) generateFVU(input FilingSagaInput, deductees []domain.Deductee, entries []domain.TDSEntry) (string, int, decimal.Decimal, error) {
	switch input.FormType {
	case domain.FormType138:
		if input.Employer == nil {
			return "", 0, decimal.Zero, fmt.Errorf("employer details required for Form 138")
		}
		payload, err := domain.GenerateForm138(domain.Form138Input{
			Employer:      *input.Employer,
			FinancialYear: input.FinancialYear,
			Quarter:       input.Quarter,
			Deductees:     deductees,
			Entries:       entries,
		})
		if err != nil {
			return "", 0, decimal.Zero, err
		}
		fvu := domain.GenerateForm138FVU(payload)
		return fvu, len(payload.Employees), payload.TotalTDS, nil

	case domain.FormType140:
		if input.Deductor == nil {
			return "", 0, decimal.Zero, fmt.Errorf("deductor details required for Form 140")
		}
		payload, err := domain.GenerateForm140(domain.Form140Input{
			Deductor:      *input.Deductor,
			FinancialYear: input.FinancialYear,
			Quarter:       input.Quarter,
			Deductees:     deductees,
			Entries:       entries,
		})
		if err != nil {
			return "", 0, decimal.Zero, err
		}
		fvu := domain.GenerateForm140FVU(payload)
		return fvu, len(payload.Deductions), payload.TotalTDS, nil

	case domain.FormType144:
		if input.Deductor == nil {
			return "", 0, decimal.Zero, fmt.Errorf("deductor details required for Form 144")
		}
		payload, err := domain.GenerateForm144(domain.Form144Input{
			Deductor:       *input.Deductor,
			FinancialYear:  input.FinancialYear,
			Quarter:        input.Quarter,
			Deductees:      deductees,
			Entries:        entries,
			CountryCodes:   input.CountryCodes,
			DTAAArticles:   input.DTAAArticles,
			DTAARates:      input.DTAARates,
			CurrencyCodes:  input.CurrencyCodes,
			ExchangeRates:  input.ExchangeRates,
			ForeignAmounts: input.ForeignAmounts,
		})
		if err != nil {
			return "", 0, decimal.Zero, err
		}
		fvu := domain.GenerateForm144FVU(payload)
		return fvu, len(payload.Remittances), payload.TotalTDS, nil

	default:
		return "", 0, decimal.Zero, fmt.Errorf("unsupported form type: %s", input.FormType)
	}
}
