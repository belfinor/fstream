// +build ignore

package main

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-11-03

import (
	"fmt"
	"os"

	"github.com/belfinor/fstream"
)

func main() {

	os.Mkdir(".data", 0777)

	r := fstream.NewReader(".data", ".data/reader.idx", func(data []byte) { fmt.Println(string(data)) })

	r.BeforeReadFile = func(filename string) {
		fmt.Printf("process file %s\n", filename)
	}

	r.AfterReadFile = func(filename string) {
		fmt.Printf("finish file %s\n", filename)
	}

	defer r.Close()

	wait := make(chan int)

	<-wait
}
