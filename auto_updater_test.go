package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizeVersionForFilename(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "v1.2.3", want: "v1_2_3"},
		{input: "  v1.2.3-beta+build  ", want: "v1_2_3_beta_build"},
		{input: "", want: "unknown"},
		{input: "###", want: "unknown"},
	}

	for _, tc := range tests {
		got := sanitizeVersionForFilename(tc.input)
		if got != tc.want {
			t.Fatalf("sanitizeVersionForFilename(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestCalculateDownloadPercent(t *testing.T) {
	if got := calculateDownloadPercent(0, 0); got != 50 {
		t.Fatalf("expected percent 50 when total is unknown, got %d", got)
	}
	if got := calculateDownloadPercent(0, 100); got != 5 {
		t.Fatalf("expected percent 5 at start, got %d", got)
	}
	if got := calculateDownloadPercent(100, 100); got != 85 {
		t.Fatalf("expected percent capped to 85 at download completion, got %d", got)
	}
}

func TestVerifyFileSHA256(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "asset.exe")
	content := []byte("nooltools update package")

	if err := os.WriteFile(filePath, content, 0o600); err != nil {
		t.Fatalf("failed to write temp test file: %v", err)
	}

	sum := sha256.Sum256(content)
	expected := hex.EncodeToString(sum[:])

	if err := verifyFileSHA256(filePath, expected); err != nil {
		t.Fatalf("verifyFileSHA256 should pass with matching digest: %v", err)
	}

	if err := verifyFileSHA256(filePath, "0000000000000000000000000000000000000000000000000000000000000000"); err == nil {
		t.Fatalf("verifyFileSHA256 should fail with mismatched digest")
	}
}
