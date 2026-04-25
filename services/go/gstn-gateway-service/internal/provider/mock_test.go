package provider

import (
	"context"
	"testing"

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
