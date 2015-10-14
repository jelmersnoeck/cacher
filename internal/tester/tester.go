package tester

import (
	"encoding/binary"
	"strconv"
	"testing"

	"github.com/jelmersnoeck/cacher"
	"github.com/jelmersnoeck/cacher/internal/numbers"
)

func RunCacher(t *testing.T, cacher cacher.Cacher) {
	TestSet(t, cacher)
	cacher.Flush()

	TestSetMulti(t, cacher)
	cacher.Flush()

	TestAdd(t, cacher)
	cacher.Flush()

	TestReplace(t, cacher)
	cacher.Flush()

	TestDecrement(t, cacher)
	cacher.Flush()

	TestIncrement(t, cacher)
	cacher.Flush()

	TestGet(t, cacher)
	cacher.Flush()

	TestGetMulti(t, cacher)
	cacher.Flush()

	TestDelete(t, cacher)
	cacher.Flush()

	TestDeleteMulti(t, cacher)
	cacher.Flush()

	TestFlush(t, cacher)
	cacher.Flush()
}

func TestAdd(t *testing.T, cache cacher.Cacher) {
	if !cache.Add("key1", []byte("value1"), 0) {
		t.Errorf("Expecting `key1` to be added to the cache")
		t.FailNow()
	}

	if cache.Add("key1", []byte("value2"), 0) {
		t.Errorf("Expecting `key1` not to be added to the cache")
		t.FailNow()
	}

	Compare(t, cache, "key1", "value1")
}

func TestSet(t *testing.T, cache cacher.Cacher) {
	if !cache.Set("key1", []byte("value"), 0) {
		t.Errorf("Expecting `key1` to be `value`")
		t.Fail()
	}

	if !cache.Set("key2", numbers.Int64Bytes(2), 0) {
		t.Errorf("Expecting `key2` to be `2`")
		t.Fail()
	}
}

func TestSetMulti(t *testing.T, cache cacher.Cacher) {
	items := map[string][]byte{
		"item1": numbers.Int64Bytes(1),
		"item2": []byte("string"),
	}

	cache.SetMulti(items, 0)

	Compare(t, cache, "item1", 1)
	Compare(t, cache, "item2", "string")
}

func TestReplace(t *testing.T, cache cacher.Cacher) {
	if cache.Replace("key1", []byte("value1"), 0) {
		t.Errorf("Key1 is not set yet, should not be able to replace.")
		t.FailNow()
	}

	cache.Set("key1", []byte("value1"), 0)
	if !cache.Replace("key1", []byte("value1"), 0) {
		t.Errorf("Key1 has been set, should be able to replace.")
		t.FailNow()
	}
}

func TestIncrement(t *testing.T, cache cacher.Cacher) {
	cache.Increment("key1", 0, 1, 0)
	Compare(t, cache, "key1", 0)

	cache.Increment("key1", 0, 1, 0)
	Compare(t, cache, "key1", 1)

	cache.Set("string", []byte("value"), 0)
	if cache.Increment("string", 0, 1, 0) {
		t.Errorf("Can't increment a string value.")
		t.FailNow()
	}

	if cache.Increment("key2", 0, 0, 0) {
		t.Errorf("Can't have an offset of <= 0")
		t.FailNow()
	}

	if cache.Increment("key3", -1, 1, 0) {
		t.Errorf("Can't have an initial value of < 0")
		t.FailNow()
	}
}

func TestDecrement(t *testing.T, cache cacher.Cacher) {
	cache.Decrement("key1", 10, 1, 0)
	Compare(t, cache, "key1", 10)

	cache.Decrement("key1", 10, 1, 0)
	Compare(t, cache, "key1", 9)

	cache.Set("string", []byte("value"), 0)
	if cache.Decrement("string", 0, 1, 0) {
		t.Errorf("Can't decrement a string value.")
		t.FailNow()
	}

	if cache.Decrement("key2", 0, 0, 0) {
		t.Errorf("Can't have an offset of <= 0")
		t.FailNow()
	}

	if cache.Decrement("key3", -1, 1, 0) {
		t.Errorf("Can't have an initial value of < 0")
		t.FailNow()
	}

	if cache.Decrement("key1", 10, 10, 0) {
		t.Errorf("Can't decrement below 0")
		t.FailNow()
	}
}

func TestGet(t *testing.T, cache cacher.Cacher) {
	cache.Set("key1", []byte("value1"), 0)
	Compare(t, cache, "key1", "value1")

	if _, ok := cache.Get("key2"); ok {
		t.Errorf("Key2 is not present, ok should be false.")
		t.FailNow()
	}
}

func TestGetMulti(t *testing.T, cache cacher.Cacher) {
	items := map[string][]byte{
		"item1": numbers.Int64Bytes(1),
		"item2": []byte("string"),
	}

	cache.SetMulti(items, 0)

	var keys []string
	for k, _ := range items {
		keys = append(keys, k)
	}

	values := cache.GetMulti(keys)

	_, val := binary.Varint(values["item1"])
	if val != 1 {
		t.Errorf("Expected `item1` to equal `1`")
		t.FailNow()
	}

	if string(values["item2"]) != "string" {
		t.Errorf("Expected `item2` to equal `string`")
		t.FailNow()
	}
}

func TestDelete(t *testing.T, cache cacher.Cacher) {
	cache.Set("key1", []byte("value1"), 0)
	Compare(t, cache, "key1", "value1")

	cache.Delete("key1")

	if _, ok := cache.Get("key1"); ok {
		t.Errorf("`key1` should be deleted from the cache.")
		t.FailNow()
	}
}

func TestDeleteMulti(t *testing.T, cache cacher.Cacher) {
	items := map[string][]byte{
		"item1": numbers.Int64Bytes(1),
		"item2": []byte("string"),
	}

	cache.SetMulti(items, 0)
	cache.Set("key1", []byte("value1"), 0)

	var keys []string
	for k, _ := range items {
		keys = append(keys, k)
	}

	cache.DeleteMulti(keys)

	if _, ok := cache.Get("item1"); ok {
		t.Errorf("`item1` should be deleted from the cache.")
		t.FailNow()
	}

	if _, ok := cache.Get("item2"); ok {
		t.Errorf("`item2` should be deleted from the cache.")
		t.FailNow()
	}

	Compare(t, cache, "key1", "value1")
}

func TestFlush(t *testing.T, cache cacher.Cacher) {
	cache.Set("key1", []byte("value1"), 0)
	Compare(t, cache, "key1", "value1")

	if !cache.Flush() {
		t.Fail()
	}

	if v, _ := cache.Get("key1"); v != nil {
		t.Errorf("Expecting `key1` to be nil")
		t.Fail()
	}
}

func Compare(t *testing.T, cache cacher.Cacher, key string, value interface{}) {
	_, ok := value.(int)
	if ok {
		val := int64(value.(int))
		v, _ := cache.Get(key)
		valInt, _ := numbers.BytesInt64(v)
		if valInt != val {
			t.Errorf("Expected `" + key + "` to equal `" + strconv.FormatInt(val, 10) + "`, is `" + strconv.FormatInt(valInt, 10) + "`")
			t.FailNow()
		}
	} else {
		value = value.(string)
		if v, _ := cache.Get(key); string(v) != value {
			t.Errorf("Expected `" + key + "` to equal `" + value.(string) + "`")
			t.FailNow()
		}
	}
}

func NotPresent(t *testing.T, cache cacher.Cacher, key string) {
	if _, ok := cache.Get(key); ok {
		t.Errorf("Expected `" + key + "` not to be present")
		t.FailNow()
	}
}
