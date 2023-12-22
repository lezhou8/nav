package main

import (
	"os"
	"log"
)

func GetFilePath() string {
	d, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return d
}
