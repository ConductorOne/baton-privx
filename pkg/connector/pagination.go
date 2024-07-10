package connector

import (
	"strconv"
)

func parsePageToken(token string) (int, error) {
	if token == "" {
		return 0, nil
	}
	return strconv.Atoi(token)
}
