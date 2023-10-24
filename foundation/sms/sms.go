// Package sms provide utilities for sending sms
// esms docs: https://esms.vn/eSMS.vn_TailieuAPI.pdf
package sms

import (
	"net/http"
)

const (
	codeSuccess = "100"
)

type (
	// Config represents the mandatory settings needed to work with sms.
	Config struct {
		Address   string
		APIKey    string
		SecretKey string
		BrandName string
		Client    *http.Client
	}
	SMS struct {
		address   string
		apiKey    string
		secretKey string
		brandName string
		client    *http.Client
	}
)

func New(cfg Config) *SMS {
	if cfg.Client == nil {
		// This provides a default client configuration, but it's recommended
		// this is replaced by the user with application specific settings using
		// the WithClient function at the time a GraphQL is constructed.
		cfg.Client = &http.Client{}
	}
	return &SMS{
		address:   cfg.Address,
		apiKey:    cfg.APIKey,
		secretKey: cfg.SecretKey,
		client:    cfg.Client,
		brandName: cfg.BrandName,
	}
}
