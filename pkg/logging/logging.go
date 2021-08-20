package logging

import (
	"io"
	"log"
	"os"
)

func ErrorHandler(msg string, e error) {
	file, err := os.OpenFile("all.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatalf("file closing fauilure: %v", err)
		}
	}()
	wrt := io.MultiWriter(file)
	log.SetOutput(wrt)
	log.Println(msg, e)
}
