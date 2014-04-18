package main

import (
	"github.com/satori/go.uuid"
	"testing"
)

func TestISUUID(t *testing.T) {
	for i := 0; i < 100; i++ {
		s := uuid.NewV4().String()
		if !IsUUID(s) {
			t.Fatalf("\"%s\" must be an uuid.", s)
		}
	}
}
