package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", ":6377")
	if err != nil {
		slog.Error(err.Error(), err)
		panic(err)
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		slog.Error(err.Error(), err)
		panic(err)
	}
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			slog.Error("error reading from client: ", err)
			os.Exit(1)
		}
	}
}
