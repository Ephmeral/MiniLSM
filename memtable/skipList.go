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
	head     *Node // 头节点
	maxLevel int   // 最大高度
}

func NewSkipList() Skiplist {
	head := &Node{[]*Node{nil}, []byte{}, []byte{}}
	return Skiplist{head, 1}
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

func (s *Skiplist) Search(key []byte) ([]byte, bool) {
	// key能找到，返回对应的value
	if node := s.getNode(key); node != nil {
		return node.value, true
	}
	return nil, false
}

func (s *Skiplist) Put(key, value []byte) {
	// key已经存在，更新value
	if node := s.getNode(key); node != nil {
		node.value = value
		return
	}

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

func (s *Skiplist) Erase(key []byte) bool {
	flag := false // 标志是否删除了节点
	// 从最高层开始遍历，每层依次插入节点
	node := s.head
	for level := len(s.head.next) - 1; level >= 0; level-- {
		// 找到下一个节点不存在或比key大
		for node.next[level] != nil && bytes.Compare(node.next[level].key, key) < 0 {
			node = node.next[level]
		}
		// 找到了下一个节点
		if node.next[level] != nil && bytes.Compare(node.next[level].key, key) == 0 {
			node.next[level] = node.next[level].next[level]
			flag = true
		}
	}
	return flag
}
