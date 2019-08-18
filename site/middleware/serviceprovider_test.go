package middleware

import (
	"goairmon/site/testhelpers"
	"testing"
)

func TestServiceProvider(t *testing.T) {
	provider := NewServiceProvider()

	provider.Register("first_key", "first_value")
	provider.Register("second_key", "second_value")

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}

	err := provider.BindServices()(testhelpers.EmptyHandler)(ctx)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if ctx.Get("first_key").(string) != "first_value" {
		t.Error("expected bound value")
	}

	if ctx.Get("second_key").(string) != "second_value" {
		t.Error("expected bound value")
	}

}
