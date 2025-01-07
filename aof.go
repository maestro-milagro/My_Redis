package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mx   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	go func() {
		aof.mx.Lock()

		err = aof.file.Sync()
		if err != nil {
			aof.mx.Unlock()
			return
		}

		aof.mx.Unlock()

		time.Sleep(time.Second)
	}()

	return aof, nil
}

func (a *Aof) Close() error {
	a.mx.Lock()
	defer a.mx.Unlock()

	return a.file.Close()
}

func (a *Aof) Write(value Value) error {
	a.mx.Lock()
	defer a.mx.Unlock()

	_, err := a.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}
