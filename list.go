package main

import (
	"os"
	"log"
)

type Entry struct {
	name string
}

type List struct {
	entries []Entry
	cursor int
}

func GetList(path string) List {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Panic(err)
	}

	l := List {
		entries: []Entry{},
		cursor: 0,
	}

	for _, file := range files {
		entry := Entry {
			name: file.Name(),
		}
		l.entries = append(l.entries, entry)
	}

	return l
}

func ListView(l List) string {
	var s string
	for i, entry := range l.entries {
		if i == l.cursor {
			s += "> " + entry.name + "\n"
		} else {
			s += "  " + entry.name + "\n"
		}
	}
	return s
}

func (l List) ListGoUp() List {
	newCursor := l.cursor - 1

	if newCursor < 0 {
		return l
	}

	return List {
		entries: l.entries,
		cursor: newCursor,
	}
}

func (l List) ListGoDown() List {
	newCursor := l.cursor + 1

	if newCursor >= len(l.entries) {
		return l
	}

	return List {
		entries: l.entries,
		cursor: newCursor,
	}
}
