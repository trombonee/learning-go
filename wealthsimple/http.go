package wealthsimple

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Request struct {
	req *http.Request
}

func (r *Request) AddHeaderGroup(headers map[string]string) {
	for key, value := range headers {
		r.req.Header.Add(key, value)
	}
}

func (r *Request) AddHeader(key string, value string) {
	r.req.Header.Add(key, value)
}

func (r *Request) MakeRequest() (*Response, error) {
	var err error

	client := &http.Client{}
	resp, err := client.Do(r.req)
	if err != nil {
		return nil, err
	}

	return &Response{resp: resp}, nil
}

func NewRequest(method string, route string, data []byte) (*Request, error) {
	request, err := http.NewRequest(method, route, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return &Request{req: request}, nil
}

type Response struct {
	resp *http.Response
}

func (r *Response) GetHeader(key string) string {
	return r.resp.Header.Get(key)
}

func (r *Response) GetBody() ([]byte, error) {
	return io.ReadAll(r.resp.Body)
}

func (r *Response) GetJsonResponse() (map[string]interface{}, error) {
	var responseBody map[string]interface{}
	body, err := r.GetBody()
	if err != nil {
		return responseBody, err
	}

	err = json.Unmarshal(body, &responseBody)

	return responseBody, err
}

func (r *Response) Close() {
	r.resp.Body.Close()
}
