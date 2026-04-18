package baron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	ErrParseURL       = "baron: client: parse url error"
	ErrJsonMarshal    = "baron: client: marshal body error"
	ErrCreateRequest  = "baron: client: create request error"
	ErrDoRequest      = "baron: client: do request error"
	ErrDoRequestFails = "baron: client: do request fails"
	ErrDecodeResponse = "baron: client: decoding response body"
)

type BaronClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func New(baseURL, apiKey string, timeout time.Duration) *BaronClient {
	return &BaronClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (b *BaronClient) createRequest(ctx context.Context, method string, endpoint string, body *bytes.Reader) (*http.Request, error) {
	url, err := url.Parse(b.baseURL + "/" + endpoint)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrParseURL, err)
	}

	// var bodyReader *bytes.Reader
	// bodyReader = bytes.NewReader(body)
	// if body != nil {
	// 	bodyBytes, err := json.Marshal(body)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("%s: %w", ErrJsonMarshal, err)
	// 	}
	// 	bodyReader = bytes.NewReader(bodyBytes)
	// } else {
	// 	bodyReader = bytes.NewReader(nil)
	// }

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCreateRequest, err)
	}

	b.setHeaders(req)

	return req, nil
}

func (b *BaronClient) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("apikey", b.apiKey)
}

func (b *BaronClient) do(req *http.Request, out any) error {
	resp, err := b.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s, %w", ErrDoRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s with status code: %d", ErrDoRequestFails, resp.StatusCode)
	}

	if out == nil {
		return nil
	}
	// bodyBytes, _ := io.ReadAll(resp.Body)
	// fmt.Println("response body: ", string(bodyBytes))
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("%s: %w", ErrDecodeResponse, err)
	}

	return nil
}
