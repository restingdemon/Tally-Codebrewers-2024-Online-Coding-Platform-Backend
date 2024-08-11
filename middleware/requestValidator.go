package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"worldwide-coders/utils"

	"github.com/gorilla/mux"
)

func CheckHTTPAuthorization(r *http.Request, ctx context.Context, userType string, userEmail string) (context.Context, error) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/users/get"):
		queryParams := r.URL.Query()
		email := queryParams.Get("email")
		if userType == utils.UserRole {
			ctx = context.WithValue(ctx, "email", userEmail)
			return ctx, nil
		} else {
			ctx = context.WithValue(ctx, "email", email)
			return ctx, nil
		}
	case strings.HasPrefix(r.URL.Path, "/users/update"):
		// extracting email from path parameters
		vars := mux.Vars(r)
		email, ok := vars["email"]
		if !ok {
			return ctx, fmt.Errorf("no email provided")
		}
		if userType == "superadmin" {
			ctx = context.WithValue(ctx, "email", email)
			return ctx, nil
		}
		if email != userEmail {
			return ctx, fmt.Errorf("you can only update your own details")
		}

		ctx = context.WithValue(ctx, "email", email)
		return ctx, nil

	case strings.HasPrefix(r.URL.Path, "/problems/upload"):
		ctx = context.WithValue(ctx, "email", userEmail)
		ctx = context.WithValue(ctx, "role", userType)
		return ctx, nil

	case strings.HasPrefix(r.URL.Path, "/problems/getnotvisible"):
		ctx = context.WithValue(ctx, "email", userEmail)
		ctx = context.WithValue(ctx, "role", userType)
		return ctx, nil

	case strings.HasPrefix(r.URL.Path, "/problems/update"):
		ctx = context.WithValue(ctx, "email", userEmail)
		ctx = context.WithValue(ctx, "role", userType)
		return ctx, nil

	case strings.HasPrefix(r.URL.Path, "/contests/register"):
		ctx = context.WithValue(ctx, "email", userEmail)
		return ctx, nil

	case strings.HasPrefix(r.URL.Path, "/contests/create"):
		ctx = context.WithValue(ctx, "email", userEmail)
		return ctx, nil

	case strings.HasPrefix(r.URL.Path, "/contests/get/registrations"):
		vars := mux.Vars(r)
		contestId, ok := vars["contestId"]
		if !ok {
			return ctx, fmt.Errorf("no contestId Id provided")
		}

		if userType == utils.SuperAdminRole {
			ctx = context.WithValue(ctx, "contestId", contestId)
			ctx = context.WithValue(ctx, "role", utils.SuperAdminRole)
			ctx = context.WithValue(ctx, "userType", userType)
			return ctx, nil
		} else if userType == utils.UserRole {
			ctx = context.WithValue(ctx, "contestId", contestId)
			ctx = context.WithValue(ctx, "role", userType)
			ctx = context.WithValue(ctx, "email", userEmail)
			ctx = context.WithValue(ctx, "userType", userType)
			return ctx, nil
		}

	case strings.HasPrefix(r.URL.Path, "/contests/check/registrations/"):
		vars := mux.Vars(r)
		contestId, ok := vars["contestId"]
		if !ok {
			return ctx, fmt.Errorf("no Contest Id provided")
		}
		ctx = context.WithValue(ctx, "contestId", contestId)
		ctx = context.WithValue(ctx, "email", userEmail)
		return ctx, nil

	}
	// Default to allowing access if the route is not explicitly handled
	return ctx, nil
}
