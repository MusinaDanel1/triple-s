package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"triple-s/pkg/server"
)

func main() {
	port := flag.String("port", "8080", "Port number")
	dir := flag.String("dir", "data", "Path to the directory")
	help := flag.Bool("help", false, "Show this help message")
	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	portNum, err := server.ValidatePort(*port)
	if err != nil {
		log.Fatalf("Error: %v/n", err)
	}

	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		err = os.Mkdir(*dir, 0o755)
		if err != nil {
			log.Fatalf("error creating data directory: %v", err)
		}
		fmt.Printf("Directoty %s created.\n", *dir)
	}

	fmt.Printf("Starting server on port %v\n", portNum)
	fmt.Printf("Using directory: %s\n", *dir)

	if err := http.ListenAndServe(":"+(*port), server.SetupRoutes(*dir)); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

func showHelp() {
	helpMessage := `Simple Storage Service.
	
	
**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
  --help     Show this screen.
  --port N   Port number
  --dir S    Path to the directory`

	fmt.Println(helpMessage)
}
