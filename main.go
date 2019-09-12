package main

import (
	"crypto/tls"
	"gopkg.in/yaml.v2"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"time"
)
/*Based on
https://github.com/bechurch/reverse-proxy-demo
https://github.com/txn2/p3y/blob/master/p3y.go
https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c
*/

type conf struct {
	Entries []string `yaml:"block"`
}

//https
func handleTunneling(res http.ResponseWriter, req *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", req.Host, 10*time.Second)
	if err != nil {
		http.Error(res, err.Error(), http.StatusServiceUnavailable)
		return
	}
	res.WriteHeader(http.StatusOK)
	hijacker, ok := res.(http.Hijacker)
	if !ok {
		http.Error(res, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(res, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}
func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

// http
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func handleRequestCustom(res http.ResponseWriter, req *http.Request) {
	transport := http.DefaultTransport
	out, err := transport.RoundTrip(req)
	if err != nil {
		http.Error(res, err.Error(), http.StatusServiceUnavailable)
		return
	}
	copyHeader(res.Header(), out.Header)
	res.WriteHeader(out.StatusCode)
	_, err = io.Copy(res, out.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	err = out.Body.Close()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}


func loadBlockedList(filename string) []regexp.Regexp {
	var c conf
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	regexps := []regexp.Regexp{}
	for _, condition := range c.Entries {
		//log.Printf("%s", condition)
		r := regexp.MustCompile(condition)
		regexps = append(regexps, *r)
	}
	return regexps
}

func shouldBlock(regList []regexp.Regexp, url string) bool {
	for _, condition := range regList {
		if condition.MatchString(url) {
			return true
		}
	}
	return false
}

func main() {
	port :=  "0.0.0.0:11666"
	log.Printf("Server will run on: %s\n", port)
	http.HandleFunc("/", handleRequestCustom)
	regList := loadBlockedList("block.yaml")
	server := &http.Server{
		Addr: port,
		Handler: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if shouldBlock(regList, req.Host) {
				//log.Printf("Blocked %s\n", req.Host)
				if wr, ok := res.(http.Hijacker); ok {
					conn, _, err := wr.Hijack()
					if err != nil {
						fmt.Fprint(res, err)
					}
					conn.Close()
				}
			} else {
				if req.Method == http.MethodConnect {
					log.Printf("proxy_url: %s\n", req.Host)
					handleTunneling(res, req)
				} else {
					log.Printf("proxy_url: %s%s\n", req.Host, req.RequestURI)
					handleRequestCustom(res, req)
				}
			}

		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	server.ListenAndServe()
}
