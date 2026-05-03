package provider

import (
	"context"
	"testing"

	"github.com/complai/complai/services/go/itr-gateway-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockProvider_ImplementsInterface(t *testing.T) {
	var _ SandboxITRProvider = (*MockProvider)(nil)
}

func TestCheckPANAadhaarLink_Linked(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.CheckPANAadhaarLink(context.Background(), domain.PANAadhaarLinkRequest{PAN: "ABCDE1234F"})
	require.NoError(t, err)
	assert.True(t, resp.Linked)
	assert.Equal(t, "ABCDE1234F", resp.PAN)
	assert.NotEmpty(t, resp.LinkDate)
}

func TestCheckPANAadhaarLink_NotLinked(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.CheckPANAadhaarLink(context.Background(), domain.PANAadhaarLinkRequest{PAN: "ABCDE1234Z"})
	require.NoError(t, err)
	assert.False(t, resp.Linked)
}

func TestCheckPANAadhaarLink_InvalidPAN(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.CheckPANAadhaarLink(context.Background(), domain.PANAadhaarLinkRequest{PAN: "SHORT"})
	require.NoError(t, err)
	assert.False(t, resp.Linked)
}

func TestFetchAIS_Success(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.FetchAIS(context.Background(), domain.AISRequest{PAN: "ABCDE1234F", TaxYear: "2026-27"})
	require.NoError(t, err)
	assert.Equal(t, "ABCDE1234F", resp.PAN)
	assert.Equal(t, "2026-27", resp.TaxYear)
	assert.NotEmpty(t, resp.Form168Ref)
	assert.Len(t, resp.TDSEntries, 2)
	assert.Equal(t, "392", resp.TDSEntries[0].Section)
}

func TestFetchAIS_InvalidPAN(t *testing.T) {
	p := NewMockProvider()
	_, err := p.FetchAIS(context.Background(), domain.AISRequest{PAN: "SHORT", TaxYear: "2026-27"})
	assert.Error(t, err)
}

func TestSubmitITR_Success(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.SubmitITR(context.Background(), domain.ITRSubmitRequest{
		PAN: "ABCDE1234F", TaxYear: "2026-27", FormType: "ITR-1", Payload: "{}",
	})
	require.NoError(t, err)
	assert.Contains(t, resp.ARN, "ARN-2026-27")
	assert.Equal(t, "SUBMITTED", resp.Status)
	assert.NotEmpty(t, resp.AcknowledgementNo)
}

func TestSubmitITR_InvalidPAN(t *testing.T) {
	p := NewMockProvider()
	_, err := p.SubmitITR(context.Background(), domain.ITRSubmitRequest{
		PAN: "BAD", TaxYear: "2026-27", FormType: "ITR-1",
	})
	assert.Error(t, err)
}

func TestSubmitITR_UnsupportedForm(t *testing.T) {
	p := NewMockProvider()
	_, err := p.SubmitITR(context.Background(), domain.ITRSubmitRequest{
		PAN: "ABCDE1234F", TaxYear: "2026-27", FormType: "ITR-99",
	})
	assert.Error(t, err)
}

func TestGenerateITRV_Success(t *testing.T) {
	p := NewMockProvider()
	sub, _ := p.SubmitITR(context.Background(), domain.ITRSubmitRequest{
		PAN: "ABCDE1234F", TaxYear: "2026-27", FormType: "ITR-1",
	})

	resp, err := p.GenerateITRV(context.Background(), domain.ITRVRequest{ARN: sub.ARN})
	require.NoError(t, err)
	assert.Equal(t, sub.ARN, resp.ARN)
	assert.Contains(t, resp.ITRVURL, sub.ARN)
}

func TestGenerateITRV_NotFound(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GenerateITRV(context.Background(), domain.ITRVRequest{ARN: "BOGUS"})
	assert.Error(t, err)
}

func TestCheckEVerification_Success(t *testing.T) {
	p := NewMockProvider()
	sub, _ := p.SubmitITR(context.Background(), domain.ITRSubmitRequest{
		PAN: "ABCDE1234F", TaxYear: "2026-27", FormType: "ITR-2",
	})

	resp, err := p.CheckEVerification(context.Background(), domain.EVerifyRequest{ARN: sub.ARN})
	require.NoError(t, err)
	assert.True(t, resp.Verified)
	assert.Equal(t, "AADHAAR_OTP", resp.Method)
}

func TestCheckEVerification_NotFound(t *testing.T) {
	p := NewMockProvider()
	_, err := p.CheckEVerification(context.Background(), domain.EVerifyRequest{ARN: "BOGUS"})
	assert.Error(t, err)
}

func TestCheckRefundStatus_Success(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.CheckRefundStatus(context.Background(), domain.RefundStatusRequest{
		PAN: "ABCDE1234F", TaxYear: "2026-27",
	})
	require.NoError(t, err)
	assert.Equal(t, "PROCESSED", resp.Status)
	assert.Equal(t, float64(15000), resp.Amount)
}

func TestCheckRefundStatus_InvalidPAN(t *testing.T) {
	p := NewMockProvider()
	_, err := p.CheckRefundStatus(context.Background(), domain.RefundStatusRequest{
		PAN: "BAD", TaxYear: "2026-27",
	})
	assert.Error(t, err)
}
