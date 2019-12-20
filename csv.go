// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

import (
	"log"
	"os"
	"reflect"
)

type Csv struct {
	FileName string
	Records  map[string][]string
	Keys     []string
}

func makeCsv(filename string) *Csv {
	keys := make([]string, 0)
	return &Csv{filename, nil, keys}
}

func (c *Csv) flush(title []string) {
	var f *os.File

	if _, err := os.Stat(c.FileName); err == nil || os.IsNotExist(err) {
		f, err = os.OpenFile(c.FileName,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	} else {
		panic(err)
	}

	s := lineToString(title)
	if _, err := f.WriteString(s + "\n"); err != nil {
		panic(err)
	}

	for _, k := range c.Keys {
		if k == title[0] && reflect.DeepEqual(c.Records[k], title[1:]) {
			continue
		}

		s = k + "\t" + lineToString(c.Records[k])
		if _, err := f.WriteString(s + "\n"); err != nil {
			panic(err)
		}
	}
}

func lineToString(line []string) string {
	s := ""
	for i, cell := range line {
		s = s + cell
		if i < len(line)-1 {
			s = s + "\t"
		}
	}

	return s
}

func (c *Csv) log() {
	var s string
	for id, line := range c.Records {
		s = id + "\t"
		for i, cell := range line {
			s = s + cell
			if i < len(line)-1 {
				s = s + "\t"
			}
		}
		log.Print(s + "\n")
	}
}

func (c *Csv) addKey(key string) {
	c.Keys = append(c.Keys, key)
}