package wal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"ephmeral.com/v2/memtable"
	"errors"
	"io"
	"os"
)

type Reader struct {
	file   string        // 预写日志文件名，绝对路径
	src    *os.File      // 预写日志文件
	reader *bufio.Reader // 基于bufio reader对日志的封装
}

// NewReader 构造器
func NewReader(file string) (*Reader, error) {
	src, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Reader{
		file:   file,
		src:    src,
		reader: bufio.NewReader(src),
	}, nil
}

// Recover 恢复文件数据到内存当中
func (r *Reader) Recover(memtable memtable.MemTable) error {
	// 读取预写日志全量内容
	body, err := io.ReadAll(r.reader)
	if err != nil {
		return err
	}
	// 保证文件偏移量重置为0
	defer func() {
		_, _ = r.src.Seek(0, io.SeekStart)
	}()

	kvs, err := r.readAll(bytes.NewReader(body))
	if err != nil {
		return err
	}
	for _, kv := range kvs {
		memtable.Put(kv.Key, kv.Value)
	}
	return nil
}

// 从文件当中读取所有的kv对
func (r *Reader) readAll(reader *bytes.Reader) ([]*memtable.KV, error) {
	var kvs []*memtable.KV
	for {
		// 从reader当中读取key的长度
		keyLen, err := binary.ReadUvarint(reader)
		// 读到EOF说明文件内容已经读完
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		// 从reader当中读取value的长度
		valueLen, err := binary.ReadUvarint(reader)
		if err != nil {
			return nil, err
		}
		// 从reader当中读取key
		key := make([]byte, keyLen)
		if _, err := io.ReadFull(reader, key); err != nil {
			return nil, err
		}
		// 从reader当中读取value
		value := make([]byte, valueLen)
		if _, err := io.ReadFull(reader, value); err != nil {
			return nil, err
		}
		kvs = append(kvs, &memtable.KV{
			Key:   key,
			Value: value,
		})
	}
	return kvs, nil
}

// Close 关闭文件流
func (r *Reader) Close() {
	r.reader.Reset(r.src)
	_ = r.src.Close()
}
