package memtable

// MemTableConstructor memtable 构造器
type MemTableConstructor func() MemTable

type MemTable interface {
	Put(key, value []byte)         // 写入数据
	Get(key []byte) ([]byte, bool) // 读取数据，第二个标识数据是否读取成功
	Erase(key []byte) bool         // 删除一条kv记录
	AllData() []*KV                // 返回所有的kv数据对
	MemorySize() int               // 内存表的数据大小，单位byte
	EntryCount() int               // kv数据对的个数
}

type KV struct {
	Key, Value []byte
}
