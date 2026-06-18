package connector

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/conductorone/baton-privx/pkg/connector/client"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/stretchr/testify/require"
)

// newMockPrivXServer creates a test server that handles the PrivX OAuth flow
// (Authorization Code Grant) and serves fixture data for API endpoints.
func newMockPrivXServer(fixturePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/api/v1/oauth/authorize":
			// The SDK sends state as a query param; echo it back as the token so
			// the login endpoint can return the same state value.
			state := r.URL.Query().Get("state")
			w.Header().Set("Location", "/privx/oauth-callback?token="+state)
			w.WriteHeader(http.StatusTemporaryRedirect)

		case "/auth/api/v1/login":
			var body struct {
				Token string `json:"token"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			w.Header().Set(uhttp.ContentType, "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"code":  "test-code",
				"state": body.Token, // token == original state, satisfying the SDK check
			})

		case "/auth/api/v1/oauth/token":
			w.Header().Set(uhttp.ContentType, "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"access_token": "test-token",
				"token_type":   "bearer",
				"expires_in":   3600,
			})

		default:
			w.Header().Set(uhttp.ContentType, "application/json")
			w.WriteHeader(http.StatusOK)
			data, _ := os.ReadFile(fixturePath)
			_, _ = w.Write(data)
		}
	}))
}

func TestUsersList(t *testing.T) {
	ctx := context.Background()
	t.Run("should receive users", func(t *testing.T) {
		server := newMockPrivXServer("./client/fixtures/search_page_0.json")
		defer server.Close()

		privXClient, err := client.NewPrivXClientWithClientSecret(
			ctx,
			server.URL,
			"clientId",
			"clientSecret",
		)
		require.Nil(t, err)
		userBuilder := newUserBuilder(*privXClient)

		resources, results, err := userBuilder.List(ctx, nil, resource.SyncOpAttrs{})
		require.Nil(t, err)

		// Assert the returned user has an ID.
		require.NotNil(t, resources)
		require.Len(t, resources, 3)
		require.NotEmpty(t, resources[0].Id)

		require.Nil(t, results)
	})

	t.Run("should paginate", func(t *testing.T) {
		server := newMockPrivXServer("./client/fixtures/search_page_0.json")
		defer server.Close()

		privXClient, err := client.NewPrivXClientWithClientSecret(
			ctx,
			server.URL,
			"clientId",
			"clientSecret",
		)
		require.Nil(t, err)
		userBuilder := newUserBuilder(*privXClient)

		opts := resource.SyncOpAttrs{}
		opts.PageToken.Token = "100"
		opts.PageToken.Size = 3

		resources, results, err := userBuilder.List(ctx, nil, opts)
		require.Nil(t, err)

		// Assert the returned user has an ID.
		require.NotNil(t, resources)
		require.Len(t, resources, 3)
		require.NotEmpty(t, resources[0].Id)

		// Should look for second page.
		require.NotNil(t, results)
		require.Equal(t, "103", results.NextPageToken)
	})
}
