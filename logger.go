package zerodown

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stderr, "[zerodown]", log.Llongfile|log.Lmicroseconds|log.Ldate)
}
