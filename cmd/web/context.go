package main

//create this new type, so that we can
//distinguish context keys made by this application
//and not have naming collisions with third party packages
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
