package client

import (
	"fmt"
	"net/http"
	"statistics/config"
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
			Timeout: time.Second * 10,
		},
	}
}

func (c *Client) Analize() {
	fmt.Println("синхронный запрос на отправку данных статистике")
}
