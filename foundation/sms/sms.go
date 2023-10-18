// Package sms provide utilities for sending sms
// esms docs: https://esms.vn/eSMS.vn_TailieuAPI.pdf
package sms

import "net/http"

// This provides a default client configuration, but it's recommended
// this is replaced by the client with application specific settings
var defaultClient = http.Client{}

// Config represents the mandatory settings needed to work with sms.
type Config struct {
	Address   string
	APIKey    string
	SecretKey string
	Client    *http.Client
}

type sms struct {
	address   string
	apiKey    string
	secretKey string
	client    *http.Client
}

func New(cfg Config) *sms {
	if cfg.Client == nil {
		cfg.Client = &defaultClient
	}
	return &sms{
		address:   cfg.Address,
		apiKey:    cfg.APIKey,
		secretKey: cfg.SecretKey,
		client:    cfg.Client,
	}
}

type request struct {
	APIKey      string `json:"ApiKey"`
	SecretKey   string `json:"SecretKey"`
	Content     string `json:"Content"`
	Phone       string `json:"Phone"`
	IsUnicode   string `json:"IsUnicode"`
	BrandName   string `json:"Brandname"`
	SMSType     string `json:"SmsType"`
	RequestID   string `json:"RequestId"`
	CallbackURL string `json:"CallBackUrl"`
	SendDate    string `json:"SendDate"`
	CampaignID  string `json:"campaignid"`
}

type response struct {
}

type OTPMessage struct {
	BrandName string
	Phone     string
	Content   string
}

func (s *sms) SendOTP(in OTPMessage) error {
	return nil
}

func (s *sms) send(req request) (response, error) {
	return response{}, nil
}
