package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"GET":     get,
	"SET":     set,
	"HSET":    hSet,
	"HGET":    hGet,
	"HGETALL": hGetAll,
}

var SetMap = make(map[string]string)
var SetMutex sync.RWMutex

var HSetMap = make(map[string]map[string]string)
var HSetMutex sync.RWMutex

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
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hSet(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "Invalid number of arguments"}
	}
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSetMutex.Lock()
	if _, ok := HSetMap[hash]; !ok {
		HSetMap[hash] = make(map[string]string)
	}
	HSetMap[hash][key] = value
	HSetMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hGet(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "Invalid number of arguments"}
	}
	hash := args[0].bulk
	key := args[1].bulk

	HSetMutex.RLock()
	value, ok := HSetMap[hash][key]
	HSetMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hGetAll(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "Invalid number of arguments"}
	}
	hash := args[0].bulk

	HSetMutex.RLock()
	valueMap, ok := HSetMap[hash]
	HSetMutex.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	values := make([]Value, 0, len(valueMap))

	for _, v := range valueMap {
		values = append(values, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: values}
}
