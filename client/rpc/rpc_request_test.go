package rpc

import "testing"

func TestRequestOptions(t *testing.T) {
	r := newRequest("service", "m1", nil, "application/protobuf")
	if r.Service() != "service" {
		t.Fatalf("expected 'service' got %s", r.Service())
	}
	if r.Method() != "m1" {
		t.Fatalf("expected 'm1' got %s", r.Method())
	}
	if r.ContentType() != "application/protobuf" {
		t.Fatalf("expected 'm1' got %s", r.ContentType())
	}
}
