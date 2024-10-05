package wal

import (
	"bytes"
	"ephmeral.com/v2/memtable"
	"testing"
)

func TestWALWrite(t *testing.T) {
	writer, err := NewWriter("./test.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer writer.Close()

	skiplist := memtable.NewSkipList()

	kvs := make([]*memtable.KV, 0, 100)
	for i := 0; i < 100; i++ {
		kvs = append(kvs, &memtable.KV{
			Key:   []byte{'a' + uint8(i)},
			Value: []byte{'b' + uint8(i)},
		})
	}
	for _, kv := range kvs {
		skiplist.Put(kv.Key, kv.Value)
		if err = writer.Write(kv.Key, kv.Value); err != nil {
			t.Error(err)
			return
		}
	}

	reader, err := NewReader("./test.txt")
	if err != nil {
		t.Error(err)
		return
	}
	defer reader.Close()

	recoverSkipList := memtable.NewSkipList()
	if reader.Recover(recoverSkipList) != nil {
		return
	}

	originKVs := skiplist.AllData()
	restoredKVs := recoverSkipList.AllData()
	if len(originKVs) != len(restoredKVs) {
		t.Errorf("not euqal len, got: %d, expect: %d", len(restoredKVs), len(originKVs))
		return
	}

	for i := 0; i < len(originKVs); i++ {
		if !bytes.Equal(originKVs[i].Key, restoredKVs[i].Key) {
			t.Errorf("not euqal, index: %d, got key: %s, expect: %s", i, restoredKVs[i].Key, originKVs[i].Key)
		}
		if !bytes.Equal(originKVs[i].Value, restoredKVs[i].Value) {
			t.Errorf("not euqal, index: %d, got val: %s, expect: %s", i, restoredKVs[i].Value, originKVs[i].Value)
		}
	}
}
