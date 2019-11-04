# fstream

*fstream* is the simple solution for collect/read input messages in binary files queue. Each file contains all messages that were received within 1 minute.

File format:

```
<uint16> <msg1> <uint16> <msg2> ...
```

Message size is stored in uint16 value.

All files are saved in the selected directory and have names similar to the following 000000, 000001, 000002, ...
The name 000000 follows name 999999. The current file name is stored in idx file.

Writer saves the last 10000 files and and deletes older.

fstream is useful when there is a large flow of events and there is no way to process them on the fly or use solutions such as RabbitMQ. It's effective for statistics aggregation and subsequent processing.

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
