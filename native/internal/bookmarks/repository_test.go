package bookmarks

import (
	"context"
	"errors"
	"testing"
)

func TestServiceRejectsMissingOwnerIdentity(t *testing.T) {
	service, err := NewService(NewMemoryRepository(nil))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.List(context.Background(), "   "); !errors.Is(err, ErrOwnerIdentityRequired) {
		t.Fatalf("expected owner identity failure, got %v", err)
	}
}

func TestServicePreservesOwnerIsolation(t *testing.T) {
	service, err := NewService(NewMemoryRepository(nil))
	if err != nil {
		t.Fatal(err)
	}
	created, err := service.Create(context.Background(), "owner-a", CreateInput{URL: "https://example.com/private", Title: "Private"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok, err := service.Get(context.Background(), "owner-b", created.ID); err != nil || ok {
		t.Fatalf("cross-owner read must fail closed: ok=%v err=%v", ok, err)
	}
	item, ok, err := service.Get(context.Background(), "owner-a", created.ID)
	if err != nil || !ok || item.OwnerID != "owner-a" {
		t.Fatalf("owner read failed: %#v ok=%v err=%v", item, ok, err)
	}
}

func TestMemoryRepositorySatisfiesRepositoryContract(t *testing.T) {
	var repository Repository = NewMemoryRepository(nil)
	service, err := NewService(repository)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Create(context.Background(), "owner", CreateInput{URL: "https://example.org"}); err != nil {
		t.Fatal(err)
	}
	items, err := service.List(context.Background(), "owner")
	if err != nil || len(items) != 1 {
		t.Fatalf("unexpected list: %#v err=%v", items, err)
	}
}
