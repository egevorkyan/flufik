package logging

import (
	"github.com/egevorkyan/flufik/core"
	"io"
	"log"
	"os"
)

func ErrorHandler(msg string, e error) {
	file, err := os.OpenFile(core.FlufikLoggingFilePath(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatalf("file closing fauilure: %v", err)
		}
	}()
	wrt := io.MultiWriter(os.Stdout, file)
	log.SetOutput(wrt)
	log.Println(msg, e)
}
