package load_test

import "testing"

func TestLoad(t *testing.T) {
	err := PerformLoadTesting()
	if err != nil {
		t.Fatalf("Load test failed: %v", err)
	}
}
