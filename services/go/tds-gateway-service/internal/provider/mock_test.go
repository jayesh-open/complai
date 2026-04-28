package provider

import (
	"context"
	"testing"

	"github.com/complai/complai/services/go/tds-gateway-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyPAN_Valid(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), domain.PANVerifyRequest{PAN: "ABCPD1234E", Name: "Test"})
	require.NoError(t, err)
	assert.Equal(t, "VALID", resp.Status)
	assert.Equal(t, "ABCPD1234E", resp.PAN)
	assert.Equal(t, "INDIVIDUAL", resp.Category)
	assert.Equal(t, "Test", resp.Name)
}

func TestVerifyPAN_InvalidLength(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), domain.PANVerifyRequest{PAN: "ABC"})
	require.NoError(t, err)
	assert.Equal(t, "INVALID", resp.Status)
}

func TestVerifyPAN_Categories(t *testing.T) {
	cases := []struct {
		pan  string
		want string
	}{
		{"ABCPD1234E", "INDIVIDUAL"},
		{"ABCCD1234E", "COMPANY"},
		{"ABCHD1234E", "HUF"},
		{"ABCFD1234E", "FIRM"},
		{"ABCTD1234E", "TRUST"},
		{"ABCAD1234E", "AOP"},
		{"ABCLD1234E", "LOCAL_AUTHORITY"},
		{"ABCGD1234E", "GOVERNMENT"},
		{"ABCXD1234E", "INDIVIDUAL"},
	}
	p := NewMockProvider()
	for _, tc := range cases {
		t.Run(tc.pan, func(t *testing.T) {
			resp, err := p.VerifyPAN(context.Background(), domain.PANVerifyRequest{PAN: tc.pan})
			require.NoError(t, err)
			assert.Equal(t, tc.want, resp.Category)
		})
	}
}

func TestVerifyPAN_DefaultName(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), domain.PANVerifyRequest{PAN: "ABCPD1234E"})
	require.NoError(t, err)
	assert.Equal(t, "Mock Entity ABCPD", resp.Name)
}

func TestVerifyTAN_Valid(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyTAN(context.Background(), domain.TANVerifyRequest{TAN: "MUMA12345A"})
	require.NoError(t, err)
	assert.Equal(t, "ACTIVE", resp.Status)
	assert.Equal(t, "MUMA12345A", resp.TAN)
	assert.Contains(t, resp.Name, "Mock Deductor")
}

func TestVerifyTAN_InvalidLength(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyTAN(context.Background(), domain.TANVerifyRequest{TAN: "SHORT"})
	require.NoError(t, err)
	assert.Equal(t, "NOT_FOUND", resp.Status)
}

func TestGenerateChallan(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.GenerateChallan(context.Background(), domain.ChallanRequest{
		TAN: "MUMA12345A", Section: "194C", Amount: 10000,
		Surcharge: 100, Cess: 50, Interest: 0, Penalty: 0,
		AssessmentYear: "2026-27",
	})
	require.NoError(t, err)
	assert.Equal(t, "SUCCESS", resp.Status)
	assert.NotEmpty(t, resp.ChallanNumber)
	assert.Equal(t, "0001234", resp.BSRCode)
	assert.InDelta(t, 10150.0, resp.Amount, 0.01)
}

func TestFileForm26Q(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FileForm26Q(context.Background(), domain.Form26QRequest{
		TAN: "MUMA12345A", FinancialYear: "2025-26", Quarter: "Q1",
		Deductions: []domain.Deduction26Q{
			{DeducteePAN: "ABCPD1234E", DeducteeName: "Test", Section: "194C", Amount: 50000, TDSAmount: 1000},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "ACCEPTED", resp.Status)
	assert.NotEmpty(t, resp.TokenNumber)
	assert.NotEmpty(t, resp.AcknowledgementNumber)
	assert.Contains(t, resp.TokenNumber, "TKN26Q")
}

func TestFileForm24Q(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FileForm24Q(context.Background(), domain.Form24QRequest{
		TAN: "MUMA12345A", FinancialYear: "2025-26", Quarter: "Q1",
		Employees: []domain.Employee24Q{
			{PAN: "ABCPD1234E", Name: "Employee1", GrossSalary: 1200000, TDSDeducted: 50000},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "ACCEPTED", resp.Status)
	assert.Contains(t, resp.TokenNumber, "TKN24Q")
}

func TestConcurrentChallans(t *testing.T) {
	p := NewMockProvider()
	done := make(chan *domain.ChallanResponse, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			resp, _ := p.GenerateChallan(context.Background(), domain.ChallanRequest{
				TAN: "MUMA12345A", Section: "194C", Amount: float64(n * 1000),
			})
			done <- resp
		}(i)
	}
	seen := map[string]bool{}
	for i := 0; i < 10; i++ {
		resp := <-done
		assert.NotNil(t, resp)
		assert.False(t, seen[resp.ChallanNumber], "duplicate challan number")
		seen[resp.ChallanNumber] = true
	}
}
