package main

import "testing"

func TestListenAddressDefaultsToLoopback(t *testing.T) {
	t.Setenv("GOREECLOUD_BOOKMARKS_LISTEN_ADDR", "")

	if got := listenAddress(); got != defaultListenAddress {
		t.Fatalf("listenAddress() = %q, want %q", got, defaultListenAddress)
	}
}

func TestListenAddressUsesExplicitOverride(t *testing.T) {
	t.Setenv("GOREECLOUD_BOOKMARKS_LISTEN_ADDR", "  127.0.0.1:9090  ")

	if got := listenAddress(); got != "127.0.0.1:9090" {
		t.Fatalf("listenAddress() = %q, want 127.0.0.1:9090", got)
	}
}
