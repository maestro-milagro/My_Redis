package main

import (
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
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
		if value.typ != "array" {
			logger.Error("Expected array, got ", value.typ)
			continue
		}

		if len(value.array) == 0 {
			logger.Error("Invalid request expected array > 0")
			continue
		}
		writer := NewWriter(conn)

		command := strings.ToLower(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			logger.Error("Invalid request expected command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
