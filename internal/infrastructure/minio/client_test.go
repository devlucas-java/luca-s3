package minio

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestMinIOClient(t *testing.T) {
	// Skip if no MinIO available
	client, err := New("localhost:9000", "username", "password", false)
	if err != nil {
		t.Skip("MinIO not available:", err)
	}

	ctx := context.Background()

	t.Run("UploadAndDownload", func(t *testing.T) {
		content := "test content"
		objectName := "test/upload.txt"

		// Upload
		err := client.UploadObject(ctx, objectName, strings.NewReader(content), int64(len(content)), "text/plain")
		if err != nil {
			t.Fatalf("upload failed: %v", err)
		}

		// Check exists
		exists, err := client.ObjectExists(ctx, objectName)
		if err != nil {
			t.Fatalf("check exists failed: %v", err)
		}
		if !exists {
			t.Fatal("object should exist")
		}

		// Delete
		err = client.DeleteObject(ctx, objectName)
		if err != nil {
			t.Fatalf("delete failed: %v", err)
		}

		// Check not exists
		exists, err = client.ObjectExists(ctx, objectName)
		if err != nil {
			t.Fatalf("check exists failed: %v", err)
		}
		if exists {
			t.Fatal("object should not exist")
		}
	})

	t.Run("PresignedURL", func(t *testing.T) {
		content := "test content"
		objectName := "test/presigned.txt"

		// Upload
		err := client.UploadObject(ctx, objectName, strings.NewReader(content), int64(len(content)), "text/plain")
		if err != nil {
			t.Fatalf("upload failed: %v", err)
		}
		defer client.DeleteObject(ctx, objectName)

		// Get presigned URL
		url, err := client.PresignedURL(ctx, objectName, 1*time.Hour)
		if err != nil {
			t.Fatalf("presigned url failed: %v", err)
		}

		if url == "" {
			t.Fatal("presigned url should not be empty")
		}
	})

	t.Run("ListObjects", func(t *testing.T) {
		prefix := "test/list/"

		// Upload multiple objects
		for i := 0; i < 3; i++ {
			objectName := prefix + string(rune('a'+i)) + ".txt"
			err := client.UploadObject(ctx, objectName, strings.NewReader("test"), 4, "text/plain")
			if err != nil {
				t.Fatalf("upload failed: %v", err)
			}
			defer client.DeleteObject(ctx, objectName)
		}

		// List objects
		objects := client.ListObjects(ctx, prefix)
		if len(objects) < 3 {
			t.Fatalf("expected at least 3 objects, got %d", len(objects))
		}
	})

	t.Run("DeleteFolder", func(t *testing.T) {
		prefix := "test/folder/"

		// Upload multiple objects
		for i := 0; i < 3; i++ {
			objectName := prefix + string(rune('a'+i)) + ".txt"
			err := client.UploadObject(ctx, objectName, strings.NewReader("test"), 4, "text/plain")
			if err != nil {
				t.Fatalf("upload failed: %v", err)
			}
		}

		// Delete folder
		err := client.DeleteFolder(ctx, prefix)
		if err != nil {
			t.Fatalf("delete folder failed: %v", err)
		}

		// Check objects deleted
		objects := client.ListObjects(ctx, prefix)
		if len(objects) > 0 {
			t.Fatalf("expected 0 objects, got %d", len(objects))
		}
	})
}
