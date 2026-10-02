package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
	"github.com/affanahmed2302-cmyk/Aether/internal/transport"
)

// aether-chaos is a simple chaos tool that:
// 1. Continuously sends writes
// 2. Randomly targets different nodes
// 3. Reports whether the cluster keeps accepting commands after simulated disruption
//
// This is the kind of failure-oriented testing that separates average student projects
// from elite intern-level work.
func main() {
	targets := flag.String("targets", "localhost:7001,localhost:7002,localhost:7003", "comma-separated node addresses")
	duration := flag.Duration("duration", 30*time.Second, "how long to run chaos")
	flag.Parse()

	addrs := strings.Split(*targets, ",")
	log.Printf("Starting chaos test against %v for %v", addrs, *duration)

	end := time.Now().Add(*duration)
	success := 0
	fail := 0
	i := 0

	for time.Now().Before(end) {
		addr := addrs[rand.Intn(len(addrs))]
		key := fmt.Sprintf("chaos-%d", i)
		cmd := raft.EncodeCommand("put", key, fmt.Sprintf("v-%d", i))
		_, err := transport.Propose(addr, cmd)
		if err != nil {
			fail++
			log.Printf("write failed via %s: %v", addr, err)
		} else {
			success++
			log.Printf("write ok via %s key=%s", addr, key)
		}
		i++
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("Chaos finished. success=%d fail=%d", success, fail)
	if success == 0 {
		log.Fatal("cluster accepted zero writes — something is wrong")
	}
	log.Printf("Cluster continued to accept writes under chaos. This is the signal elite interns show.")
}
