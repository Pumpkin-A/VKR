package client

import (
	"billing/config"
	"billing/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(cfg config.Config) *Client {
	return &Client{
		baseURL: cfg.Client.BaseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10, // Устанавливаем таймаут для запросов
		},
	}
}

// DoPayment - метод для отправки запроса на проведение платежа
func (c *Client) DoPayment(payment models.Payment) (models.ResultOfRequestFromBank, error) {
	requestBody, err := json.Marshal(payment)
	if err != nil {
		return models.ResultOfRequestFromBank{}, fmt.Errorf("failed to marshal payment: %w", err)
	}

	requestURL := c.baseURL + "/dopayment" // Формируем URL запроса

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return models.ResultOfRequestFromBank{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.ResultOfRequestFromBank{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Читаем тело ошибки, если есть
		return models.ResultOfRequestFromBank{}, fmt.Errorf("server returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result models.ResultOfRequestFromBank
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return models.ResultOfRequestFromBank{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
