package main

import (
	"context"
	"movieDB/internal/data"
	"net/http"
)

// contextKey avoids naming collision. Set the user context as type: contextKey.
type contextKey string

const userContextKey = contextKey("users")

// Return a new Context with User embedded in the contextKey.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// Retrieve the user from the current context.
func (app *application) contextGetUser(r *http.Request) *data.User {
	//
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
