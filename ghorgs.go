// Copyright (c) 2019 Sony Mobile Communications Inc.
// All rights reserved.

package main

import (
	"ghorgs/cmd"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	cmd.Execute()
}
