//go:build integration

package provider_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/complai/complai/services/go/gstn-gateway-service/internal/domain"
	"github.com/complai/complai/services/go/gstn-gateway-service/internal/provider"
)

func newAdaequareProvider(t *testing.T) *provider.AdaequareProvider {
	t.Helper()
	clientID := os.Getenv("ADAEQUARE_CLIENT_ID")
	clientSecret := os.Getenv("ADAEQUARE_CLIENT_SECRET")
	baseURL := os.Getenv("ADAEQUARE_BASE_URL")
	if clientID == "" || clientSecret == "" {
		t.Skip("ADAEQUARE_CLIENT_ID and ADAEQUARE_CLIENT_SECRET not set; skipping integration test")
	}
	if baseURL == "" {
		baseURL = "https://gsp.adaequare.com/test"
	}
	return provider.NewAdaequareProvider(baseURL, clientID, clientSecret)
}

func TestAdaequareAuthenticate(t *testing.T) {
	p := newAdaequareProvider(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := p.Authenticate(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.NotEmpty(t, resp.AccessToken, "access_token should not be empty")
	assert.Equal(t, "bearer", resp.TokenType)
	assert.Greater(t, resp.ExpiresIn, 0)
	assert.Equal(t, "gsp", resp.Scope)
	t.Logf("authenticated: token=%s...%s, expires_in=%d, jti=%s",
		resp.AccessToken[:10], resp.AccessToken[len(resp.AccessToken)-10:],
		resp.ExpiresIn, resp.JTI)
}

func TestAdaequareAuthTokenCaching(t *testing.T) {
	p := newAdaequareProvider(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp1, err := p.Authenticate(ctx)
	require.NoError(t, err)

	resp2, err := p.Authenticate(ctx)
	require.NoError(t, err)

	assert.Equal(t, resp1.AccessToken, resp2.AccessToken, "second call should return cached token")
}

func TestAdaequareGSTR1Status(t *testing.T) {
	p := newAdaequareProviderWithCreds(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	gstin := os.Getenv("ADAEQUARE_TEST_GSTIN")
	if gstin == "" {
		t.Skip("ADAEQUARE_TEST_GSTIN not set")
	}

	resp, err := p.GSTR1Status(ctx, &domain.GSTR1StatusRequest{
		GSTIN:     gstin,
		RetPeriod: "012025",
		RequestID: "inttest-status-1",
	})
	// We expect either success or a provider-level error (e.g. RET11402);
	// what matters is the request reached Adaequare without a client-side failure.
	if err != nil {
		t.Logf("GSTR1Status returned provider error (expected in sandbox): %v", err)
		assert.Contains(t, err.Error(), "adaequare:", "error should be from adaequare, not a client bug")
		return
	}
	t.Logf("GSTR1Status succeeded: %+v", resp)
}

func newAdaequareProviderWithCreds(t *testing.T) *provider.AdaequareProvider {
	t.Helper()
	clientID := os.Getenv("ADAEQUARE_CLIENT_ID")
	clientSecret := os.Getenv("ADAEQUARE_CLIENT_SECRET")
	baseURL := os.Getenv("ADAEQUARE_BASE_URL")
	username := os.Getenv("ADAEQUARE_TEST_USERNAME")
	password := os.Getenv("ADAEQUARE_TEST_PASSWORD")
	if clientID == "" || clientSecret == "" {
		t.Skip("ADAEQUARE_CLIENT_ID and ADAEQUARE_CLIENT_SECRET not set")
	}
	if baseURL == "" {
		baseURL = "https://gsp.adaequare.com/test"
	}
	var opts []provider.AdaequareOption
	if username != "" && password != "" {
		opts = append(opts, provider.WithTaxpayerCredentials(username, password))
	}
	return provider.NewAdaequareProvider(baseURL, clientID, clientSecret, opts...)
}
