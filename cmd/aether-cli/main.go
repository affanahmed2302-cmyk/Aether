package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/affanahmed2302-cmyk/Aether/internal/raft"
	"github.com/affanahmed2302-cmyk/Aether/internal/transport"
)

func main() {
	addr := flag.String("addr", "localhost:7001", "any node address (preferably leader)")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: aether-cli -addr host:port put|get|delete ...")
		os.Exit(1)
	}

	switch args[0] {
	case "put":
		if len(args) < 3 {
			fmt.Println("usage: put <key> <value>")
			os.Exit(1)
		}
		cmd := raft.EncodeCommand("put", args[1], args[2])
		idx, err := transport.Propose(*addr, cmd)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Printf("OK (index %d)\n", idx)
	case "get":
		fmt.Println("get is applied on the state machine — check node logs for applied entries in this educational version")
	case "delete":
		if len(args) < 2 {
			fmt.Println("usage: delete <key>")
			os.Exit(1)
		}
		cmd := raft.EncodeCommand("delete", args[1], "")
		idx, err := transport.Propose(*addr, cmd)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		fmt.Printf("OK (index %d)\n", idx)
	default:
		fmt.Println("unknown command")
		os.Exit(1)
	}
}
