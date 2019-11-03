// +build ignore

package main

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-11-03

import (
	"bufio"
	"os"
	"strings"

	"github.com/belfinor/fstream"
)

func main() {

	os.Mkdir(".data", 0777)

	w := fstream.NewWriter(".data", ".data/writer.idx")
	defer w.Close()

	br := bufio.NewReader(os.Stdin)

	for {

		str, err := br.ReadString('\n')
		if err != nil && str == "" {
			break
		}

		str = strings.TrimSpace(str)

		if str == "" {
			continue
		}

		w.Write([]byte(str))

	}
}
