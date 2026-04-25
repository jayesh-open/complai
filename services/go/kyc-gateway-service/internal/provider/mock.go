package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/complai/complai/services/go/kyc-gateway-service/internal/domain"
)

var _ KYCProvider = (*MockProvider)(nil)

var (
	panRegex  = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]$`)
	gstinRegex = regexp.MustCompile(`^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z][A-Z0-9]{1}[Z][A-Z0-9]$`)
	tanRegex  = regexp.MustCompile(`^[A-Z]{4}[0-9]{5}[A-Z]$`)
	ifscRegex = regexp.MustCompile(`^[A-Z]{4}0[A-Z0-9]{6}$`)

	stateCodeMap = map[string]string{
		"01": "Jammu & Kashmir",
		"02": "Himachal Pradesh",
		"03": "Punjab",
		"04": "Chandigarh",
		"05": "Uttarakhand",
		"06": "Haryana",
		"07": "Delhi",
		"08": "Rajasthan",
		"09": "Uttar Pradesh",
		"10": "Bihar",
		"11": "Sikkim",
		"12": "Arunachal Pradesh",
		"13": "Nagaland",
		"14": "Manipur",
		"15": "Mizoram",
		"16": "Tripura",
		"17": "Meghalaya",
		"18": "Assam",
		"19": "West Bengal",
		"20": "Jharkhand",
		"21": "Odisha",
		"22": "Chhattisgarh",
		"23": "Madhya Pradesh",
		"24": "Gujarat",
		"25": "Daman & Diu",
		"26": "Dadra & Nagar Haveli",
		"27": "Maharashtra",
		"28": "Andhra Pradesh (Old)",
		"29": "Karnataka",
		"30": "Goa",
		"31": "Lakshadweep",
		"32": "Kerala",
		"33": "Tamil Nadu",
		"34": "Puducherry",
		"35": "Andaman & Nicobar Islands",
		"36": "Telangana",
		"37": "Andhra Pradesh",
	}

	panCategoryMap = map[byte]string{
		'A': "AOP (Association of Persons)",
		'B': "BOI (Body of Individuals)",
		'C': "Company",
		'F': "Firm",
		'G': "Government",
		'H': "HUF (Hindu Undivided Family)",
		'J': "Artificial Juridical Person",
		'L': "Local Authority",
		'P': "Individual",
		'T': "Trust",
	}

	bankNameMap = map[string]string{
		"SBIN": "State Bank of India",
		"HDFC": "HDFC Bank",
		"ICIC": "ICICI Bank",
		"UTIB": "Axis Bank",
		"KKBK": "Kotak Mahindra Bank",
		"PUNB": "Punjab National Bank",
		"BARB": "Bank of Baroda",
		"CNRB": "Canara Bank",
		"UBIN": "Union Bank of India",
		"IOBA": "Indian Overseas Bank",
		"IDIB": "Indian Bank",
		"BKID": "Bank of India",
		"CBIN": "Central Bank of India",
		"YESB": "Yes Bank",
		"INDB": "IndusInd Bank",
	}
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (m *MockProvider) VerifyPAN(_ context.Context, req *domain.PANVerifyRequest) (*domain.PANVerifyResponse, error) {
	pan := strings.ToUpper(strings.TrimSpace(req.PAN))

	if !panRegex.MatchString(pan) {
		return &domain.PANVerifyResponse{
			PAN:       req.PAN,
			Name:      "",
			Category:  "",
			Status:    "invalid",
			Valid:     false,
			RequestID: req.RequestID,
		}, nil
	}

	category := "Unknown"
	if cat, ok := panCategoryMap[pan[3]]; ok {
		category = cat
	}

	name := req.Name
	if name == "" {
		name = fmt.Sprintf("Mock Entity %s", pan[:5])
	}

	return &domain.PANVerifyResponse{
		PAN:       pan,
		Name:      name,
		Category:  category,
		Status:    "valid",
		Valid:     true,
		RequestID: req.RequestID,
	}, nil
}

func (m *MockProvider) VerifyGSTIN(_ context.Context, req *domain.GSTINVerifyRequest) (*domain.GSTINVerifyResponse, error) {
	gstin := strings.ToUpper(strings.TrimSpace(req.GSTIN))

	if !gstinRegex.MatchString(gstin) {
		return &domain.GSTINVerifyResponse{
			GSTIN:     req.GSTIN,
			Status:    "invalid",
			Valid:     false,
			RequestID: req.RequestID,
		}, nil
	}

	stateCode := gstin[:2]
	stateName, validState := stateCodeMap[stateCode]
	if !validState {
		return &domain.GSTINVerifyResponse{
			GSTIN:     gstin,
			Status:    "invalid",
			Valid:     false,
			RequestID: req.RequestID,
		}, nil
	}

	pan := gstin[2:12]

	return &domain.GSTINVerifyResponse{
		GSTIN:            gstin,
		LegalName:        fmt.Sprintf("Mock Legal Entity %s", pan[:5]),
		TradeName:        fmt.Sprintf("Mock Trade Name %s", pan[:5]),
		Status:           "Active",
		RegistrationType: "Regular",
		StateCode:        stateCode,
		State:            stateName,
		PAN:              pan,
		Valid:            true,
		RequestID:        req.RequestID,
	}, nil
}

func (m *MockProvider) VerifyTAN(_ context.Context, req *domain.TANVerifyRequest) (*domain.TANVerifyResponse, error) {
	tan := strings.ToUpper(strings.TrimSpace(req.TAN))

	if !tanRegex.MatchString(tan) {
		return &domain.TANVerifyResponse{
			TAN:       req.TAN,
			Name:      "",
			Status:    "invalid",
			Valid:     false,
			RequestID: req.RequestID,
		}, nil
	}

	return &domain.TANVerifyResponse{
		TAN:       tan,
		Name:      fmt.Sprintf("Mock Deductor %s", tan[:4]),
		Status:    "valid",
		Valid:     true,
		RequestID: req.RequestID,
	}, nil
}

func (m *MockProvider) VerifyBank(_ context.Context, req *domain.BankVerifyRequest) (*domain.BankVerifyResponse, error) {
	ifsc := strings.ToUpper(strings.TrimSpace(req.IFSC))
	accountNumber := strings.TrimSpace(req.AccountNumber)

	if !ifscRegex.MatchString(ifsc) {
		return &domain.BankVerifyResponse{
			AccountNumber: req.AccountNumber,
			IFSC:          req.IFSC,
			BankName:      "",
			BranchName:    "",
			NameAtBank:    "",
			Valid:         false,
			RequestID:     req.RequestID,
		}, nil
	}

	if accountNumber == "" {
		return &domain.BankVerifyResponse{
			AccountNumber: req.AccountNumber,
			IFSC:          ifsc,
			BankName:      "",
			BranchName:    "",
			NameAtBank:    "",
			Valid:         false,
			RequestID:     req.RequestID,
		}, nil
	}

	bankPrefix := ifsc[:4]
	bankName := "Unknown Bank"
	if name, ok := bankNameMap[bankPrefix]; ok {
		bankName = name
	}

	branchCode := ifsc[5:]
	branchName := fmt.Sprintf("%s Branch %s", bankName, branchCode)

	return &domain.BankVerifyResponse{
		AccountNumber: accountNumber,
		IFSC:          ifsc,
		BankName:      bankName,
		BranchName:    branchName,
		NameAtBank:    fmt.Sprintf("Account Holder %s", accountNumber[len(accountNumber)-4:]),
		Valid:         true,
		RequestID:     req.RequestID,
	}, nil
}
