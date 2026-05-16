package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
}

func New(addr, password string, db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis: ping: %w", err)
	}

	return &Client{client: rdb}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Client() *redis.Client {
	return c.client
}
