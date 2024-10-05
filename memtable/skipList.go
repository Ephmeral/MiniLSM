package memtable

import (
	"bytes"
	"math/rand"
	"time"
)

type Node struct {
	next       []*Node
	key, value []byte
}

type Skiplist struct {
	head       *Node // 头节点
	entryCount int   // kv数据对的个数
	memorySize int   // 跳表数据占用内存大小，单位byte
}

func NewSkipList() MemTable {
	head := &Node{[]*Node{nil}, []byte{}, []byte{}}
	return &Skiplist{head, 0, 0}
}

// roll 出一个节点的高度. 最小为 1，每提高 1 层，概率减少为 1//2
func (s *Skiplist) roll() int {
	var level int
	rander := rand.New(rand.NewSource(time.Now().Unix()))
	for rander.Intn(2) == 1 {
		level++
	}
	return level + 1
}

// getNode 内部查找某个key，返回对应的节点指针
func (s *Skiplist) getNode(key []byte) *Node {
	// 从最高层开始查找
	node := s.head
	for level := len(node.next) - 1; level >= 0; level-- {
		// 找到下一个节点不存在或比key大，此时下一个节点的位置
		for node.next[level] != nil && bytes.Compare(node.next[level].key, key) < 0 {
			node = node.next[level]
		}
		// 找到了key对应的节点
		if node.next[level] != nil && bytes.Compare(node.next[level].key, key) == 0 {
			return node.next[level]
		}
	}
	return nil
}

// Get 查找key，返回value，第二个bool标识是否有数据
func (s *Skiplist) Get(key []byte) ([]byte, bool) {
	// key能找到，返回对应的value
	if node := s.getNode(key); node != nil {
		return node.value, true
	}
	return nil, false
}

// Put 插入一条kv记录
func (s *Skiplist) Put(key, value []byte) {
	// key已经存在，更新value
	if node := s.getNode(key); node != nil {
		// 更新跳表数据内存大小
		s.memorySize += len(value) - len(node.value)
		node.value = value
		return
	}

	// key不存在，插入新的节点
	// 更新跳表数据内存大小，kv对的数量
	s.memorySize += len(key) + len(value)
	s.entryCount++

	newHeight := s.roll()
	head := s.head
	if len(head.next) < newHeight {
		diff := make([]*Node, newHeight+1-len(head.next))
		head.next = append(head.next, diff...)
	}

	newNode := &Node{
		next:  make([]*Node, newHeight),
		key:   key,
		value: value,
	}

	// 从最高层开始遍历，每层依次插入节点
	node := head
	for level := newHeight - 1; level >= 0; level-- {
		// 找到下一个节点不存在或比key大，此时下一个节点的位置就是新插入的位置
		for node.next[level] != nil && bytes.Compare(node.next[level].key, key) < 0 {
			node = node.next[level]
		}
		// 插入新节点
		newNode.next[level] = node.next[level]
		node.next[level] = newNode
	}
}

// Erase 删除key对应的节点
func (s *Skiplist) Erase(key []byte) bool {
	flag := false // 标志是否删除了节点
	// 从最高层开始遍历，每层依次插入节点
	node := s.head
	for level := len(s.head.next) - 1; level >= 0; level-- {
		// 找到下一个不存在或比key大的节点
		for node.next[level] != nil && bytes.Compare(node.next[level].key, key) < 0 {
			node = node.next[level]
		}
		// 找到了需要删除的节点
		if node.next[level] != nil && bytes.Compare(node.next[level].key, key) == 0 {
			// 更新删除节点后占用的内存大小
			if !flag {
				s.memorySize = s.memorySize - len(node.next[level].key) - len(node.next[level].value)
			}
			node.next[level] = node.next[level].next[level]
			flag = true
		}
	}
	return flag
}

// AllData 返回所有的kv数据对
func (s *Skiplist) AllData() []*KV {
	if len(s.head.next) == 0 {
		return nil
	}
	kvs := make([]*KV, 0, s.entryCount)
	// 从最底层开始遍历，获取所有的kv数据对
	for node := s.head.next[0]; node != nil; node = node.next[0] {
		kvs = append(kvs, &KV{
			Key:   node.key,
			Value: node.value,
		})
	}
	return kvs
}

// EntryCount 返回kv数据对的个数
func (s *Skiplist) EntryCount() int {
	return s.entryCount
}

// MemorySize 返回跳表内存占用大小，单位byte
func (s *Skiplist) MemorySize() int {
	return s.memorySize
}
