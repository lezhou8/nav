package main

import (
	"os"
	"log"
)

type FilePath string

func GetFilePath() FilePath {
	d, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return FilePath(d)
}
