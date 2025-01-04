package main

import (
	"bufio"
	"encoding/json"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	INTEGER = ':'
	BULK    = '$'
	ERROR   = '-'
	ARRAY   = '*'
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

func NewResp(reader io.Reader) *Resp {
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
	switch _type {
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

// TODO: move to separate package
// Write RESP

func (v Value) marshalString() []byte {
	var result []byte
	result = append(result, STRING)
	result = append(result, v.str...)
	result = append(result, '\r', '\n')

	return result
}

func (v Value) marshalInteger() []byte {
	var result []byte

	result = append(result, INTEGER)
	result = append(result, strconv.Itoa(v.num)...)
	result = append(result, '\r', '\n')

	return result
}

func (v Value) marshalBulk() []byte {
	var result []byte

	result = append(result, BULK)
	result = append(result, strconv.Itoa(len(v.bulk))...)
	result = append(result, '\r', '\n')
	result = append(result, v.bulk...)
	result = append(result, '\r', '\n')

	return result
}

func (v Value) marshalArray() []byte {
	var result []byte

	result = append(result, ARRAY)
	result = append(result, strconv.Itoa(len(v.array))...)
	result = append(result, '\r', '\n')

	for _, value := range v.array {
		result = append(result, value.Marshal()...)
	}

	return result
}

func (v Value) marshalError() []byte {
	var result []byte
	result = append(result, ERROR)
	result = append(result, v.str...)
	result = append(result, '\r', '\n')

	return result
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "string":
		return v.marshalString()
	case "integer":
		return v.marshalInteger()
	case "bulk":
		return v.marshalBulk()
	case "array":
		return v.marshalArray()
	case "error":
		return v.marshalError()
	case "null":
		return v.marshalNull()
	default:
		return nil
	}
}

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w *Writer) Write(value Value) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
