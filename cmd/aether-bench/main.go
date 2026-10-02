package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
	"github.com/affanahmed2302-cmyk/Aether/internal/transport"
)

func main() {
	addr := flag.String("addr", "localhost:7001", "leader or any node address")
	clients := flag.Int("clients", 10, "concurrent clients")
	ops := flag.Int("ops", 1000, "total operations")
	flag.Parse()

	var wg sync.WaitGroup
	var success atomic.Int64
	var failed atomic.Int64
	latencies := make(chan time.Duration, *ops)

	start := time.Now()
	opsPerClient := *ops / *clients

	for c := 0; c < *clients; c++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < opsPerClient; i++ {
				key := fmt.Sprintf("k-%d-%d", id, i)
				cmd := raft.EncodeCommand("put", key, "v")
				t0 := time.Now()
				_, err := transport.Propose(*addr, cmd)
				latencies <- time.Since(t0)
				if err != nil {
					failed.Add(1)
				} else {
					success.Add(1)
				}
			}
		}(c)
	}
	wg.Wait()
	close(latencies)
	elapsed := time.Since(start)

	var total time.Duration
	var count int64
	var max time.Duration
	for d := range latencies {
		total += d
		count++
		if d > max {
			max = d
		}
	}
	avg := time.Duration(0)
	if count > 0 {
		avg = total / time.Duration(count)
	}
	throughput := float64(success.Load()) / elapsed.Seconds()

	log.Printf("=== Aether Benchmark ===")
	log.Printf("Total ops (requested): %d", *ops)
	log.Printf("Success: %d | Failed: %d", success.Load(), failed.Load())
	log.Printf("Elapsed: %v", elapsed)
	log.Printf("Throughput: %.2f ops/sec", throughput)
	log.Printf("Avg latency: %v", avg)
	log.Printf("Max latency: %v", max)
}
