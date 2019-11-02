package writer

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-11-02

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/belfinor/fcounter"
	"github.com/belfinor/log"
	"github.com/belfinor/pack"
)

const (
	MessageLimit  int   = 2048
	FileNumberMod int64 = 1000000
	SaveFiles     int64 = 10000
	SavePeriod    int64 = 60
)

type Writer struct {
	Input    chan []byte
	File     *os.File
	Path     string
	Cnt      *fcounter.Counter
	cancel   context.CancelFunc
	finished chan bool
}

func Open(path string, idx string) *Writer {

	w := &Writer{
		Input:    make(chan []byte, MessageLimit),
		Path:     path,
		Cnt:      fcounter.New(idx, FileNumberMod),
		finished: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())

	w.cancel = cancel

	go w.writer(ctx)

	return w
}

func (w *Writer) Close() {
	w.cancel()
	<-w.finished
}

func (w *Writer) openLog() {

	file_name := fmt.Sprintf("%s/%d.tmp", w.Path, w.Cnt.Get())
	var err error

	if w.File != nil {
		w.File.Close()
		last_name := fmt.Sprintf("%s/%d", w.Path, (w.Cnt.Get()+FileNumberMod-1)%FileNumberMod)
		os.Rename(last_name+".tmp", last_name)
	}

	if w.File, err = os.OpenFile(file_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	log.Info("start log " + file_name)
}

func (w *Writer) rotate() {

	w.Cnt.Inc()

	w.openLog()

	remove_name := fmt.Sprintf("%s/%d", w.Path, (w.Cnt.Get()+FileNumberMod-SaveFiles)%FileNumberMod)
	os.Remove(remove_name)
}

func (w *Writer) Write(data []byte) (int, error) {

	size := len(data)

	if size == 0 {
		return 0, errors.New("no input data")
	}

	block := make([]byte, size)
	copy(block, data)

	w.Input <- block

	return size, nil
}

func (w *Writer) writer(ctx context.Context) {
	log.Info("start storage writer")

	w.openLog()

	start := time.Now().Unix()
	last := start

	for {
		select {

		case data := <-w.Input:

			size := len(data)
			if size > 0 {

				data = pack.Encode(uint16(size), data)

				if _, err := w.File.Write(data); err != nil {
					log.Fatal(err.Error())
				}
			}

		case <-time.After(time.Second):

		case <-ctx.Done():

			if w.File != nil {
				w.File.Close()
				last_name := fmt.Sprintf("%s/%d", w.Path, w.Cnt.Get())
				os.Rename(last_name+".tmp", last_name)
			}

			w.Cnt.Inc()

			w.finished <- true

			return

		}

		last = time.Now().Unix()

		if last-start >= SavePeriod {
			w.rotate()
			start = last
		}
	}
}