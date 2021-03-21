package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/req", handReq)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Printf("http server failed,error:%v\n", err)
	}
}

func sayHello(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "<h1>Hello Golang</h1><script>alert('xss攻击')</script>")
}

func handReq(w http.ResponseWriter, r *http.Request) {
	fmt.Println("进入、、、")

	method := r.Method
	switch method {
	case "GET", "DELETE":
		fmt.Println("get/delete :", method)
		values := r.URL.Query()
		fmt.Printf("%v", values)
	case "POST", "PUT":
		fmt.Println("post/put :", method)
		contentType, ok := r.Header["Content-Type"]
		if !ok {
			fmt.Println("Content-Type is not set")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("Content-Type:%v\n", contentType)
		for _, v1 := range contentType {
			switch {
			case v1 == "text/plain":
				fmt.Println("text/plain")
				all, err := ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Println("read body error:", err)
				}
				fmt.Println("text:", string(all))
			case v1 == "application/json":
				fmt.Println("application/json")
				all, err := ioutil.ReadAll(r.Body)
				if err != nil {
					fmt.Println("read body error:", err)
				}
				fmt.Println("json:", string(all))
			case strings.Contains(v1, "multipart/form-data"):
				fmt.Println("multipart/form-data")
				err := r.ParseMultipartForm(128)
				if err != nil {
					fmt.Println("parse multipartForm error:", err)
				}
				marshal, err := json.Marshal(r.MultipartForm)

				fmt.Println("multipart :", string(marshal))
			case v1 == "application/x-www-form-urlencoded":
				fmt.Println("application/x-www-form-urlencoded")
				err := r.ParseForm()
				if err != nil {
					fmt.Println("ParseForm error:", err)
				}
				marshal, err := json.Marshal(r.PostForm)
				if err != nil {
					fmt.Println("json error:", err)
				}

				fmt.Printf("urlencoded:%v\n", string(marshal))
			}

		}

	}

}
