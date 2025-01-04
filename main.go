package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

func main() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	l, err := net.Listen("tcp", ":6377")
	if err != nil {
		logger.Error(err.Error(), err)
		panic(err)
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		logger.Error(err.Error(), err)
		panic(err)
	}
	defer conn.Close()
	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			logger.Error("error reading from client: ", err)
			os.Exit(1)
		}
		fmt.Println(value)

		writer := NewWriter(conn)
		writer.Write(Value{str: "OK", typ: "string"})
	}
}
