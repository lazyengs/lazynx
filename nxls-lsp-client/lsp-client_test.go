package lspclient

import (
	"testing"
)

func TestNxlsLspClient(t *testing.T) {
	result := NxlsLspClient("works")
	if result != "NxlsLspClient works" {
		t.Error("Expected NxlsLspClient to append 'works'")
	}
}
