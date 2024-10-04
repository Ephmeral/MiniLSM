package memtable

import (
	"github.com/emirpasic/gods/maps/treemap"
	"math/rand"
	"testing"
)

// BenchmarkSkiplistInsert 对比跳表的插入性能
func BenchmarkSkiplistInsert(b *testing.B) {
	s := NewSkipList()
	for i := 0; i < b.N; i++ {
		key := []byte{byte(rand.Int())}
		value := []byte{byte(rand.Int())}
		s.Put(key, value)
	}
}

// BenchmarkMapInsert 对比map的插入性能
func BenchmarkMapInsert(b *testing.B) {
	treeMap := treemap.NewWithIntComparator()
	for i := 0; i < b.N; i++ {
		key := rand.Int()
		value := rand.Int()
		treeMap.Put(key, value)
	}
}

// BenchmarkSkiplistSearch 对比跳表的搜索性能
func BenchmarkSkiplistSearch(b *testing.B) {
	s := NewSkipList()
	// 插入大量数据
	for i := 0; i < 100000; i++ {
		key := []byte{byte(rand.Int())}
		value := []byte{byte(rand.Int())}
		s.Put(key, value)
	}

	// 开始基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte{byte(i % 100000)} // 随机搜索已插入的键
		s.Search(key)
	}
}

// BenchmarkMapSearch 对比map的搜索性能
func BenchmarkMapSearch(b *testing.B) {
	// TreeMap: 键有序
	treeMap := treemap.NewWithIntComparator()

	// 插入大量数据
	for i := 0; i < 100000; i++ {
		key := rand.Int()
		value := rand.Int()
		treeMap.Put(key, value)
	}
	// 开始基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := i % 100000 // 随机搜索已插入的键
		_, _ = treeMap.Get(key)
	}
}

// BenchmarkSkiplistErase 对比跳表的删除性能
func BenchmarkSkiplistErase(b *testing.B) {
	s := NewSkipList()
	// 插入大量数据
	for i := 0; i < 100000; i++ {
		key := []byte{byte(rand.Int())}
		value := []byte{byte(rand.Int())}
		s.Put(key, value)
	}

	// 开始基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := []byte{byte(i % 100000)} // 随机删除已插入的键
		s.Erase(key)
	}
}

// BenchmarkMapErase 对比map的删除性能
func BenchmarkMapErase(b *testing.B) {
	// TreeMap: 键有序
	treeMap := treemap.NewWithIntComparator()
	// 插入大量数据
	for i := 0; i < 100000; i++ {
		key := rand.Int()
		value := rand.Int()
		treeMap.Put(key, value)
	}

	// 开始基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := i % 100000 // 随机删除已插入的键
		treeMap.Remove(key)
	}
}
