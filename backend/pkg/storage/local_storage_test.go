package storage

import (
	"os"
	"testing"
)

func TestUploadFoto(t *testing.T) {
	tmpDir := t.TempDir()
	svc := NewLocalStorageService(tmpDir)

	content := []byte("fake image data")
	path, err := svc.UploadFoto("photo.jpg", content)
	if err != nil {
		t.Fatalf("UploadFoto returned error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("expected file at %s but it does not exist", path)
	}

	saved, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	if string(saved) != string(content) {
		t.Errorf("file content mismatch: got %q, want %q", saved, content)
	}
}

func TestGetFoto(t *testing.T) {
	tmpDir := t.TempDir()
	svc := NewLocalStorageService(tmpDir)

	content := []byte("fake image data")
	path, err := svc.UploadFoto("photo.jpg", content)
	if err != nil {
		t.Fatalf("UploadFoto returned error: %v", err)
	}

	got, err := svc.GetFoto(path)
	if err != nil {
		t.Fatalf("GetFoto returned error: %v", err)
	}

	if string(got) != string(content) {
		t.Errorf("GetFoto content mismatch: got %q, want %q", got, content)
	}
}

func TestGetFoto_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	svc := NewLocalStorageService(tmpDir)

	_, err := svc.GetFoto(tmpDir + "/nonexistent.jpg")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
