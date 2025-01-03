package main

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
