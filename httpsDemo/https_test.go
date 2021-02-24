package httpsDemo

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

// TestMain 初始化：启动server
func TestMain(m *testing.M) {
	fmt.Println("begin")
	//启动server
	go Server()
	// 休眠1秒，确保server启动
	time.Sleep(time.Second)
	m.Run()
	fmt.Println("end")
}

// TestServer 正常流程
func TestClient(t *testing.T) {
	// 启动client
	got, err := Client([]string{"aa", "bb", "aa"})
	if err != nil {
		t.Error(err)
	}
	want := []bool{false, false, true}
	// 验证结果
	if !reflect.DeepEqual(got, want) {
		t.Errorf("excepted:%v, got:%v", want, got)
	}
}

// TestMultiClient 多用例
func TestMultiClient(t *testing.T) {
	type test struct {
		input []string
		want  []bool
	}
	tests := []test{
		{input: []string{"a", "b", "c"}, want: []bool{false, false, false}},
		{input: []string{"a", "b", "c"}, want: []bool{true, true, true}},
		{input: []string{"d", "e", "f"}, want: []bool{false, false, false}},
		{input: []string{"a", "d", "g"}, want: []bool{true, true, false}},
	}
	for _, tc := range tests {
		// 启动client
		got, err := Client(tc.input)
		if err != nil {
			t.Error(err)
		}
		// 验证结果
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("excepted:%v, got:%v", tc.want, got)
		}
	}
}

// TestMultiClientSync  模拟多个client同时 请求
func TestMultiClientSync(t *testing.T) {
	strings := []string{"A", "B", "C", "D", "E"}
	allResult := []bool{}
	mutex := &sync.Mutex{}
	group := &sync.WaitGroup{}
	group.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			bools, err := Client(strings)
			if err != nil {
				t.Error(err)
			}
			mutex.Lock()
			allResult = append(allResult, bools...)
			mutex.Unlock()
			group.Done()
		}()
	}
	group.Wait()
	//如果只存在5个false
	count := 0
	for _, v := range allResult {
		if v == false {
			count++
		}
	}
	if count != len(strings) {
		t.Error("多个client同时请求 ，返回结果错误")
	}
}
