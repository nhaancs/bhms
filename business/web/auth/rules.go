package auth

import (
	_ "embed"
)

// These the current set of rules we have for auth.
const (
	RuleAuthenticate = "auth"
)

// Package name of our rego code.
const (
	opaPackage string = "nhaancs.rego"
)

// Core OPA policies.
var (
	//go:embed rego/authentication.rego
	opaAuthentication string
)
