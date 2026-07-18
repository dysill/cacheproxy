package cache

import (
	"testing"
	"time"
)

// Basic Cache
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

func TestBasicCache_Delete(t *testing.T) {
	var c Cache = NewBasicCache()
	c.Set("testkey", []byte("testval"), time.Minute)

	c.Delete("testkey")

	_, ok := c.Get("testkey")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}
}

func TestBasicCache_Clear(t *testing.T) {
	var c Cache = NewBasicCache()
	c.Set("testkey1", []byte("testval1"), time.Minute)
	c.Set("testkey2", []byte("testval2"), time.Minute)

	c.Clear()

	_, ok := c.Get("testkey1")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}

	_, ok = c.Get("testkey2")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}
}

// FIFO Cache
func TestFIFOCache_SetGet(t *testing.T) {
	var c Cache = NewFIFOCache(100)
	c.Set("testkey", []byte("testval"), time.Minute)

	val, ok := c.Get("testkey")
	if !ok {
		t.Fatalf("expected key to be found, but it was not")
	}
	if string(val) != "testval" {
		t.Errorf("expected value %q, got %q", "testval", val)
	}
}

func TestFIFOCache_Expiry(t *testing.T) {
	var c Cache = NewFIFOCache(100)
	c.Set("testkey", []byte("testval"), 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	_, found := c.Get("testkey")
	if found {
		t.Errorf("expected key to be expired, but it was still found")
	}
}

func TestFIFOCache_Delete(t *testing.T) {
	var c Cache = NewFIFOCache(100)
	c.Set("testkey", []byte("testval"), time.Minute)

	c.Delete("testkey")

	_, ok := c.Get("testkey")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}
}

func TestFIFOCache_Clear(t *testing.T) {
	var c Cache = NewFIFOCache(100)
	c.Set("testkey1", []byte("testval1"), time.Minute)
	c.Set("testkey2", []byte("testval2"), time.Minute)

	c.Clear()

	_, ok := c.Get("testkey1")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}

	_, ok = c.Get("testkey2")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}
}

func TestFIFOCache_Eviction(t *testing.T) {
	var c Cache = NewFIFOCache(2)

	c.Set("testkey1", []byte("testval1"), time.Minute)
	c.Set("testkey2", []byte("testval2"), time.Minute)
	c.Set("testkey3", []byte("testval3"), time.Minute)

	if _, ok := c.Get("testkey1"); ok {
		t.Errorf("expected testkey1 to be evicted")
	}

	if _, ok := c.Get("testkey2"); !ok {
		t.Errorf("expected testkey2 to be found")
	}

	if _, ok := c.Get("testkey3"); !ok {
		t.Errorf("expected testkey3 to be found")
	}
}

func TestFIFOCache_UpdateExisting(t *testing.T) {
	var c Cache = NewFIFOCache(2)

	c.Set("testkey1", []byte("testval1"), time.Minute)
	c.Set("testkey2", []byte("testval2"), time.Minute)
	c.Set("testkey1", []byte("new_testval1"), time.Minute)

	if _, ok := c.Get("testkey2"); !ok {
		t.Errorf("expected testkey2 to be found")
	}

	val, _ := c.Get("testkey1")
	if string(val) != "new_testval1" {
		t.Errorf("expected upated value, got %q", val)
	}
}

// LRU Cache
func TestLRUCache_SetGet(t *testing.T) {
	var c Cache = NewLRUCache(100)
	c.Set("testkey", []byte("testval"), time.Minute)

	val, ok := c.Get("testkey")
	if !ok {
		t.Fatalf("expected key to be found, but it was not")
	}
	if string(val) != "testval" {
		t.Errorf("expected value %q, got %q", "testval", val)
	}
}

func TestLRUCache_Expiry(t *testing.T) {
	var c Cache = NewLRUCache(100)
	c.Set("testkey", []byte("testval"), 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	_, found := c.Get("testkey")
	if found {
		t.Errorf("expected key to be expired, but it was still found")
	}
}

func TestLRUCache_Delete(t *testing.T) {
	var c Cache = NewLRUCache(100)
	c.Set("testkey", []byte("testval"), time.Minute)

	c.Delete("testkey")

	_, ok := c.Get("testkey")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}
}

func TestLRUCache_Clear(t *testing.T) {
	var c Cache = NewLRUCache(100)
	c.Set("testkey1", []byte("testval1"), time.Minute)
	c.Set("testkey2", []byte("testval2"), time.Minute)

	c.Clear()

	_, ok := c.Get("testkey1")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}

	_, ok = c.Get("testkey2")
	if ok {
		t.Errorf("expected key to not be found, but it was")
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	var c Cache = NewLRUCache(2)

	c.Set("testkey1", []byte("testval1"), time.Minute)
	c.Set("testkey2", []byte("testval2"), time.Minute)
	_, _ = c.Get("testkey1")
	c.Set("testkey3", []byte("testval3"), time.Minute)

	if _, ok := c.Get("testkey2"); ok {
		t.Errorf("expected testkey2 to be evicted")
	}

	if _, ok := c.Get("testkey1"); !ok {
		t.Errorf("expected testkey1 to be found")
	}

	if _, ok := c.Get("testkey3"); !ok {
		t.Errorf("expected testkey3 to be found")
	}
}

func TestLRUCache_UpdateExisting(t *testing.T) {
	var c Cache = NewLRUCache(2)

	c.Set("testkey1", []byte("testval1"), time.Minute)
	c.Set("testkey2", []byte("testval2"), time.Minute)
	c.Set("testkey1", []byte("new_testval1"), time.Minute)

	if _, ok := c.Get("testkey2"); !ok {
		t.Errorf("expected testkey2 to be found")
	}

	val, _ := c.Get("testkey1")
	if string(val) != "new_testval1" {
		t.Errorf("expected upated value, got %q", val)
	}
}
