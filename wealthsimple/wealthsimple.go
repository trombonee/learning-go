package wealthsimple

import (
	"encoding/json"
	"fmt"
	"learning-go/db"
	"learning-go/env"
	"net/http"
	"strconv"
)

const authTokenRoute = "https://api.production.wealthsimple.com/v1/oauth/v2/token"
const otpClaimHeader = "X-Wealthsimple-Otp-Options"

var baseHeaders = map[string]string{
	"accept":                "application/json",
	"accept-language":       "en-US,en;q=0.9",
	"content-type":          "application/json",
	"sec-ch-ua":             "\"Chromium\";v=\"122\", \"Not(A:Brand\";v=\"24\", \"Microsoft Edge\";v=\"122\"",
	"sec-ch-ua-mobile":      "?0",
	"sec-ch-ua-platform":    "\"Windows\"",
	"sec-fetch-dest":        "empty",
	"sec-fetch-mode":        "cors",
	"sec-fetch-site":        "same-site",
	"x-wealthsimple-client": "@wealthsimple/wealthsimple",
}

type WealthSimpleLogin struct {
	Email    string
	Password string
	OtpClaim string
}

func (ws *WealthSimpleLogin) InitOtpClaim() error {
	req, err := ws.getBaseLoginRequest()
	if err != nil {
		return err
	}

	resp, err := req.MakeRequest()
	if err != nil {
		return err
	}

	ws.OtpClaim = resp.GetHeader(otpClaimHeader)

	if ws.OtpClaim == "" {
		return fmt.Errorf("something happened... could not get the otp claim")
	}

	return nil
}

func (ws *WealthSimpleLogin) LoginWithOtp(otp string) (*db.AuthToken, error) {
	req, err := ws.getBaseLoginRequest()
	if err != nil {
		return nil, err
	}

	req.AddHeader("x-wealthsimple-otp", otp)
	req.AddHeader("x-wealthsimple-otp-authenticated-claim", ws.OtpClaim)

	resp, err := req.MakeRequest()
	if err != nil {
		return nil, err
	}

	respBody, err := resp.GetJsonResponse()
	if err != nil {
		return nil, err
	}

	fmt.Println(respBody)

	return &db.AuthToken{
		AccessToken:  respBody["access_token"].(string),
		RefreshToken: respBody["refresh_token"].(string),
		CreatedAt:    createdAtToString(respBody["created_at"].(float64)),
	}, nil
}

func (ws *WealthSimpleLogin) getBaseLoginRequest() (*Request, error) {
	jsonForm, err := json.Marshal(map[string]string{
		"grant_type":     "password",
		"username":       ws.Email,
		"password":       ws.Password,
		"skip_provision": "true",
		"otp_claim":      "null",
		"scope":          "invest.read invest.write trade.read trade.write tax.read tax.write",
		"client_id":      env.GetEnv("CLIENT_ID"),
	})

	if err != nil {
		return nil, err
	}

	req, err := NewRequest(http.MethodPost, authTokenRoute, jsonForm)
	if err != nil {
		return nil, err
	}

	req.AddHeaderGroup(baseHeaders)

	return req, nil
}

func createdAtToString(createdAt float64) string {
	return strconv.FormatFloat(createdAt, 'f', -1, 64)
}
