//go:build go1.18
// +build go1.18

package trans

import "testing"

func TestPtr(t *testing.T) {
	b := true
	pb := Ptr(b)
	if pb == nil {
		t.Fatal("unexpected nil conversion")
	}
	if *pb != b {
		t.Fatalf("got %v, want %v", *pb, b)
	}
}

func TestClonePtr(t *testing.T) {
	original := Ptr("value")
	cloned := ClonePtr(original)
	if cloned == nil {
		t.Fatal("unexpected nil clone")
	}
	if cloned == original {
		t.Fatal("expected a new pointer")
	}
	if *cloned != *original {
		t.Fatalf("got %q, want %q", *cloned, *original)
	}

	*cloned = "changed"
	if *original != "value" {
		t.Fatalf("clone mutation changed original to %q", *original)
	}

	if ClonePtr[string](nil) != nil {
		t.Fatal("expected nil clone for nil input")
	}
}

func TestSliceOfPtrs(t *testing.T) {
	arr := SliceOfPtrs[int]()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = SliceOfPtrs(1, 2, 3, 4, 5)
	for i, v := range arr {
		if *v != i+1 {
			t.Fatal("values don't match")
		}
	}
}
