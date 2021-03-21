package httpsDemo

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var data map[string]struct{}
var mutex sync.Mutex

func Server() {
	http.HandleFunc("/", handler)
	s := &http.Server{
		Addr: ":8080",
	}
	data=make(map[string]struct{})
	// server 任何情况下均不能退出进程
	for true {
		err := s.ListenAndServeTLS("cert/server.crt", "cert/server_no_passwd.key")
		if err != nil {
			log.Println("server: ListenAndServeTLS error:", err)
		}
	}

}

// 处理函数
func handler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("ParseForm error:", err)
		return
	}

	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Println("r.Body.Close error:", err)
		}
	}()
	strs := r.PostForm["data"]
	result := make([]bool, len(strs))
	// 加锁, (1.读取key，2.若key不存在则赋值。) 这两步合起来不是原子性的
	mutex.Lock()
	for i, v := range strs {
		_, ok := data[v]
		if ok {
			result[i] = true
		} else {
			result[i] = false
			data[v]= struct{}{}
		}
	}
	mutex.Unlock()
	_, err = fmt.Fprint(w, result)
	if err != nil {
		log.Println("return result error:", err)
	}

}
