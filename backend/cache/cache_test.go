package cache

import (
	"os"
	"testing"
	"time"
)

func TestCache_SetGet(t *testing.T) {
	c := New()

	c.Set("key1", "value1")

	val, found := c.Get("key1")
	if !found {
		t.Fatal("key1 should be found")
	}
	if val != "value1" {
		t.Errorf("value = %v, want 'value1'", val)
	}
}

func TestCache_GetMiss(t *testing.T) {
	c := New()

	_, found := c.Get("nonexistent")
	if found {
		t.Error("nonexistent key should not be found")
	}
}

func TestCache_Delete(t *testing.T) {
	c := New()

	c.Set("key1", "value1")
	c.Delete("key1")

	_, found := c.Get("key1")
	if found {
		t.Error("deleted key should not be found")
	}
}

func TestCache_Flush(t *testing.T) {
	c := New()

	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Flush()

	keys := c.Keys()
	if len(keys) != 0 {
		t.Errorf("after flush, keys = %d, want 0", len(keys))
	}
}

func TestCache_Keys(t *testing.T) {
	c := New()

	c.Set("x", 1)
	c.Set("y", 2)

	keys := c.Keys()
	if len(keys) != 2 {
		t.Errorf("keys count = %d, want 2", len(keys))
	}
}

func TestCache_IsEnabled(t *testing.T) {
	c := New()
	if !c.IsEnabled() {
		t.Error("default cache should be enabled")
	}
}

func TestCache_GetTTL(t *testing.T) {
	c := New()
	expected := 30 * time.Minute
	if c.GetTTL() != expected {
		t.Errorf("TTL = %v, want %v", c.GetTTL(), expected)
	}
}

func TestCache_DisabledWhenTTLZero(t *testing.T) {
	os.Setenv("CACHE_TTL_MINUTES", "0")
	defer os.Unsetenv("CACHE_TTL_MINUTES")

	c := New()

	if c.IsEnabled() {
		t.Error("cache with TTL=0 should be disabled")
	}

	c.Set("key", "value")
	_, found := c.Get("key")
	if found {
		t.Error("disabled cache should not store values")
	}
}

func TestCache_CustomTTL(t *testing.T) {
	os.Setenv("CACHE_TTL_MINUTES", "5")
	defer os.Unsetenv("CACHE_TTL_MINUTES")

	c := New()
	expected := 5 * time.Minute
	if c.GetTTL() != expected {
		t.Errorf("TTL = %v, want %v", c.GetTTL(), expected)
	}
}

func TestCache_ComplexValues(t *testing.T) {
	c := New()

	type TestStruct struct {
		Name  string
		Value int
	}

	c.Set("struct", TestStruct{Name: "test", Value: 42})

	val, found := c.Get("struct")
	if !found {
		t.Fatal("struct should be found")
	}

	ts, ok := val.(TestStruct)
	if !ok {
		t.Fatalf("expected TestStruct, got %T", val)
	}
	if ts.Name != "test" || ts.Value != 42 {
		t.Errorf("struct = %+v", ts)
	}
}
