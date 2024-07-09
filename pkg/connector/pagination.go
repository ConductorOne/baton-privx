package connector

import (
	"strconv"

	"github.com/conductorone/baton-sdk/pkg/pagination"
)

const (
	ResourcePageSizeDefault = 100
)

// parsePageToken takes a pagination and parses it into offset and limit (in
// that order) and picks sensible defaults.
func parsePageToken(pageToken *pagination.Token) (int, int, error) {
	offset := 0
	limit := ResourcePageSizeDefault

	if pageToken == nil {
		return offset, limit, nil
	}

	if pageToken.Token != "" {
		offsetValue, err := strconv.Atoi(pageToken.Token)
		if err != nil {
			return 0, 0, err
		}
		offset = offsetValue
	}

	if pageToken.Size > 0 {
		limit = pageToken.Size
	}

	return offset, limit, nil
}
