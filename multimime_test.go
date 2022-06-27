package multimime

import (
	"fmt"
	"os"
	"testing"
)

func openTestEmail(t *testing.T) *os.File {
	testEmail := "test.eml"
	r, err := os.Open(testEmail)
	if err != nil {
		t.Fatalf("Couldn't read %s", testEmail)
	}
	return r
}

func TestGetInlineText(t *testing.T) {
	r := openTestEmail(t)
	text, err := GetInlineText(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(text)
}

func TestGetAllParts(t *testing.T) {
	r := openTestEmail(t)
	parts, err := GetAllParts(r)
	if err != nil {
		t.Fatal(err)
	}
	expected := 4
	if len(parts) != expected {
		t.Fatalf("Too many parts: expected %d, got %d", expected, len(parts))
	}
}

func TestGetTextParts(t *testing.T) {
	r := openTestEmail(t)
	parts, err := GetTextParts(r)
	if err != nil {
		t.Fatal(err)
	}
	expected := 3
	if len(parts) != expected {
		t.Fatalf("Too many text parts: expected %d, got %d", expected, len(parts))
	}
}

func TestGetAttachments(t *testing.T) {
	r := openTestEmail(t)
	parts, err := GetAttachments(r)
	if err != nil {
		t.Fatal(err)
	}
	expected := 2
	if len(parts) != expected {
		t.Fatalf("Too many attachments: expected %d, got %d", expected, len(parts))
	}
}
