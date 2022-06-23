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

func printPartsContents(t *testing.T, parts []Part) {
	var out string
	for _, part := range parts {
		out = fmt.Sprintf("%s\n%s", out, string(part.Content))
	}
	fmt.Println(out)
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
	printPartsContents(t, parts)
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
