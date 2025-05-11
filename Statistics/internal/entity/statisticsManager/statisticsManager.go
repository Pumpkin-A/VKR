package statisticsManager

import (
	"context"
)

type Client interface {
	Analize()
}

type statisticsManager struct {
	client Client
}

func New(c Client) *statisticsManager {
	return &statisticsManager{
		client: c,
	}
}

func (pm *statisticsManager) Analize(ctx context.Context) {
	pm.client.Analize()

}
