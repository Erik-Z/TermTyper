package main

import (
	"log"
	"termtyper/cmd"
)

func main() {

	cmd.OsInit()
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
