package main

import (
	"./proxy"
	"crypto/tls"
	"fmt"
	"github.com/akamensky/argparse"
	"net/http"
	"os"
)

func readArgs() (string, string, string) {
	parser := argparse.NewParser("blocking-http-proxy", "HTTP/S proxy that blocks")
	host := parser.String("", "host", &argparse.Options{Help: "host", Required: false})
	port := parser.String("", "port", &argparse.Options{Help: "port", Required: false})
	blockFile := parser.String("", "block", &argparse.Options{Help: "YAML block file", Required: false})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
	}
	if *host == "" {
		*host = "0.0.0.0"
	}
	if *port == "" {
		*port = "11666"
	}
	if *blockFile == "" {
		*blockFile = "block.yaml"
	}
	return *host, *port, *blockFile
}

func main() {
	logger := proxy.NewLogger()
	blockedLogger := proxy.NewFileLogger("block.log")
	allowLogger := proxy.NewFileLogger("allow.log")
	host, port, blockFile := readArgs()
	address :=  host+":"+port
	c := proxy.NewConfig()
	fmt.Println(len(os.Args), os.Args)
	logger.Printf("Server will run on: %s\n", address)
	regList := c.LoadBlockedList(blockFile)
	server := &http.Server{
		Addr: address,
		Handler: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if proxy.ShouldBlock(regList, req.Host) {
				blockedLogger.WriteLine("Blocked %s\n", false, req.Host)
				if wr, ok := res.(http.Hijacker); ok {
					conn, _, err := wr.Hijack()
					if err != nil {
						fmt.Fprint(res, err)
					}
					conn.Close()
				}
			} else {
				if req.Method == http.MethodConnect {
					allowLogger.WriteLine("proxy_url: %s\n", true, req.Host)
					proxy.HandleTunneling(res, req)
				} else {
					allowLogger.WriteLine("proxy_url: %s %s\n", true, req.Host, req.RequestURI)
					proxy.HandleRequestCustom(res, req)
				}
			}

		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	server.ListenAndServe()
}
