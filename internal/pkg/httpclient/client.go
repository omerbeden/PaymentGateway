package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestParam[B any] struct {
	Client       *http.Client
	Header       *http.Header
	Method       string
	URL          string
	Body         B
	Ctx          context.Context
	ClientID     string
	ClientSecret string
}

func MakeRequest[B, T any](p RequestParam[B], result T) error {

	var body io.Reader

	bodyBytes, err := json.Marshal(p.Body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	body = bytes.NewReader(bodyBytes)

	req, err := http.NewRequestWithContext(p.Ctx, p.Method, p.URL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if p.ClientID != "" && p.ClientSecret != "" {
		req.SetBasicAuth(p.ClientID, p.ClientSecret)
	}

	if p.Header != nil {
		for k, vv := range *p.Header {
			for _, v := range vv {
				req.Header.Add(k, v)
			}
		}
	}

	resp, err := p.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil

}
