package main

import (
	"os"
	"log"
)

type Entry struct {
	name string
}

type List []Entry

func GetList(path string) List {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}

	var l List

	for _, file := range files {
		entry := Entry {
			name: file.Name(),
		}
		l = append(l, entry)
	}

	return l
}

func ListView(l List) string {
	var s string
	for _, entry := range l {
		s += entry.name + "\n"
	}
	return s
}
