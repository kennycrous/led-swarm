package main

import (
	"flag"
	"os"
	"testing"
)

func TestMain_Flags(t *testing.T) {
	// Test flag definition parsing
	origArgs := os.Args
	t.Cleanup(func() {
		os.Args = origArgs
	})

	os.Args = []string{"led-swarm", "--server", "--port=9090", "--db=test-run.db"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	serverMode := flag.Bool("server", false, "Run in server mode")
	port := flag.Int("port", 8080, "Port to listen on")
	dbPath := flag.String("db", "led-swarm.db", "DB path")

	flag.Parse()

	if !*serverMode {
		t.Errorf("Expected serverMode=true")
	}
	if *port != 9090 {
		t.Errorf("Expected port=9090, got %d", *port)
	}
	if *dbPath != "test-run.db" {
		t.Errorf("Expected dbPath='test-run.db', got '%s'", *dbPath)
	}
}
