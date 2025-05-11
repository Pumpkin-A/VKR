package notificationsManager

import (
	"context"
)

type Client interface {
	Notificate()
}

type notificationsManager struct {
	client Client
}

func New(c Client) *notificationsManager {
	return &notificationsManager{
		client: c,
	}
}

func (pm *notificationsManager) Notificate(ctx context.Context) {
	pm.client.Notificate()

}
