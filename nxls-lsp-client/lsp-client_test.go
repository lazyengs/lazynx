package lspclient

import (
	"testing"
)

func TestLspClient(t *testing.T) {
	result := LspClient("works")
	if result != "NxlsLspClient works" {
		t.Error("Expected NxlsLspClient to append 'works'")
	}
}
