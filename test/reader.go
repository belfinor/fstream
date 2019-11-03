// +build ignore

package main

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-11-03

import (
	"fmt"
	"os"
	"time"

	"github.com/belfinor/fstream"
)

func main() {

	os.Mkdir(".data", 0777)

	w := fstream.NewReader(".data", ".data/reader.idx", func(data []byte) { fmt.Println(string(data)) })
	defer w.Close()

	for {
		<-time.After(time.Second)
	}
}
