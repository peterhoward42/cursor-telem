package main

import (
	"log"
	"os"

	function "cursor-telem"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	if err := function.Register(); err != nil {
		log.Fatalf("function.Register: %v", err)
	}

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	hostname := ""
	if os.Getenv("LOCAL_ONLY") == "true" {
		hostname = "127.0.0.1"
	}

	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v", err)
	}
}
