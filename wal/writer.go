package wal

import (
	"encoding/binary"
	"os"
)

type Writer struct {
	file      string   // 预写日志文件名，绝对路径
	dest      *os.File // 预写日志文件
	lenBuffer [30]byte // 存放key和value的长度缓冲区
}

func NewWriter(file string) (*Writer, error) {
	// 打开预写日志文件
	dest, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file: file,
		dest: dest,
	}, nil
}

func (w *Writer) Write(key, value []byte) error {
	// 将key和value的长度写入缓冲区当中
	n := binary.PutUvarint(w.lenBuffer[0:], uint64(len(key)))
	n += binary.PutUvarint(w.lenBuffer[n:], uint64(len(value)))

	var buffer []byte
	buffer = append(buffer, w.lenBuffer[:n]...)
	buffer = append(buffer, key...)
	buffer = append(buffer, value...)
	_, err := w.dest.Write(buffer)
	return err
}

func (w *Writer) Close() {
	_ = w.dest.Close()
}
