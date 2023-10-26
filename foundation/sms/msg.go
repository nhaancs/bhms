package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type (

	// Message represents an sms message
	Message struct {
		RequestID uuid.UUID
		Phone     string
		Content   string
		Type      Type
		DryRun    bool
	}
	msgReqData struct {
		APIKey      string `json:"ApiKey,omitempty"`      // *ApiKey của tài khoản
		SecretKey   string `json:"SecretKey,omitempty"`   // *Secret key của tài khoản
		Content     string `json:"Content,omitempty"`     // *Nội dung tin nhắn
		Phone       string `json:"Phone,omitempty"`       // *Số điện thoại nhận tin
		IsUnicode   int32  `json:"IsUnicode,omitempty"`   // Gửi nội dung có dấu 1: Có dấu 0: Không dấu
		Sandbox     int32  `json:"Sandbox,omitempty"`     // 1: Tin thử nghiệm, không gửi tin nhắn, chỉ trả về kết quả SMS, tin không lưu hệ thống và không trừ tiền. 0: Không thử nghiệm, tin đi thật.
		BrandName   string `json:"Brandname,omitempty"`   // *Tên Brandname (tên công ty hay tổ chức khi gửi tin sẽ hiển thị trên tin nhắn đó). Chú ý: sẽ phải đăng ký trước khi sử dụng.
		SMSType     string `json:"SmsType,omitempty"`     // *Loại tin nhắn
		RequestID   string `json:"RequestId,omitempty"`   // ID tin nhắn của đối tác, dùng để kiểm tra xem tin nhắn này đã được tiếp nhận trước đó hay chưa
		CallbackURL string `json:"CallBackUrl,omitempty"` // Khi có kết quả của tin nhắn sẽ gọi về URL này
		SendDate    string `json:"SendDate,omitempty"`    // Thời gian hẹn gửi của tin. Không truyền khi tin muốn tin nhắn gửi đi liền. Định dạng: yyyy-mm-dd hh:MM:ss
		CampaignID  string `json:"campaignid,omitempty"`  // Tên chiến dịch gửi tin
	}

	msgRespData struct {
		CodeResult      string `json:"CodeResult"`
		CountRegenerate int32  `json:"CountRegenerate"`
		SMSID           string `json:"SMSID"` // ID của tin nhắn mới được tạo ra trên hệ thống eSMS. Dùng ID này để query lấy trạng thái tin nhắn.
	}
)

// Send an SMS message to a phone number
func (s *SMS) Send(ctx context.Context, msg Message) (smsID string, err error) {
	var sandbox int32
	if msg.DryRun {
		sandbox = 1
	}

	resp, err := s.send(ctx, msgReqData{
		RequestID: msg.RequestID.String(),
		Phone:     msg.Phone,
		SMSType:   msg.Type.ID(),
		Content:   msg.Content,
		APIKey:    s.apiKey,
		SecretKey: s.secretKey,
		BrandName: s.brandName,
		Sandbox:   sandbox,
		IsUnicode: 1,
	})
	if err != nil {
		return "", fmt.Errorf("sms.Send error=%w", err)
	}

	if resp.CodeResult != codeSuccess {
		return "", fmt.Errorf("sms.Send received failed result, resp=%+v", resp)
	}

	if len(resp.SMSID) == 0 {
		return "", fmt.Errorf("sms.Send SMSID is empty, resp=%+v", resp)
	}

	return resp.SMSID, nil
}

func (s *SMS) send(ctx context.Context, body msgReqData) (msgRespData, error) {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(body); err != nil {
		return msgRespData{}, fmt.Errorf("encode data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.address+"/MainService.svc/json/SendMultipleMessage_V4_post_json", &b)
	if err != nil {
		return msgRespData{}, fmt.Errorf("create request error=%w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return msgRespData{}, fmt.Errorf("do error=%w", err)
	}
	defer resp.Body.Close()

	var response msgRespData
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return msgRespData{}, fmt.Errorf("json decode error=%w", err)
	}
	return response, nil
}
