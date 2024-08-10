package middleware

import (
	"context"
	"net/http"
	
)

func CheckHTTPAuthorization(r *http.Request, ctx context.Context, userType string, userEmail string) (context.Context, error) {
	switch {
	// case strings.HasPrefix(r.URL.Path, "/users/get"):
	// 	queryParams := r.URL.Query()
	// 	email := queryParams.Get("email")
	// 	if userType == utils.UserRole {
	// 		ctx = context.WithValue(ctx, "email", userEmail)
	// 		return ctx, nil
	// 	} else {
	// 		ctx = context.WithValue(ctx, "email", email)
	// 		return ctx, nil
	// 	}
	
	}
	// Default to allowing access if the route is not explicitly handled
	return ctx, nil
}
