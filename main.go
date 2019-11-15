package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
)

func main() {
	// parse command line
	flag.Parse()
	files, err := getFiles(flag.Args())
	if err != nil {
		log.Fatalln("Errors inspecting arguments, aborting: " + err.Error())
	}
	// read and inflate files
	for _, file := range files {
		err = processFile(file)
		if err != nil {
			log.Printf("File \"%s\" failed: %s", file, err.Error())
		}
	}
}

func processFile(file string) (err error) {
	fi, err := os.Open(file)
	if err != nil {
		log.Fatalln("Could not open file: " + err.Error())
	}
	defer fi.Close()
	log.Printf("Processing \"%s\"", file)
	gzipReader, err := gzip.NewReader(fi)
	if err != nil {
		log.Fatalln("Gzip Error: " + err.Error())
	}
	tarReader := tar.NewReader(gzipReader)
	if err != nil {
		log.Fatalln("Tar Error: " + err.Error())
	}
	for {
		tarHeader, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				// we reached the end of the tar archive
				return nil
			}
			log.Fatalln("Tar Next() Error: " + err.Error())
		}
		if !tarHeader.FileInfo().IsDir() {
			// we have a file
			msg, err := mail.ReadMessage(tarReader)
			if err != nil {
				log.Println("ReadMessage Error: " + err.Error())
				continue
			}
			date, err := msg.Header.Date()
			if err != nil {
				log.Println("Header.Date Error: " + err.Error())
				continue
			}
			fmt.Printf("Date Sent: %s\n", date)
			fmt.Printf("Sender: %s\n", msg.Header.Get("From"))
			fmt.Printf("Subject: %s\n\n", msg.Header.Get("Subject"))
		}
	}
}

// Validate that the files in the list exist
func getFiles(fileList []string) (files []string, err error) {
	var fileInfo os.FileInfo
	for _, file := range fileList {
		fileInfo, err = os.Stat(file)
		if err != nil {
			log.Println(err.Error())
		} else if fileInfo.IsDir() {
			err = fmt.Errorf("%s is a Directory", file)
			log.Println(err.Error())
		}
		files = append(files, file)
	}
	return files, err
}
