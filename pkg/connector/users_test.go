package connector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/conductorone/baton-privx/pkg/connector/client"
	"github.com/conductorone/baton-sdk/pkg/pagination"
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

		resources, token, annotations, err := userBuilder.List(ctx, nil, &pagination.Token{})
		require.Nil(t, err)

		// Assert the returned user has an ID.
		require.NotNil(t, resources)
		require.Len(t, resources, 3)
		require.NotEmpty(t, resources[0].Id)

		require.Equal(t, "", token)
		require.Len(t, annotations, 0)
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

		paginationToken := pagination.Token{
			Token: "100",
			Size:  3,
		}

		resources, token, annotations, err := userBuilder.List(ctx, nil, &paginationToken)
		require.Nil(t, err)

		// Assert the returned user has an ID.
		require.NotNil(t, resources)
		require.Len(t, resources, 3)
		require.NotEmpty(t, resources[0].Id)

		// Should look for second page.
		require.Equal(t, "103", token)

		require.Len(t, annotations, 0)
	})
}
