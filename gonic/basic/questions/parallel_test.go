package questions

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type testStruct struct {
	TestI   int
	TestStr string
}

func Test_parseSliceInterface(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				v: []string{"12312", "aksjal", "JHASKD"},
			},
			want:    []interface{}{"12312", "aksjal", "JHASKD"},
			wantErr: false,
		},
		{
			name: "case 2",
			args: args{
				v: []testStruct{{1, "1"}, {2, "2"}, {3, "3"}},
			},
			want:    []interface{}{testStruct{1, "1"}, testStruct{2, "2"}, testStruct{3, "3"}},
			wantErr: false,
		},
		{
			name: "case 3",
			args: args{
				v: &[]testStruct{{1, "1"}, {2, "2"}, {3, "3"}},
			},
			want:    []interface{}{testStruct{1, "1"}, testStruct{2, "2"}, testStruct{3, "3"}},
			wantErr: false,
		},
		{
			name: "case 4",
			args: args{
				v: "12312",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSliceInterface(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSliceInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSliceInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 测试用工作函数 对items求和
func testWorkerFuncSum(items []interface{}, params ...interface{}) (interface{}, error) {
	if len(params) == 0 {
		return nil, errors.New("want addtional params")
	}

	sum := 0
	for _, v := range items {
		vint, _ := v.(int)
		sum += vint
	}
	return sum, nil
}

func TestNew(t *testing.T) {
	type args struct {
		name          string
		parallelCount int
		batchSize     int
		worker        WorkFunc
		itemList      interface{}
		params        []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				name:          "task 1",
				parallelCount: 3,
				batchSize:     3,
				worker:        testWorkerFuncSum,
				itemList:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				params:        nil,
			},
			wantErr: true,
		},
		{
			name: "case 1",
			args: args{
				name:          "task 2",
				parallelCount: 4,
				batchSize:     4,
				worker:        testWorkerFuncSum,
				itemList:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				params:        nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.name, tt.args.parallelCount, tt.args.batchSize,
				tt.args.worker, tt.args.itemList, tt.args.params...); (got.Error != nil) != tt.wantErr {
				t.Errorf("New() got err = %v, %v, wantErr %v", got.Error, (got.Error != nil), tt.wantErr)
			}
		})
	}
}

func Test_ParallelTask_case1(t *testing.T) {
	var (
		name          = "task 1"
		parallelCount = 2
		batchSize     = 7
		timeout       = 5 * time.Second
		worker        = testWorkerFuncSum
		itemList      = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
		params        = []interface{}{}
	)

	task := New(name, parallelCount, batchSize, worker, itemList, params...)
	if task.Error != nil {
		t.Errorf("new task got err: %v", task.Error)
		t.FailNow()
	}

	// t.Logf("error: %v", task.Error)

	// run
	go func() {
		task.Run(timeout)
	}()

	select {
	case <-task.Done():
		result := task.GetResult()
		for idx, v := range result {
			t.Logf("result: %d, result: %v, error: %v, item list offset: %d", idx, v.Data(), v.Error(), v.Offset())
		}
	}
}

// 测试用工作函数 对items求和 测试超时控制
func testWorkerFuncSumDelay(items []interface{}, params ...interface{}) (interface{}, error) {
	if len(params) == 0 {
		return nil, errors.New("want addtional params")
	}

	sum := 0
	for _, v := range items {
		vint, _ := v.(int)
		sum += vint
	}

	time.Sleep(time.Second * 5)

	return sum, nil
}

func Test_ParallelTask_case2(t *testing.T) {
	var (
		name          = "task 1"
		parallelCount = 4
		batchSize     = 5
		timeout       = time.Duration(0)
		worker        = testWorkerFuncSumDelay
		itemList      = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		params        = []interface{}{}
	)

	task := New(name, parallelCount, batchSize, worker, itemList, params...)
	if task.Error != nil {
		t.Errorf("new task got err: %v", task.Error)
		t.FailNow()
	}

	// t.Logf("error: %v", task.Error)

	// run
	go func() {
		task.Run(timeout)
	}()

	select {
	case <-task.Done():
		result := task.GetResult()
		for idx, v := range result {
			t.Logf("result: %d, result: %v, error: %v, item list offset: %d", idx, v.Data(), v.Error(), v.Offset())
		}
	}
}
