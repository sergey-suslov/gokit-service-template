package endpoints

import (
	"context"
	"encoding/json"
	"net/http"
)

func DefaultRequestEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

type UserClaimable interface {
	SetUserClaim(claim UserClaim)
}

func DefaultRequestDecoder(decoder func(r *http.Request) (UserClaimable, error)) func(_ context.Context, r *http.Request) (interface{}, error) {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		userClaim, err := GetUserClaimFromRequest(r)
		if err != nil {
			return nil, err
		}

		request, err := decoder(r)
		if err != nil {
			return nil, err
		}
		request.SetUserClaim(*userClaim)
		return request, nil
	}
}
