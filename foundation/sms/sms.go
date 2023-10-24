// Package sms provide utilities for sending sms
// esms docs: https://developers.esms.vn
package sms

import (
	"net/http"
)

// Codes and definitions: https://developers.esms.vn/esms-api/bang-ma-loi
const (
	codeSuccess                 = "100" // Request được gửi đến ViHAT thành công
	codeInvalidAPISecretKey     = "101" // Sai thông tin ApiKey hoặc SecretKey
	codePriceTableNotFound      = "102" // Không có bảng giá
	codeInsufficientBalance     = "103" // Không đủ số dư để gửi tin
	codeBrandNameNotFound       = "104" // Brandname/ Mã cuộc gọi không tồn tại
	codeSMSIDNotFound           = "105" // Không tìm thấy mã tin nhắn/ Mã cuộc gọi trên hệ thống
	codeRecordFileNotFound      = "106" // File ghi âm không tồn tại
	codeMinNotMatched           = "107" // Mỗi request phải có ít nhất 30 số mới được duyệt
	codeDuplicatedRequestID     = "124" // Trùng RequestId khi gửi tin
	codeInvalidTemplate         = "146" // Sai template chăm sóc khách hàng
	codeCarrierNotRegister      = "177" // Nhà mạng chưa được đăng ký
	codeRequestProcessing       = "199" // Yêu cầu này đã tồn tại, đang xử lý
	codeMissingSMSType          = "300" // Thiếu loại tin nhắn
	codeServiceTypeNotSupported = "300" // Mã của loại dịch vụ chưa hỗ trợ
	codeInvalidOTP              = "171" // Mã đã được sử dụng, mã hết hạn hoặc mã không áp dụng cho số điện thoại cần check.
	codeMaxNumberExceeded       = "120" // Danh sách số điện thoại gửi tin vượt quá giới hạn.
	codeInvalidPayload          = "201" // Không đúng payload theo quy định
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
