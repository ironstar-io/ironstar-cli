package api

import (
	"encoding/json"
	"time"

	"gitlab.com/ironstar-io/ironstar-cli/internal/types"
)

func postAuthTokenRefresh(creds types.Keylink) (*types.AuthResponseBody, error) {
	newReq := Request{
		RunTokenRefresh: false,
		Credentials:     creds,
		Method:          "POST",
		Path:            "/auth/token/refresh",
		MapStringPayload: map[string]interface{}{
			"expiry": time.Now().AddDate(0, 0, 14).UTC().Format(time.RFC3339),
		},
	}

	res, err := newReq.NankaiSend()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, res.HandleFailure()
	}

	b := &types.AuthResponseBody{}
	err = json.Unmarshal(res.Body, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
