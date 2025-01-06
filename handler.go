package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"GET":  get,
	"SET":  set,
}

var SetMap map[string]string

var SetMutex sync.RWMutex

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Invalid number of arguments"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SetMutex.Lock()
	SetMap[key] = value
	SetMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Invalid number of arguments"}
	}

	key := args[0].bulk

	SetMutex.RLock()
	value, ok := SetMap[key]
	SetMutex.RUnlock()

	if !ok {
		return Value{str: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}
