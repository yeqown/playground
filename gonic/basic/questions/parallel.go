package questions

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// WorkFuncResult 任务处理结果
type WorkFuncResult interface {
	// data 包含任务函数的结果
	Data() interface{}
	// get 任务函数的异常
	Error() error
	// item 任务函数处理的item偏移量，用于检查结果
	Offset() int
}

var (
	_ WorkFuncResult = &defaultWorkFuncResult{}
)

type defaultWorkFuncResult struct {
	v      interface{} // worker func result
	err    error       // worker func error
	offset int         // worker func handle itemList offset
}

func (r *defaultWorkFuncResult) Data() interface{} {
	return r.v
}

func (r *defaultWorkFuncResult) Error() error {
	return r.err
}

func (r *defaultWorkFuncResult) Offset() int {
	return r.offset
}

// ParallelTask ...
type ParallelTask struct {
	Name          string   // 任务名称
	ParallelCount int      // 并行数
	BatchSize     int      // 批量数
	Error         error    // 执行过程中出现的错误
	Func          WorkFunc // 任务函数

	mu            sync.Mutex          // lock
	done          chan struct{}       // 是否完成
	itemList      []interface{}       // itemList
	itemListCount int                 // len(itemList)
	resultChan    chan WorkFuncResult // 任务执行结果
	result        []WorkFuncResult    // result slice
	params        []interface{}       // 其他参数
	// timeout       time.Duration       // 超时设置
}

// WorkFunc 执行子任务的函数原型
// batchItemList 子任务参数Slice
// params 附加参数
// result 子任务的返回值
type WorkFunc func(batchItemList []interface{}, params ...interface{}) (interface{}, error)

// New 创建一新的并行任务
// name 任务名
// parallelCount 并行数（即多少个goroutine同时在执行这个任务）
// batchSize 单次任务访问的item数量，
// itemList 总的任务信息，它应该是一个Slice，这里写成interface{}类型是为了兼容所有的不同Slice类型
// params 其它参数
func New(name string, parallelCount, batchSize int, worker WorkFunc, itemList interface{}, params ...interface{}) *ParallelTask {

	task := &ParallelTask{
		Name:          name,
		ParallelCount: parallelCount,
		BatchSize:     batchSize,
		Error:         nil,
		Func:          worker,
		done:          make(chan struct{}),
		params:        params,
		mu:            sync.Mutex{},
		// timeout:       timeout,
	}

	// 解析slice interface 到 []interface{}
	task.itemList, task.Error = parseSliceInterface(itemList)
	task.itemListCount = len(task.itemList)
	// 检查并行数和批量数是否合法
	if parallelCount*batchSize < task.itemListCount {
		task.Error = errItemListCountOverflow
	}

	task.resultChan = make(chan WorkFuncResult, task.ParallelCount)
	return task
}

// Run 执行任务
func (t *ParallelTask) Run(timeout time.Duration) *ParallelTask {
	wg := sync.WaitGroup{}
	// 0 表示无限等待
	if timeout == 0 {
		timeout = time.Minute * 9999
	}

	// 并行处理任务
	for i := 0; i < t.ParallelCount; i++ {
		// get batchItemList
		start, end := i*t.BatchSize, (i+1)*t.BatchSize
		grCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if start >= t.itemListCount {
			break
		}
		// end overflow
		if end > len(t.itemList) {
			end = len(t.itemList)
		}

		wg.Add(1)
		debugF("start: %d, end: %d, t.itemListCount: %d\n", start, end, t.itemListCount)

		// start an goroutine
		go func(ctx context.Context, batchItemList []interface{}, offset int, params ...interface{}) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// timeout
				t.resultChan <- &defaultWorkFuncResult{
					v:      nil,
					err:    errWorkerTimeout,
					offset: offset,
				}
			default:
				// handler worker func result
				r, err := t.Func(batchItemList, params)
				t.resultChan <- &defaultWorkFuncResult{
					v:      r,
					err:    err,
					offset: offset,
				}
			}
			return
		}(grCtx, t.itemList[start:end], start, t.params)
	}

	wg.Wait()
	debugF("goroutine done\n")
	close(t.resultChan)

	// 合并任务函数处理结果
	for r := range t.resultChan {
		debugF("result: %v\n", r)
		t.result = append(t.result, r)
	}

	// it blocked
	t.done <- struct{}{}

	return t
}

// GetResult 得到并行任务的执行结果
// duration是等待的时间，为0表示无限等待，如果超时后任务没有完成，需要返回相应的错误提示
// result 为并行任务成功完成后的返回值，没有返回值可以返回nil
func (t *ParallelTask) GetResult() (result []WorkFuncResult) {
	return t.result
}

// Done ...
func (t *ParallelTask) Done() chan struct{} {
	t.mu.Lock()
	if t.done == nil {
		t.done = make(chan struct{})
	}
	d := t.done
	t.mu.Unlock()
	return d
}

var (
	errWorkerTimeout         = errors.New("worker timeout")                                     // 任务函数处理超时错误
	errInvalidTypOfItemList  = errors.New("invalid item list type, must be slice")              // itemList类型错误
	errItemListCountOverflow = errors.New("items count is bigger than parallelCount*batchSize") // itemList数量超出并行数
)

// 解析一个interface{} 到 []interface{}, 如果不是slice类型则报错
func parseSliceInterface(v interface{}) ([]interface{}, error) {
	// 检查是否slice类型
	val := reflect.ValueOf(v)
	typ := val.Type()

	// 如果是指针
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	// 如果不是slice类型则报错
	if typ.Kind() != reflect.Slice {
		return nil, errInvalidTypOfItemList
	}

	// fmt.Println(typ.Kind(), val.Len())
	// 获取整个slice的item，并存放在result中
	result := make([]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		result[i] = val.Index(i).Interface()
	}
	return result, nil
}

var (
	debug = true
)

func debugF(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format, args...)
	}
}
