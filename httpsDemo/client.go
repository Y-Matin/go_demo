package httpsDemo

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Client(req []string) ([]bool, error) {
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile("cert/ca.pem")
	if err != nil {
		log.Println("read ca.crt file error:", err.Error())
		return nil, err
	}
	pool.AppendCertsFromPEM(caCrt)
	cliCrt, err := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	if err != nil {
		log.Println("LoadX509KeyPair error:", err.Error())
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
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
