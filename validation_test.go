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
