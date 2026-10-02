package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/affanahmed2302-cmyk/Aether/internal/kv"
	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
	"github.com/affanahmed2302-cmyk/Aether/internal/transport"
)

func main() {
	id := flag.String("id", "node1", "node id")
	addr := flag.String("addr", ":7001", "listen address")
	peers := flag.String("peers", "", "comma-separated peer addresses")
	flag.Parse()

	var peerList []string
	if *peers != "" {
		peerList = strings.Split(*peers, ",")
	}

	log.Printf("Starting Aether node %s on %s | peers: %v", *id, *addr, peerList)

	applyCh := make(chan raft.LogEntry, 256)
	store := kv.NewStore()

	go func() {
		for e := range applyCh {
			store.Apply(e)
			log.Printf("[%s] applied index=%d", *id, e.Index)
		}
	}()

	node := raft.NewNode(*id, peerList, applyCh)
	node.SendRequestVote = transport.SendRequestVote
	node.SendAppendEntries = transport.SendAppendEntries
	node.Start()

	srv := transport.NewServer(node)
	if err := srv.Start(*addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	defer srv.Close()

	log.Printf("[%s] listening on %s", *id, *addr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down...")
	node.Stop()
}
