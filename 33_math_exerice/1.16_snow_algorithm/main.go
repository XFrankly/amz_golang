package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	WORKER_ID_BITS int64 = 4
	SEQUENCE_BITS  int64 = 5

	MAX_WORKER_ID        int64 = -1 ^ (-1 << WORKER_ID_BITS)
	WORKER_ID_SHIFT      int64 = SEQUENCE_BITS
	TIMESTAMP_LEFT_SHIFT int64 = SEQUENCE_BITS + WORKER_ID_SHIFT
	SEQUENCE_MASK        int64 = -1 ^ (-1 << SEQUENCE_BITS)
	TWEPOCH              int64 = 1632325994945
)

var (
	Logg = log.New(os.Stderr, "INFO -:", 18)
)

type SystemClockError error

type IdWorker struct {
	worker_id      int64
	sequence       int64
	last_timestamp int64
}

func GenIdWorker(mac_id int64, sequence int64) *IdWorker {
	return &IdWorker{worker_id: mac_id, sequence: sequence, last_timestamp: -1}
}

func (iw *IdWorker) GetTimestamp() int64 {
	/// 返回从1970到的当前 13位时间戳，now.UnixMilli()
	// fmt.Println(time.Now().UnixMilli())
	return time.Now().UnixMilli()
}

func (iw *IdWorker) WaitNextMillis(last_timestamp int64) int64 {
	/// 等待下一毫秒
	timestamp := iw.GetTimestamp()
	for timestamp <= last_timestamp {
		timestamp = iw.GetTimestamp()
	}
	return timestamp
}

func (iw *IdWorker) GetId(s int64) int64 {
	// 计算下一个Id
	timestamp := iw.GetTimestamp()
	if timestamp < iw.last_timestamp {
		iw.sequence = (iw.sequence + 1) & SEQUENCE_MASK
		if iw.sequence == 0 {
			timestamp = iw.WaitNextMillis(iw.last_timestamp)
		}
	} else {
		iw.sequence = 0
	}
	iw.last_timestamp = timestamp
	first := (timestamp-TWEPOCH)<<TIMESTAMP_LEFT_SHIFT - s
	second := iw.worker_id >> WORKER_ID_SHIFT
	new_id := first | second | iw.sequence
	return new_id
}

func main() {
	worker := GenIdWorker(10, 1)
	var near_id int64 = 0
	for i := 0; i < 100; i += 2 {
		new_id := worker.GetId(int64(i))
		if near_id == new_id {
			i -= 2
			continue
		} else {
			fmt.Println(new_id)
		}
	}

}
