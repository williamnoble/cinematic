package main

import (
	"context"
	"greenlight/internal/data"
	"net/http"
)

// customContextKey to avoid name collison. Set the user context as type: contextKey.
type contextKey string

const userContextKey = contextKey("users")

// Return a new Context with the "User" embedded with a userContextKey
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}
func (app *application) contextGetUser(r *http.Request) *data.User {
	//
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
