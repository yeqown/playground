package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	leveldbopt "github.com/syndtr/goleveldb/leveldb/opt"
)

type FIFOQueue struct {
	queueName string
	queuePath string
	db        *leveldb.DB

	low           uint64 // low is the lowest key in the queue
	high          uint64 // high is the highest key in the queue
	lowRW, highRW sync.RWMutex
}

func NewFIFOQueue(path, name string) (*FIFOQueue, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}

	if name == "" {
		name = "fifo-queue"
	}

	leveldbOpts := &leveldbopt.Options{
		Compression: leveldbopt.NoCompression,
	}
	queuePath := filepath.Join(path, name)
	db, err := leveldb.OpenFile(queuePath, leveldbOpts)
	if err != nil {
		return nil, err
	}
	q := &FIFOQueue{
		queuePath: queuePath,
		queueName: name,
		db:        db,
	}
	q.init()
	return q, nil
}

// init initializes the queue or restores the queue from the database
func (q *FIFOQueue) init() {
	q.lowRW.Lock()
	defer q.lowRW.Unlock()
	q.highRW.Lock()
	defer q.highRW.Unlock()

	q.low = q.getLow()
	q.high = q.getHigh()
}

var (
	// lowKey is the key of the lowest key in the queue
	lowKey = []byte("low")
	// highKey is the key of the highest key in the queue
	highKey = []byte("high")
	// keyOffset is the offset of the key in the queue
	keyOffset uint64 = 210705
)

func (q *FIFOQueue) prefix(key []byte) []byte {
	return append([]byte(q.queueName+"-"), key...)
}

func (q *FIFOQueue) getLow() uint64 {
	lowBytes, err := q.db.Get(q.prefix(lowKey), nil)
	if err != nil {
		return keyOffset
	}

	return bytesToUint64(lowBytes)
}

func (q *FIFOQueue) getHigh() uint64 {
	highBytes, err := q.db.Get(q.prefix(highKey), nil)
	if err != nil {
		return keyOffset
	}

	return bytesToUint64(highBytes)
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)
	return b
}

func (q *FIFOQueue) Enqueue(data []byte) error {
	q.highRW.Lock()
	defer q.highRW.Unlock()

	if q.db == nil {
		return fmt.Errorf("db is nil")
	}

	key := q.prefix(uint64ToBytes(q.high))
	q.high++
	// TODO: Is there a better way to do this? update the high key and then insert the data?
	q.db.Put(q.prefix(highKey), uint64ToBytes(q.high), nil)
	if err := q.db.Put(key, data, nil); err != nil {
		return fmt.Errorf("enqueue(%x): %v", key, err)
	}

	return nil
}

func (q *FIFOQueue) Dequeue() ([]byte, error) {
	q.lowRW.Lock()
	defer q.lowRW.Unlock()

	if q.db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	key := q.prefix(uint64ToBytes(q.low))
	q.low++

	q.db.Put(q.prefix(lowKey), uint64ToBytes(q.low), nil)
	data, err := q.db.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("dequeue(%x): %v", key, err)
	}
	if err = q.db.Delete(key, nil); err != nil {
		return nil, fmt.Errorf("dequeue(rm %x): %v", key, err)
	}

	return data, nil
}

func (q *FIFOQueue) Close() error {
	if q.db == nil {
		return nil
	}

	if err := q.db.Close(); err != nil {
		return err
	}

	q.db = nil
	return nil
}

// Clear clears the queue by deleting all keys in the queue
// and remove the level db file
func (q *FIFOQueue) Clear() error {
	q.lowRW.Lock()
	defer q.lowRW.Unlock()
	q.highRW.Lock()
	defer q.highRW.Unlock()

	q.low = 0
	q.high = 0

	q.Close()

	// TODO: Is there a better way to do this?
	if err := os.RemoveAll(q.queuePath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func main() {
	q, err := NewFIFOQueue("./testdata", "my-queue")
	if err != nil {
		panic(err)
	}

	defer q.Close()

	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()

			err := q.Enqueue([]byte(fmt.Sprintf("hello-%d", i)))
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Printf("Enqueue: hello-%d\n", i)
		}(i)
	}
	wg.Wait()

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			data, err := q.Dequeue()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println("Dequeue: " + string(data))
		}()
	}
	wg.Wait()

	if err = q.Clear(); err != nil {
		panic(err)
	}

	if err = q.Close(); err != nil {
		panic(err)
	}
}
