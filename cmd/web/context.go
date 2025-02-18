package main

// contextKey is a type for keys used in context.WithValue.
type contextKey string

const (
	// isAuthenticatedContextKey is the key used to store authentication status in the context.
	isAuthenticatedContextKey = contextKey("isAuthenticated")

	// nonceContextKey is the key used to store the nonce in the context.
	nonceContextKey = contextKey("nonce")
)
