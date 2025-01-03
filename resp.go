package main

import (
	"bufio"
	"io"
	"strconv"
)

const (
	STRING  = "+"
	INTEGER = ":"
	BULK    = "$"
	ERROR   = "-"
	ARRAY   = "*"
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func (r *Resp) New(reader io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(reader),
	}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n++
		line = append(line, b)
		if b == '\r' && line[n-1] == '\n' && len(line) >= 2 {
			break
		}
	}
	return line, n, nil
}

func (r *Resp) readInteger() (number int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	_, err = strconv.ParseInt(string(line), 0, 64)
	if err != nil {
		return 0, n, err
	}
	return number, n, nil
}

func (r *Resp) readBulk() (Value, error) {
	var value Value
	value.typ = "bulk"

	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	bulk := make([]byte, length)

	_, err = r.reader.Read(bulk)
	if err != nil {
		return Value{}, err
	}

	value.bulk = string(bulk)

	r.readLine()

	return value, nil
}

func (r *Resp) readArray() (Value, error) {
	var value Value
	value.typ = "array"

	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	value.array = make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		value.array[i] = val
	}
	return value, nil
}

func (r *Resp) Read() (Value, error) {
	var value Value
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch string(_type) {
	case ARRAY:
		value, err = r.readArray()
		if err != nil {
			return Value{}, err
		}
	case BULK:
		value, err = r.readBulk()
		if err != nil {
			return Value{}, err
		}
	default:
		return Value{}, nil
	}
	return value, nil
}
