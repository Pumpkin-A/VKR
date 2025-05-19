package client

import (
	"billing/config"
	"billing/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	tracer     trace.Tracer
}

func New(cfg config.Config) *Client {
	return &Client{
		baseURL: cfg.Client.BaseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10, // Устанавливаем таймаут для запросов
		},
		tracer: otel.Tracer("BankClient_billing"),
	}
}

// DoPayment - метод для отправки запроса на проведение платежа
func (c *Client) DoPayment(ctx context.Context, payment models.Payment) (models.ResultOfRequestFromBank, error) {
	_, sp := c.tracer.Start(ctx, "BankClient.DoPayment")
	sp.SetAttributes(attribute.String("paymentId", payment.UUID))
	defer sp.End()

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
