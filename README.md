# fstream

Simple solution for collect/read input messages in binary files queue. Each file contains all messages received in 1 minute.

File format:

```
<uint16> <msg1> <uint16> <msg2> ...
```

Message size is stored in uint16 value.

fstream is useful when there is a large flow of events and there is no way to process them on the fly or use solutions such as RabbitMQ.

## Writer example

```go
package main

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
```

## Reader example

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/belfinor/fstream"
)

func handler(data []byte) {
  fmt.Println(string(data)
}


func main() {

	os.Mkdir(".data", 0777)

	w := fstream.NewReader(".data", ".data/reader.idx", handler)
	defer w.Close()

  wait := make(chan int)

  <-wait
}
```
