package provider

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/services/go/gstn-gateway-service/internal/domain"
)

func TestMockProvider_Authenticate(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.Authenticate(context.Background())
	require.NoError(t, err)
	assert.Contains(t, resp.AccessToken, "mock-gsp-token-")
	assert.Equal(t, "bearer", resp.TokenType)
	assert.Equal(t, 86399, resp.ExpiresIn)
}

func TestMockProvider_GSTR1Save_NewFiling(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.GSTR1Save(context.Background(), &domain.GSTR1SaveRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Section: "b2b",
		Data: map[string]interface{}{"invoices": []interface{}{}}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.NotEmpty(t, resp.Token)
}

func TestMockProvider_GSTR1Save_Idempotent(t *testing.T) {
	p := NewMockProvider()
	reqID := uuid.New().String()
	req := &domain.GSTR1SaveRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Section: "b2b",
		Data: map[string]interface{}{}, RequestID: reqID,
	}

	resp1, err := p.GSTR1Save(context.Background(), req)
	require.NoError(t, err)
	resp2, err := p.GSTR1Save(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR1Get_Empty(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.GSTR1Get(context.Background(), &domain.GSTR1GetRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "draft", resp.Status)
	assert.Empty(t, resp.Data)
}

func TestMockProvider_GSTR1Get_WithSection(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	_, err := p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)

	_, err = p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Section: "hsn",
		Data: map[string]interface{}{"y": 2}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)

	resp, err := p.GSTR1Get(ctx, &domain.GSTR1GetRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Section: "b2b", RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Contains(t, resp.Data, "b2b")
	assert.NotContains(t, resp.Data, "hsn")

	respAll, err := p.GSTR1Get(ctx, &domain.GSTR1GetRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Contains(t, respAll.Data, "b2b")
	assert.Contains(t, respAll.Data, "hsn")
}

func TestMockProvider_GSTR1Reset(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	_, err := p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", Section: "b2b",
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)

	resp, err := p.GSTR1Reset(ctx, &domain.GSTR1ResetRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)

	getResp, err := p.GSTR1Get(ctx, &domain.GSTR1GetRequest{
		GSTIN: "29AABCA1234A1Z5", RetPeriod: "042026", RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "draft", getResp.Status)
	assert.Empty(t, getResp.Data)
}

func TestMockProvider_GSTR1Reset_NoFiling(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.GSTR1Reset(context.Background(), &domain.GSTR1ResetRequest{
		GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "Nothing to reset", resp.Message)
}

func TestMockProvider_GSTR1Reset_Idempotent(t *testing.T) {
	p := NewMockProvider()
	reqID := uuid.New().String()
	req := &domain.GSTR1ResetRequest{GSTIN: "X", RetPeriod: "042026", RequestID: reqID}

	resp1, err := p.GSTR1Reset(context.Background(), req)
	require.NoError(t, err)
	resp2, err := p.GSTR1Reset(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR1Reset_AfterFiled(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	gstin := "29AABCA1234A1Z5"

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: gstin, RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: gstin, RetPeriod: "042026", RequestID: uuid.New().String()})
	p.GSTR1File(ctx, &domain.GSTR1FileRequest{GSTIN: gstin, RetPeriod: "042026", SignType: "EVC", EVOTP: "123456", PAN: "AABCA1234A", RequestID: uuid.New().String()})

	_, err := p.GSTR1Reset(ctx, &domain.GSTR1ResetRequest{GSTIN: gstin, RetPeriod: "042026", RequestID: uuid.New().String()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already filed")
}

func TestMockProvider_GSTR1Submit_NoFiling(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR1Submit(context.Background(), &domain.GSTR1SubmitRequest{
		GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no draft found")
}

func TestMockProvider_GSTR1Submit_NoSections(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Reset(ctx, &domain.GSTR1ResetRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})

	_, err := p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no sections saved")
}

func TestMockProvider_GSTR1Submit_Idempotent(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String()})

	reqID := uuid.New().String()
	resp1, err := p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: reqID})
	require.NoError(t, err)
	resp2, err := p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: reqID})
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR1Submit_AlreadyFiled(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})
	p.GSTR1File(ctx, &domain.GSTR1FileRequest{GSTIN: "X", RetPeriod: "042026", SignType: "EVC", EVOTP: "123456", PAN: "X", RequestID: uuid.New().String()})

	_, err := p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already filed")
}

func TestMockProvider_GSTR1File_NoFiling(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR1File(context.Background(), &domain.GSTR1FileRequest{
		GSTIN: "X", RetPeriod: "042026", SignType: "EVC", EVOTP: "123456", PAN: "X", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no filing found")
}

func TestMockProvider_GSTR1File_NotSubmitted(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})

	_, err := p.GSTR1File(ctx, &domain.GSTR1FileRequest{
		GSTIN: "X", RetPeriod: "042026", SignType: "EVC", EVOTP: "123456", PAN: "X", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must submit before filing")
}

func TestMockProvider_GSTR1File_AlreadyFiled(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})
	p.GSTR1File(ctx, &domain.GSTR1FileRequest{GSTIN: "X", RetPeriod: "042026", SignType: "EVC", EVOTP: "123456", PAN: "X", RequestID: uuid.New().String()})

	_, err := p.GSTR1File(ctx, &domain.GSTR1FileRequest{
		GSTIN: "X", RetPeriod: "042026", SignType: "EVC", EVOTP: "123456", PAN: "X", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already filed")
}

func TestMockProvider_GSTR1File_InvalidSignType(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})

	_, err := p.GSTR1File(ctx, &domain.GSTR1FileRequest{
		GSTIN: "X", RetPeriod: "042026", SignType: "INVALID", PAN: "X", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sign_type")
}

func TestMockProvider_GSTR1File_EVCMissingOTP(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})

	_, err := p.GSTR1File(ctx, &domain.GSTR1FileRequest{
		GSTIN: "X", RetPeriod: "042026", SignType: "EVC", EVOTP: "", PAN: "X", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EVC OTP required")
}

func TestMockProvider_GSTR1File_DSCNoOTPNeeded(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})

	resp, err := p.GSTR1File(ctx, &domain.GSTR1FileRequest{
		GSTIN: "X", RetPeriod: "042026", SignType: "DSC", PAN: "X", RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.NotEmpty(t, resp.ARN)
}

func TestMockProvider_GSTR1File_Idempotent(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})

	reqID := uuid.New().String()
	req := &domain.GSTR1FileRequest{GSTIN: "X", RetPeriod: "042026", SignType: "EVC", EVOTP: "123456", PAN: "X", RequestID: reqID}

	resp1, err := p.GSTR1File(ctx, req)
	require.NoError(t, err)
	resp2, err := p.GSTR1File(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR1Status_NoFiling(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.GSTR1Status(context.Background(), &domain.GSTR1StatusRequest{
		GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "draft", resp.Status)
	assert.Empty(t, resp.ARN)
}

func TestMockProvider_GSTR1Status_Filed(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})
	p.GSTR1File(ctx, &domain.GSTR1FileRequest{GSTIN: "X", RetPeriod: "042026", SignType: "EVC", EVOTP: "123456", PAN: "X", RequestID: uuid.New().String()})

	resp, err := p.GSTR1Status(ctx, &domain.GSTR1StatusRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})
	require.NoError(t, err)
	assert.Equal(t, "filed", resp.Status)
	assert.NotEmpty(t, resp.ARN)
	assert.NotNil(t, resp.FiledAt)
}

func TestMockProvider_GSTR1Save_AfterSubmitted(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})

	_, err := p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "hsn", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot save")
}

func TestMockProvider_ResetState(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.ResetState()

	resp, err := p.GSTR1Get(ctx, &domain.GSTR1GetRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})
	require.NoError(t, err)
	assert.Equal(t, "draft", resp.Status)
	assert.Empty(t, resp.Data)
}

func TestMockProvider_GSTR1Submit_AlreadySubmitted(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	p.GSTR1Save(ctx, &domain.GSTR1SaveRequest{GSTIN: "X", RetPeriod: "042026", Section: "b2b", Data: map[string]interface{}{}, RequestID: uuid.New().String()})
	p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})

	_, err := p.GSTR1Submit(ctx, &domain.GSTR1SubmitRequest{GSTIN: "X", RetPeriod: "042026", RequestID: uuid.New().String()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already submitted")
}

// ---------------------------------------------------------------------------
// GSTR-9 Provider Tests
// ---------------------------------------------------------------------------

const testGSTIN = "27AABCU9603R1ZM"
const testFY = "2025-26"

func TestMockProvider_GSTR9Save_Success(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.GSTR9Save(context.Background(), &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"tables": "data"}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Contains(t, resp.Reference, "GSTR9-")
	assert.NotEmpty(t, resp.SavedAt)
}

func TestMockProvider_GSTR9Save_Idempotent(t *testing.T) {
	p := NewMockProvider()
	reqID := uuid.New().String()
	req := &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: reqID,
	}
	resp1, err := p.GSTR9Save(context.Background(), req)
	require.NoError(t, err)
	resp2, err := p.GSTR9Save(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR9Save_InvalidGSTIN(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9Save(context.Background(), &domain.GSTR9SaveRequest{
		GSTIN: "SHORT", FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid GSTIN")
}

func TestMockProvider_GSTR9Save_InvalidFY(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9Save(context.Background(), &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: "2025",
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid financial_year")
}

func TestMockProvider_GSTR9Save_UpdateExisting(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	resp1, err := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"v": 1}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)

	resp2, err := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"v": 2}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, resp1.Reference, resp2.Reference)
	assert.Equal(t, "GSTR-9 draft updated", resp2.Message)
}

func TestMockProvider_GSTR9Save_AfterSubmitted(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	resp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: resp.Reference, RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot save")
}

func TestMockProvider_GSTR9Submit_Success(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"x": 1}, RequestID: uuid.New().String(),
	})

	resp, err := p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
}

func TestMockProvider_GSTR9Submit_NotFound(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9Submit(context.Background(), &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: "NONEXIST", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no GSTR-9 draft found")
}

func TestMockProvider_GSTR9Submit_AlreadySubmitted(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already submitted")
}

func TestMockProvider_GSTR9Submit_Idempotent(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	reqID := uuid.New().String()
	resp1, err := p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: reqID,
	})
	require.NoError(t, err)
	resp2, err := p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: reqID,
	})
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR9File_FullLifecycle(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"tables": "ok"}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})

	fileResp, err := p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "AABCU9603R",
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", fileResp.Status)
	assert.Contains(t, fileResp.ARN, "AR27")
	assert.NotEmpty(t, fileResp.FiledAt)
}

func TestMockProvider_GSTR9File_NotSubmitted(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "X",
		RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must submit GSTR-9")
}

func TestMockProvider_GSTR9File_AlreadyFiled(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "X",
		RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "X",
		RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already filed")
}

func TestMockProvider_GSTR9File_InvalidSignType(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "INVALID", PAN: "X",
		RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sign_type")
}

func TestMockProvider_GSTR9File_EVCMissingOTP(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "EVC", EVOTP: "",
		PAN: "X", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EVC OTP required")
}

func TestMockProvider_GSTR9File_NotFound(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9File(context.Background(), &domain.GSTR9FileRequest{
		Reference: "BAD", SignType: "DSC", PAN: "X",
		GSTIN: testGSTIN, FinancialYear: testFY, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no GSTR-9 found")
}

func TestMockProvider_GSTR9File_Idempotent(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})

	reqID := uuid.New().String()
	fileReq := &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "X",
		RequestID: reqID,
	}
	resp1, err := p.GSTR9File(ctx, fileReq)
	require.NoError(t, err)
	resp2, err := p.GSTR9File(ctx, fileReq)
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR9Status_Found(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})

	resp, err := p.GSTR9Status(ctx, &domain.GSTR9StatusRequest{
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "saved", resp.Status)
	assert.Empty(t, resp.ARN)
}

func TestMockProvider_GSTR9Status_NotFound(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9Status(context.Background(), &domain.GSTR9StatusRequest{
		Reference: "NONEXIST", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no GSTR-9 found")
}

func TestMockProvider_GSTR9Status_Filed(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "X",
		RequestID: uuid.New().String(),
	})

	resp, err := p.GSTR9Status(ctx, &domain.GSTR9StatusRequest{
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "filed", resp.Status)
	assert.NotEmpty(t, resp.ARN)
	assert.NotNil(t, resp.FiledAt)
}

// ---------------------------------------------------------------------------
// GSTR-9C Provider Tests
// ---------------------------------------------------------------------------

func TestMockProvider_GSTR9CSave_Success(t *testing.T) {
	p := NewMockProvider()
	resp, err := p.GSTR9CSave(context.Background(), &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"reconciliation": "data"}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Contains(t, resp.Reference, "GSTR9C-")
}

func TestMockProvider_GSTR9CSave_Idempotent(t *testing.T) {
	p := NewMockProvider()
	reqID := uuid.New().String()
	req := &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: reqID,
	}
	resp1, err := p.GSTR9CSave(context.Background(), req)
	require.NoError(t, err)
	resp2, err := p.GSTR9CSave(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR9CSave_InvalidGSTIN(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9CSave(context.Background(), &domain.GSTR9CSaveRequest{
		GSTIN: "BAD", FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid GSTIN")
}

func TestMockProvider_GSTR9CSave_InvalidFY(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9CSave(context.Background(), &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: "BADFY",
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid financial_year")
}

func TestMockProvider_GSTR9CSave_UpdateExisting(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	resp1, _ := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"v": 1}, RequestID: uuid.New().String(),
	})
	resp2, err := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"v": 2}, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, resp1.Reference, resp2.Reference)
	assert.Equal(t, "GSTR-9C draft updated", resp2.Message)
}

func TestMockProvider_GSTR9CSave_AfterFiled(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9CFile(ctx, &domain.GSTR9CFileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, PAN: "X", RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already filed")
}

func TestMockProvider_GSTR9CFile_Success(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{"recon": "ok"}, RequestID: uuid.New().String(),
	})

	resp, err := p.GSTR9CFile(ctx, &domain.GSTR9CFileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, PAN: "AABCU9603R",
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "success", resp.Status)
	assert.Contains(t, resp.ARN, "AC27")
	assert.Contains(t, resp.Message, "DSC")
}

func TestMockProvider_GSTR9CFile_NotFound(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9CFile(context.Background(), &domain.GSTR9CFileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: "NONEXIST", PAN: "X", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no GSTR-9C found")
}

func TestMockProvider_GSTR9CFile_AlreadyFiled(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9CFile(ctx, &domain.GSTR9CFileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, PAN: "X", RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9CFile(ctx, &domain.GSTR9CFileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, PAN: "X", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already filed")
}

func TestMockProvider_GSTR9CFile_Idempotent(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	reqID := uuid.New().String()
	fileReq := &domain.GSTR9CFileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, PAN: "X", RequestID: reqID,
	}
	resp1, err := p.GSTR9CFile(ctx, fileReq)
	require.NoError(t, err)
	resp2, err := p.GSTR9CFile(ctx, fileReq)
	require.NoError(t, err)
	assert.Equal(t, resp1, resp2)
}

func TestMockProvider_GSTR9CStatus_Found(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})

	resp, err := p.GSTR9CStatus(ctx, &domain.GSTR9CStatusRequest{
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "saved", resp.Status)
}

func TestMockProvider_GSTR9CStatus_NotFound(t *testing.T) {
	p := NewMockProvider()
	_, err := p.GSTR9CStatus(context.Background(), &domain.GSTR9CStatusRequest{
		Reference: "NONEXIST", RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no GSTR-9C found")
}

func TestMockProvider_GSTR9CStatus_Filed(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9CSave(ctx, &domain.GSTR9CSaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9CFile(ctx, &domain.GSTR9CFileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, PAN: "X", RequestID: uuid.New().String(),
	})

	resp, err := p.GSTR9CStatus(ctx, &domain.GSTR9CStatusRequest{
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	require.NoError(t, err)
	assert.Equal(t, "filed", resp.Status)
	assert.NotEmpty(t, resp.ARN)
	assert.NotNil(t, resp.FiledAt)
}

// ---------------------------------------------------------------------------
// Real (Adaequare) Provider Tests
// ---------------------------------------------------------------------------

func TestAdaequareProvider_FailsWithBadCredentials(t *testing.T) {
	a := NewAdaequareProvider("https://gsp.adaequare.com/test", "bad-id", "bad-secret")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := a.Authenticate(ctx)
	assert.Error(t, err, "should fail with bad credentials")
}

func TestMockProvider_GSTR9_StateMachine(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()

	_, err := p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: "NONEXIST", SignType: "DSC", PAN: "X",
		RequestID: uuid.New().String(),
	})
	require.Error(t, err, "can't file without save")

	_, err = p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: "NONEXIST", RequestID: uuid.New().String(),
	})
	require.Error(t, err, "can't submit without save")

	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})

	_, err = p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "X",
		RequestID: uuid.New().String(),
	})
	require.Error(t, err, "can't file before submit")

	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})

	_, err = p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "X",
		RequestID: uuid.New().String(),
	})
	require.NoError(t, err, "can file after submit")
}

func TestMockProvider_GSTR9Submit_AlreadyFiled(t *testing.T) {
	p := NewMockProvider()
	ctx := context.Background()
	saveResp, _ := p.GSTR9Save(ctx, &domain.GSTR9SaveRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Data: map[string]interface{}{}, RequestID: uuid.New().String(),
	})
	p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	p.GSTR9File(ctx, &domain.GSTR9FileRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, SignType: "DSC", PAN: "X",
		RequestID: uuid.New().String(),
	})

	_, err := p.GSTR9Submit(ctx, &domain.GSTR9SubmitRequest{
		GSTIN: testGSTIN, FinancialYear: testFY,
		Reference: saveResp.Reference, RequestID: uuid.New().String(),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already filed")
}
