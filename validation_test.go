package main

import (
	"github.com/satori/go.uuid"
	"strings"
	"testing"
)

func TestISUUID(t *testing.T) {
	for i := 0; i < 50; i++ {
		s := uuid.NewV4().String()
		if !IsUUID(s) {
			t.Fatalf("\"%s\" must be an uuid.", s)
		}
	}
	for i := 0; i < 50; i++ {
		s := strings.ToUpper(uuid.NewV4().String())
		if !IsUUID(s) {
			t.Fatalf("\"%s\" must be an uuid.", s)
		}
	}
}

func TestValidationError(t *testing.T) {
	var err ValidationError
	if err != nil {
		t.Fatal("Initial ValidationError should be nil")
	}
	err = err.Append("gallery id should be a uuid")
	err = err.Append("exhibition id shouldn't be empty")
	msg := err.Error()
	contains := []string{
		"Validation Error:",
		"gallery id should be a uuid",
		"exhibition id shouldn't be empty",
	}
	for _, s := range contains {
		if !strings.Contains(msg, s) {
			t.Fatalf("%s should contains \"%s\"", msg, s)
		}
	}
}
