package httpsDemo

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Client(req []string) ([]bool, error) {
	//跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.PostForm("https://localhost:8080/", url.Values{"data": req})
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println("close error:", err.Error())
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read response Body error:", err)
		return nil, err
	}
	split := strings.Split(string(body[1:len(body)-1]), " ")
	result := make([]bool, len(req))
	for i, v := range split {
		parseBool, err := strconv.ParseBool(v)
		if err != nil {
			log.Println("strconv.ParseBool error:", err)
		}
		result[i] = parseBool
	}
	return result, err
}
