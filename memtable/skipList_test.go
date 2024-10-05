package memtable

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	skipList := NewSkipList()
	skipList.Put([]byte{1}, []byte{1})
	skipList.Put([]byte{2}, []byte{2})
	skipList.Put([]byte{3}, []byte{3})
	value, ok := skipList.Get([]byte{0})
	fmt.Println("search 0 = ", value, " ok = ", ok)
	skipList.Put([]byte{4}, []byte{4})
	value, ok = skipList.Get([]byte{1})
	fmt.Println("search 1 = ", value, " ok = ", ok)
	skipList.Erase([]byte{0})
	skipList.Erase([]byte{1})
	value, ok = skipList.Get([]byte{1})
	fmt.Println("search 1 = ", value, " ok = ", ok)
}

func TestEmptySkiplistSearch(t *testing.T) {
	s := NewSkipList()
	key := []byte("key1")
	if _, found := s.Get(key); found {
		t.Errorf("Expected key %s to not be found in empty skiplist", key)
	}
}

func TestSkiplistInsertAndSearch(t *testing.T) {
	s := NewSkipList()
	key := []byte("key1")
	value := []byte("value1")

	s.Put(key, value)
	got, found := s.Get(key)

	if !found {
		t.Errorf("Expected key %s to be found", key)
	}

	if !bytes.Equal(got, value) {
		t.Errorf("Expected value %s, but got %s", value, got)
	}
}

func TestSkiplistOverwriteValue(t *testing.T) {
	s := NewSkipList()
	key := []byte("key1")
	value1 := []byte("value1")
	value2 := []byte("value2")

	s.Put(key, value1)
	s.Put(key, value2)

	got, found := s.Get(key)

	if !found {
		t.Errorf("Expected key %s to be found", key)
	}

	if !bytes.Equal(got, value2) {
		t.Errorf("Expected value %s, but got %s", value2, got)
	}
}

func TestSkiplistErase(t *testing.T) {
	s := NewSkipList()
	key1 := []byte("key1")
	value1 := []byte("value1")
	key2 := []byte("key2")
	value2 := []byte("value2")

	// 插入两个键值对
	s.Put(key1, value1)
	s.Put(key2, value2)

	// 删除其中一个
	s.Erase(key1)

	// key1 应该被删除，找不到
	if _, found := s.Get(key1); found {
		t.Errorf("Expected key %s to be erased", key1)
	}

	// key2 应该仍然存在
	if got, found := s.Get(key2); !found || !bytes.Equal(got, value2) {
		t.Errorf("Expected key %s to be found with value %s", key2, value2)
	}
}

func TestSkiplistEraseNonExistentKey(t *testing.T) {
	s := NewSkipList()
	key := []byte("nonexistent")

	// 删除一个不存在的键
	if erased := s.Erase(key); erased {
		t.Errorf("Expected key %s to not be erased", key)
	}
}

func TestSkiplistRandomInsertErase(t *testing.T) {
	s := NewSkipList()
	keys := [][]byte{
		[]byte("a"), []byte("b"), []byte("c"),
	}
	values := [][]byte{
		[]byte("value1"), []byte("value2"), []byte("value3"),
	}

	// 插入多个键值对
	for i := range keys {
		s.Put(keys[i], values[i])
	}

	// 删除部分键
	if erased := s.Erase(keys[1]); !erased {
		t.Errorf("Expected key %s to be erased", keys[1])
	}

	// 剩余键应该依然可以找到
	for i, key := range keys {
		if i == 1 {
			continue // 跳过已删除的键
		}
		if got, found := s.Get(key); !found || !bytes.Equal(got, values[i]) {
			t.Errorf("Expected key %s to be found with value %s", key, values[i])
		}
	}
}
