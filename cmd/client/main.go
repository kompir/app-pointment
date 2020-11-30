package main

import (
	"app-pointment/client"
	"flag"
	"fmt"
	"os"
)

var (
	backendURIFlag = flag.String("backend", "http://localhost:8008", "Backend API URL")
	helpFlag       = flag.Bool("help", false, "Display helpful message")
)

func main() {
	flag.Parse()
	s := client.NewSwitch(*backendURIFlag)

	if *helpFlag || len(os.Args) == 1 {
		s.Help()
		return
	}
	err := s.Switch()
	if err != nil {
		fmt.Printf("Cmd switch error: %v\n", err)
		os.Exit(2)
	}
}
