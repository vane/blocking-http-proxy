package proxy

import (
	"io"
	"net/http"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func HandleRequestCustom(res http.ResponseWriter, req *http.Request) {
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
