package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/services/go/kyc-gateway-service/internal/domain"
)

// ---------------------------------------------------------------------------
// Tests: VerifyPAN
// ---------------------------------------------------------------------------

func TestMockProvider_VerifyPAN_ValidCompany(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), &domain.PANVerifyRequest{
		PAN:       "AABCA1234A",
		Name:      "Test Company Pvt Ltd",
		RequestID: "req-pan-1",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "valid", resp.Status)
	assert.Equal(t, "Company", resp.Category)
	assert.Equal(t, "AABCA1234A", resp.PAN)
	assert.Equal(t, "Test Company Pvt Ltd", resp.Name)
	assert.Equal(t, "req-pan-1", resp.RequestID)
}

func TestMockProvider_VerifyPAN_ValidIndividual(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), &domain.PANVerifyRequest{
		PAN:       "ABCPD1234E",
		Name:      "Deepak Kumar",
		RequestID: "req-pan-2",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "valid", resp.Status)
	assert.Equal(t, "Individual", resp.Category)
}

func TestMockProvider_VerifyPAN_ValidHUF(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), &domain.PANVerifyRequest{
		PAN:       "ABCHD1234E",
		Name:      "Kumar HUF",
		RequestID: "req-pan-3",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "HUF (Hindu Undivided Family)", resp.Category)
}

func TestMockProvider_VerifyPAN_Invalid(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), &domain.PANVerifyRequest{
		PAN:       "123",
		Name:      "Bad",
		RequestID: "req-pan-4",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid", resp.Status)
	assert.Empty(t, resp.Category)
}

func TestMockProvider_VerifyPAN_EmptyName(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), &domain.PANVerifyRequest{
		PAN:       "AABCA1234A",
		Name:      "",
		RequestID: "req-pan-5",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "Mock Entity AABCA", resp.Name)
}

func TestMockProvider_VerifyPAN_LowercaseNormalized(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyPAN(context.Background(), &domain.PANVerifyRequest{
		PAN:       "aabca1234a",
		Name:      "Test",
		RequestID: "req-pan-6",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "AABCA1234A", resp.PAN)
}

// ---------------------------------------------------------------------------
// Tests: VerifyGSTIN
// ---------------------------------------------------------------------------

func TestMockProvider_VerifyGSTIN_ValidKarnataka(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyGSTIN(context.Background(), &domain.GSTINVerifyRequest{
		GSTIN:     "29AABCA1234A1Z5",
		RequestID: "req-gstin-1",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "Active", resp.Status)
	assert.Equal(t, "Karnataka", resp.State)
	assert.Equal(t, "29", resp.StateCode)
	assert.Equal(t, "AABCA1234A", resp.PAN)
	assert.Equal(t, "Regular", resp.RegistrationType)
	assert.NotEmpty(t, resp.LegalName)
	assert.NotEmpty(t, resp.TradeName)
}

func TestMockProvider_VerifyGSTIN_ValidMaharashtra(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyGSTIN(context.Background(), &domain.GSTINVerifyRequest{
		GSTIN:     "27AABCA1234A1Z5",
		RequestID: "req-gstin-2",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "Maharashtra", resp.State)
	assert.Equal(t, "27", resp.StateCode)
}

func TestMockProvider_VerifyGSTIN_Invalid(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyGSTIN(context.Background(), &domain.GSTINVerifyRequest{
		GSTIN:     "XX",
		RequestID: "req-gstin-3",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid", resp.Status)
}

func TestMockProvider_VerifyGSTIN_InvalidStateCode(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyGSTIN(context.Background(), &domain.GSTINVerifyRequest{
		GSTIN:     "99AABCA1234A1Z5",
		RequestID: "req-gstin-4",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid", resp.Status)
}

// ---------------------------------------------------------------------------
// Tests: VerifyTAN
// ---------------------------------------------------------------------------

func TestMockProvider_VerifyTAN_Valid(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyTAN(context.Background(), &domain.TANVerifyRequest{
		TAN:       "BLRA12345B",
		RequestID: "req-tan-1",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "valid", resp.Status)
	assert.Equal(t, "BLRA12345B", resp.TAN)
	assert.Contains(t, resp.Name, "BLRA")
}

func TestMockProvider_VerifyTAN_Invalid(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyTAN(context.Background(), &domain.TANVerifyRequest{
		TAN:       "123",
		RequestID: "req-tan-2",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Equal(t, "invalid", resp.Status)
}

func TestMockProvider_VerifyTAN_LowercaseNormalized(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyTAN(context.Background(), &domain.TANVerifyRequest{
		TAN:       "blra12345b",
		RequestID: "req-tan-3",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "BLRA12345B", resp.TAN)
}

// ---------------------------------------------------------------------------
// Tests: VerifyBank
// ---------------------------------------------------------------------------

func TestMockProvider_VerifyBank_ValidSBI(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyBank(context.Background(), &domain.BankVerifyRequest{
		AccountNumber: "1234567890",
		IFSC:          "SBIN0001234",
		RequestID:     "req-bank-1",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "State Bank of India", resp.BankName)
	assert.Contains(t, resp.BranchName, "State Bank of India")
	assert.Contains(t, resp.NameAtBank, "7890")
	assert.Equal(t, "1234567890", resp.AccountNumber)
	assert.Equal(t, "SBIN0001234", resp.IFSC)
}

func TestMockProvider_VerifyBank_ValidHDFC(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyBank(context.Background(), &domain.BankVerifyRequest{
		AccountNumber: "9876543210",
		IFSC:          "HDFC0004567",
		RequestID:     "req-bank-2",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "HDFC Bank", resp.BankName)
}

func TestMockProvider_VerifyBank_ValidICICI(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyBank(context.Background(), &domain.BankVerifyRequest{
		AccountNumber: "5555555555",
		IFSC:          "ICIC0007890",
		RequestID:     "req-bank-3",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "ICICI Bank", resp.BankName)
}

func TestMockProvider_VerifyBank_InvalidIFSC(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyBank(context.Background(), &domain.BankVerifyRequest{
		AccountNumber: "1234567890",
		IFSC:          "123",
		RequestID:     "req-bank-4",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
	assert.Empty(t, resp.BankName)
}

func TestMockProvider_VerifyBank_EmptyAccountNumber(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyBank(context.Background(), &domain.BankVerifyRequest{
		AccountNumber: "",
		IFSC:          "SBIN0001234",
		RequestID:     "req-bank-5",
	})
	require.NoError(t, err)
	assert.False(t, resp.Valid)
}

func TestMockProvider_VerifyBank_UnknownBankPrefix(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.VerifyBank(context.Background(), &domain.BankVerifyRequest{
		AccountNumber: "1234567890",
		IFSC:          "XYZQ0001234",
		RequestID:     "req-bank-6",
	})
	require.NoError(t, err)
	assert.True(t, resp.Valid)
	assert.Equal(t, "Unknown Bank", resp.BankName)
}
