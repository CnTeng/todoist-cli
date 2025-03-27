package client

import (
	"context"

	"github.com/CnTeng/todoist-api-go/sync/v9"
	"github.com/CnTeng/todoist-cli/internal/db"
	"github.com/CnTeng/todoist-cli/internal/model"
)

type Client struct {
	db *db.DB
	sc *sync.Client
}

func NewClient(db *db.DB, sc *sync.Client) *Client {
	return &Client{
		db: db,
		sc: sc,
	}
}

func (c *Client) ListTasks(ctx context.Context) ([]*model.Task, error) {
	return c.db.ListTasks()
}
