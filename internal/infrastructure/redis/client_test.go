package redis

import (
	"context"
	"testing"
)

func TestRedisClient(t *testing.T) {
	// Skip if no Redis available
	client, err := New("localhost:6379", "", 0)
	if err != nil {
		t.Skip("Redis not available:", err)
	}
	defer client.Close()

	ctx := context.Background()

	t.Run("SetAndGet", func(t *testing.T) {
		key := "test:key"
		value := "test value"

		// Set
		err := client.client.Set(ctx, key, value, 0).Err()
		if err != nil {
			t.Fatalf("set failed: %v", err)
		}

		// Get
		result, err := client.client.Get(ctx, key).Result()
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}

		if result != value {
			t.Fatalf("expected %s, got %s", value, result)
		}

		// Delete
		err = client.client.Del(ctx, key).Err()
		if err != nil {
			t.Fatalf("delete failed: %v", err)
		}
	})

	t.Run("Ping", func(t *testing.T) {
		err := client.client.Ping(ctx).Err()
		if err != nil {
			t.Fatalf("ping failed: %v", err)
		}
	})
}
