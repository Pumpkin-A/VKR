package client

import (
	"fmt"
	"net/http"
	"notifications/config"
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

func (c *Client) Notificate() {
	fmt.Println("синхронный запрос на отравку уведомления")
}
