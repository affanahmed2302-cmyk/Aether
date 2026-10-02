package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/affanahmed2302-cmyk/Aether/internal/kv"
	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
)

func main() {
	id := flag.String("id", "node1", "unique node id")
	port := flag.String("port", "7001", "listen port")
	flag.Parse()

	log.Printf("Starting Aether node %s on port %s", *id, *port)

	applyCh := make(chan raft.LogEntry, 100)
	store := kv.NewStore()

	// Apply committed entries to the state machine
	go func() {
		for entry := range applyCh {
			store.Apply(entry)
			log.Printf("[%s] applied index %d", *id, entry.Index)
		}
	}()

	node := raft.NewNode(*id, []string{}, applyCh)
	node.Start()

	// Wait for interrupt
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("Shutting down...")
	node.Stop()
}
