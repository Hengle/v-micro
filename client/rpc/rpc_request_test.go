package rpc

import "testing"

func TestRequestOptions(t *testing.T) {
	r := newRequest("1", "service", "m1", nil, "application/protobuf")
	if r.Service() != "service" {
		t.Fatalf("expected 'service' got %s", r.Service())
	}
	if r.Method() != "m1" {
		t.Fatalf("expected 'm1' got %s", r.Method())
	}
	if r.ContentType() != "application/protobuf" {
		t.Fatalf("expected 'application/protobuf' got %s", r.ContentType())
	}
	if r.ID() != "1" {
		t.Fatalf("expected '1' got %s", r.ID())
	}
}
