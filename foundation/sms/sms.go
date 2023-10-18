// Package sms provide utilities for sending sms
// esms docs: https://esms.vn/eSMS.vn_TailieuAPI.pdf
package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

// todo: create http client wrapper support automatic tracing and logging
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

	// Message represents an sms message
	Message struct {
		RequestID uuid.UUID
		Phone     string
		Content   string
		Type      Type
		DryRun    bool
	}
)

type (
	reqData struct {
		APIKey      string `json:"ApiKey,omitempty"`      // *ApiKey của tài khoản
		SecretKey   string `json:"SecretKey,omitempty"`   // *Secret key của tài khoản
		Content     string `json:"Content,omitempty"`     // *Nội dung tin nhắn
		Phone       string `json:"Phone,omitempty"`       // *Số điện thoại nhận tin
		IsUnicode   int32  `json:"IsUnicode,omitempty"`   // Gửi nội dung có dấu 1: Có dấu 0: Không dấu
		Sandbox     int32  `json:"Sandbox,omitempty"`     // 1: Tin thử nghiệm, không gửi tin nhắn, chỉ trả về kết quả SMS, tin không lưu hệ thống và không trừ tiền. 0: Không thử nghiệm, tin đi thật.
		BrandName   string `json:"Brandname,omitempty"`   // *Tên Brandname (tên công ty hay tổ chức khi gửi tin sẽ hiển thị trên tin nhắn đó). Chú ý: sẽ phải đăng ký trước khi sử dụng.
		SMSType     Type   `json:"SmsType,omitempty"`     // *Loại tin nhắn
		RequestID   string `json:"RequestId,omitempty"`   // ID tin nhắn của đối tác, dùng để kiểm tra xem tin nhắn này đã được tiếp nhận trước đó hay chưa
		CallbackURL string `json:"CallBackUrl,omitempty"` // Khi có kết quả của tin nhắn sẽ gọi về URL này
		SendDate    string `json:"SendDate,omitempty"`    // Thời gian hẹn gửi của tin. Không truyền khi tin muốn tin nhắn gửi đi liền. Định dạng: yyyy-mm-dd hh:MM:ss
		CampaignID  string `json:"campaignid,omitempty"`  // Tên chiến dịch gửi tin
	}

	respData struct {
		CodeResult      string `json:"CodeResult"`
		CountRegenerate int32  `json:"CountRegenerate"`
		SMSID           string `json:"SMSID"` // ID của tin nhắn mới được tạo ra trên hệ thống eSMS. Dùng ID này để query lấy trạng thái tin nhắn.
	}

	sms struct {
		address   string
		apiKey    string
		secretKey string
		brandName string
		client    *http.Client
	}
)

func New(cfg Config) *sms {
	if cfg.Client == nil {
		cfg.Client = &http.Client{}
	}
	return &sms{
		address:   cfg.Address,
		apiKey:    cfg.APIKey,
		secretKey: cfg.SecretKey,
		client:    cfg.Client,
		brandName: cfg.BrandName,
	}
}

func (s *sms) Send(ctx context.Context, msg Message) (smsID string, err error) {
	var sandbox int32
	if msg.DryRun {
		sandbox = 1
	}

	resp, err := s.send(ctx, reqData{
		RequestID: msg.RequestID.String(),
		Phone:     msg.Phone,
		SMSType:   msg.Type,
		Content:   msg.Content,
		APIKey:    s.apiKey,
		SecretKey: s.secretKey,
		BrandName: s.brandName,
		Sandbox:   sandbox,
		IsUnicode: 1,
	})
	if err != nil {
		return "", fmt.Errorf("sms.Send error=%+v", err)
	}

	if resp.CodeResult != codeSuccess {
		return "", fmt.Errorf("sms.Send received failed result, resp=%+v", resp)
	}

	if len(resp.SMSID) == 0 {
		return "", fmt.Errorf("sms.Send SMSID is empty, resp=%+v", resp)
	}

	return resp.SMSID, nil
}

func (s *sms) send(ctx context.Context, body reqData) (respData, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(body); err != nil {
		return respData{}, fmt.Errorf("encode data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.address+"/SendMultipleMessage_V4_post_json", &b)
	if err != nil {
		return respData{}, fmt.Errorf("create request error=%+v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return respData{}, fmt.Errorf("do error=%+v", err)
	}
	defer resp.Body.Close()

	var response respData
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return respData{}, fmt.Errorf("json decode error=%+v", err)
	}
	return response, nil
}
