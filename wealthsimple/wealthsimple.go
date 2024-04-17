package wealthsimple

import (
	"encoding/json"
	"net/http"
)

const authTokenRoute = "https://api.production.wealthsimple.com/v1/oauth/v2/token"
const otpClaimHeader = "X-Wealthsimple-Otp-Options"

type WealthSimpleLogin struct {
	Email    string
	Password string
	OtpClaim string
}

func (ws *WealthSimpleLogin) InitOtpClaim() error {
	jsonForm, err := json.Marshal(map[string]string{
		"grant_type":     "password",
		"username":       ws.Email,
		"password":       ws.Password,
		"skip_provision": "true",
		"otp_claim":      "null",
		"scope":          "invest.read invest.write trade.read trade.write tax.read tax.write",
	})
	if err != nil {
		return err
	}

	req, err := NewRequest(http.MethodPost, authTokenRoute, jsonForm)
	if err != nil {
		return err
	}

	req.AddBaseHeaders()

	resp, err := req.MakeRequest()
	if err != nil {
		return err
	}

	ws.OtpClaim = resp.GetHeader(otpClaimHeader)

	return nil
}
