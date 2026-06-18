package connector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/conductorone/baton-privx/pkg/connector/client"
	"github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/stretchr/testify/require"
)

func TestUsersList(t *testing.T) {
	ctx := context.Background()
	t.Run("should receive users", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(
				func(writer http.ResponseWriter, request *http.Request) {
					writer.Header().Set(uhttp.ContentType, "application/json")
					writer.WriteHeader(http.StatusOK)
					json, err := os.ReadFile("./client/fixtures/search_page_0.json")
					require.Nil(t, err)
					_, err = writer.Write(json)
					if err != nil {
						return
					}
				},
			),
		)
		defer server.Close()

		privXClient, err := client.NewPrivXClient(
			ctx,
			server.URL,
			"apiClientId",
			"apiClientSecret",
			"oauthClientId",
			"oauthClientSecret",
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
		server := httptest.NewServer(
			http.HandlerFunc(
				func(writer http.ResponseWriter, request *http.Request) {
					writer.Header().Set(uhttp.ContentType, "application/json")
					writer.WriteHeader(http.StatusOK)
					json, err := os.ReadFile("./client/fixtures/search_page_0.json")
					require.Nil(t, err)
					_, err = writer.Write(json)
					if err != nil {
						return
					}
				},
			),
		)
		defer server.Close()

		privXClient, err := client.NewPrivXClient(
			ctx,
			server.URL,
			"apiClientId",
			"apiClientSecret",
			"oauthClientId",
			"oauthClientSecret",
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
