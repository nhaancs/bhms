package sms

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/http"
)

type (
	// OTPInfo provides info for sending OTP
	OTPInfo struct {
		Phone string
	}

	// VerifyOTPInfo provides info for verifying OTP
	VerifyOTPInfo struct {
		Phone string
		Code  string
	}

	otpReqData struct {
		Phone         string `url:"Phone"`         // *Số điện thoại nhận code
		APIKey        string `url:"ApiKey"`        // *ApiKey của tài khoản
		SecretKey     string `url:"SecretKey"`     // *Secret key của tài khoản
		TimeAlive     int32  `url:"TimeAlive"`     // Thời gian hiệu lực của mã code. Đơn vị tính: Phút. Giới hạn từ 2 phút đến 15 phút
		NumCharOfCode int32  `url:"NumCharOfCode"` // Số lượng ký tự của mã code. Mặc định khi không truyền là 6 ký tự
		BrandName     string `url:"Brandname"`     // *Tên Brandname (tên công ty hay tổ chức khi gửi tin sẽ hiển thị trên tin nhắn đó). Chú ý: sẽ phải đăng ký trước khi sử dụng.
		Type          string `url:"Type"`          // *Loại tin nhắn
		Message       string `url:"message"`       // Nội dung tin nhắn. Thay thế mã code OTP bằng {OTP} - Ví dụ: Muốn nhận tin nhắn về máy với nội dung: "686868 la ma xac minh dang ky Baotrixemay cua ban". Thì ở mục message truyền: "{OTP} la ma xac minh dang ky Baotrixemay cua ban".
		IsNumber      int32  `url:"IsNumber"`      // 0: code chứa chữ và số 1: code chỉ chứa số
	}

	otpRespData struct {
		CodeResult      string `json:"CodeResult"`
		CountRegenerate int32  `json:"CountRegenerate"`
		SMSID           string `json:"SMSID"` // ID của tin nhắn mới được tạo ra trên hệ thống eSMS. Dùng ID này để query lấy trạng thái tin nhắn.
	}

	checkOTPReqData struct {
		Phone     string `url:"Phone"`     // *Số điện thoại nhận code
		APIKey    string `url:"ApiKey"`    // *ApiKey của tài khoản
		SecretKey string `url:"SecretKey"` // *Secret key của tài khoản
		Code      string `url:"Code"`      // Mã code cần check
	}

	checkOTPRespData struct {
		CodeResult      string `json:"CodeResult"`
		CountRegenerate int32  `json:"CountRegenerate"`
		ErrorMessage    string `json:"ErrorMessage"`
	}
)

// SendOTP sent an OTP message to a phone number
func (s *SMS) SendOTP(ctx context.Context, otp OTPInfo) (smsID string, err error) {
	msg := fmt.Sprintf("{OTP} la ma xac nhan %s cua ban", s.brandName)
	resp, err := s.sendOTP(ctx, otpReqData{
		Phone:         otp.Phone,
		APIKey:        s.apiKey,
		SecretKey:     s.secretKey,
		TimeAlive:     5,
		NumCharOfCode: 6,
		BrandName:     s.brandName,
		Type:          TypeBrandName.ID(),
		Message:       msg,
		IsNumber:      1,
	})
	if err != nil {
		return "", fmt.Errorf("sms.SendOTP error=%+v", err)
	}

	if resp.CodeResult != codeSuccess {
		return "", fmt.Errorf("sms.SendOTP received failed result, resp=%+v", resp)
	}

	if len(resp.SMSID) == 0 {
		return "", fmt.Errorf("sms.SendOTP SMSID is empty, resp=%+v", resp)
	}

	return resp.SMSID, nil
}

// CheckOTP verify user OTP
func (s *SMS) CheckOTP(ctx context.Context, otp VerifyOTPInfo) error {
	resp, err := s.checkOTP(ctx, checkOTPReqData{
		Phone:     otp.Phone,
		Code:      otp.Code,
		APIKey:    s.apiKey,
		SecretKey: s.secretKey,
	})
	if err != nil {
		return fmt.Errorf("sms.CheckOTP error=%+v", err)
	}

	if resp.CodeResult != codeSuccess {
		return fmt.Errorf("sms.CheckOTP received failed result, resp=%+v", resp)
	}

	return nil
}

func (s *SMS) sendOTP(ctx context.Context, body otpReqData) (otpRespData, error) {
	v, err := query.Values(body)
	if err != nil {
		return otpRespData{}, fmt.Errorf("query.Values: %w", err)
	}

	url := fmt.Sprintf("%s/MainService.svc/json/SendMessageAutoGenCode_V4_get?%s", s.address, v.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return otpRespData{}, fmt.Errorf("create request error=%+v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return otpRespData{}, fmt.Errorf("do error=%+v", err)
	}
	defer resp.Body.Close()

	var response otpRespData
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return otpRespData{}, fmt.Errorf("json decode error=%+v", err)
	}
	return response, nil
}

func (s *SMS) checkOTP(ctx context.Context, body checkOTPReqData) (checkOTPRespData, error) {
	v, err := query.Values(body)
	if err != nil {
		return checkOTPRespData{}, fmt.Errorf("query.Values: %w", err)
	}

	url := fmt.Sprintf("%s/MainService.svc/json/CheckCodeGen_V4_get?%s", s.address, v.Encode())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return checkOTPRespData{}, fmt.Errorf("create request error=%+v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return checkOTPRespData{}, fmt.Errorf("do error=%+v", err)
	}
	defer resp.Body.Close()

	var response checkOTPRespData
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return checkOTPRespData{}, fmt.Errorf("json decode error=%+v", err)
	}
	return response, nil
}
