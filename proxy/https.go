package proxy

import (
	"io"
	"net"
	"net/http"
	"time"
)

func HandleTunneling(res http.ResponseWriter, req *http.Request) {
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
