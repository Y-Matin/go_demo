package httpsDemo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var data sync.Map

func Server() {
	pool := x509.NewCertPool()
	crt, err := ioutil.ReadFile("cert/ca.pem")
	if err != nil {
		log.Fatalln("读取证书失败！", err.Error())
	}
	pool.AppendCertsFromPEM(crt)
	http.HandleFunc("/", handler)
	s := &http.Server{
		Addr: ":8080",
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert, // 检验客户端证书
		},
	}
	// server 任何情况下均不能退出进程
	for true {
		err = s.ListenAndServeTLS("cert/server.pem", "cert/server.key")
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
	for i, v := range strs {
		_, ok := data.Load(v)
		if ok {
			result[i] = true
		} else {
			result[i] = false
			data.Store(v, struct{}{})
		}
	}
	_, err = fmt.Fprint(w, result)
	if err != nil {
		log.Println("return result error:", err)
	}

}
