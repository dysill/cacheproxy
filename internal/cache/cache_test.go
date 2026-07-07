package cache

import (
	"testing"
	"time"
)

func TestBasicCache_SetGet(t *testing.T) {
	var c Cache = NewBasicCache()
	c.Set("testkey", []byte("testval"), time.Minute)

	val, ok := c.Get("testkey")
	if !ok {
		t.Fatalf("expected key to be found, but it was not")
	}
	if string(val) != "testval" {
		t.Errorf("expected value %q, got %q", "testval", val)
	}
}

func TestBasicCache_Expiry(t *testing.T) {
	var c Cache = NewBasicCache()
	c.Set("testkey", []byte("testval"), 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	_, found := c.Get("testkey")
	if found {
		t.Errorf("expected key to be expired, but it was still found")
	}
}
