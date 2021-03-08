package util

import (
	"io"
	"log"
)

func CloseWithLogOnErr(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Println(err)
	}
}
