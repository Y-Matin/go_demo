package main

import (
	"fmt"
	"httpsDemo"
	"sync"
	"time"
)

func main() {
	go httpsDemo.Server()
	time.Sleep(time.Millisecond * 500)

	strings := []string{"A", "B", "C", "D", "E"}
	allResult := []bool{}
	mutex := &sync.Mutex{}
	group := &sync.WaitGroup{}
	group.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			bools, err := httpsDemo.Client(strings)
			if err != nil {
				fmt.Println(err)
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
		fmt.Println("多个client同时请求 ，返回结果错误")
	}
	time.Sleep(time.Second)
}
