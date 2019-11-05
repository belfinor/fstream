package fstream

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.002
// @date    2019-11-04

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/belfinor/fcounter"
	"github.com/belfinor/log"
	"github.com/belfinor/pack"
)

type Reader struct {
	path     string
	handler  func([]byte)
	cnt      *fcounter.Counter
	data     []byte
	cancel   context.CancelFunc
	finished chan bool
}

func NewReader(path string, idx string, handler func([]byte)) *Reader {

	r := &Reader{
		path:     path,
		cnt:      fcounter.New(idx, FileNumberMod),
		handler:  handler,
		data:     []byte{},
		finished: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())

	r.cancel = cancel

	go r.reader(ctx)

	return r
}

func (r *Reader) Close() {
	r.cancel()
	<-r.finished
}

func (r *Reader) reader(ctx context.Context) {

	for {

		filename := fmt.Sprintf("%s/%06d", r.path, r.cnt.Get())

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Info("wait file " + filename)

			for i := 0; i < 60; i++ {

				select {

				case <-ctx.Done():
					close(r.finished)
					return

				case <-time.After(time.Second):

				}

			}

			continue
		}

		r.ReadFile(filename)

		r.cnt.Inc()
	}

}

func (r *Reader) ReadFile(filename string) {

	log.Info("process file " + filename)

	f, err := os.Open(filename)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	buffer := make([]byte, 10240)

	for {

		n, err := br.Read(buffer)
		if err != nil && n == 0 {
			break
		}

		r.onData(buffer[:n])
	}
}

func (r *Reader) onData(data []byte) {
	r.data = bytes.Join([][]byte{r.data, data}, []byte{})
	size := uint16(0)

	list := r.data

	for len(list) > 2 {
		if pack.Decode(list, &size) != nil {
			break
		}
		size = size + 2
		if len(list) > int(size) {
			r.handler(list[2:size])
			list = list[size:]
		} else if len(list) == int(size) {
			r.handler(list[2:])
			list = []byte{}
		} else {
			break
		}
	}

	if len(list) > 0 {
		r.data = list
	} else {
		r.data = []byte{}
	}
}
