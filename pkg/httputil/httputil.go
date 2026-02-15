// Package httputil fornece utilitários para operações HTTP.
package httputil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client é um wrapper em torno do http.Client com configurações padrão
type Client struct {
	client *http.Client
}

// NewClient cria um novo HTTP client com timeout
func NewClient(timeout time.Duration) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// DefaultClient retorna um client com timeout de 10 segundos
func DefaultClient() *Client {
	return NewClient(10 * time.Second)
}

// Get realiza uma requisição GET
func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status HTTP inesperado: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	return body, nil
}

// GetWithHeaders realiza uma requisição GET com headers customizados
func (c *Client) GetWithHeaders(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status HTTP inesperado: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %w", err)
	}

	return body, nil
}

// DownloadFile baixa um arquivo da URL e retorna os bytes
func (c *Client) DownloadFile(ctx context.Context, url string) ([]byte, error) {
	return c.Get(ctx, url)
}

// IsReachable verifica se uma URL está acessível
func (c *Client) IsReachable(ctx context.Context, url string) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}
