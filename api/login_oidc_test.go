package api

import (
	"context"
	"testing"

	"github.com/semaphoreui/semaphore/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// d2c2b4a6: an OIDC provider's auth_style picks how the client secret is sent to the token endpoint. Unset, oauth2
// detects it, as before the option existed.
func TestGetOidcProvider_AuthStyle(t *testing.T) {
	original := util.Config
	t.Cleanup(func() { util.Config = original })

	for style, want := range map[util.OidcProviderAuthStyle]oauth2.AuthStyle{
		"":                             oauth2.AuthStyleAutoDetect,
		util.OidcProviderAuthInParams:  oauth2.AuthStyleInParams,
		util.OidcProviderAuthSInHeader: oauth2.AuthStyleInHeader,
	} {
		p := util.OidcProvider{RedirectURL: "https://semaphore.example/api/auth/oidc/p/redirect", AuthStyle: style}
		p.Endpoint.AuthURL = "https://idp.example/auth"
		p.Endpoint.TokenURL = "https://idp.example/token"
		util.Config = &util.ConfigType{OidcProviders: map[string]util.OidcProvider{"p": p}}

		_, cfg, err := getOidcProvider("p", context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, want, cfg.Endpoint.AuthStyle, "auth_style %q", style)
		assert.Equal(t, "https://idp.example/token", cfg.Endpoint.TokenURL)
	}
}
