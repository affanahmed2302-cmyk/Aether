package client

import (
	"encoding/json"

	"github.com/affanahmed2302-cmyk/Aether/internal/kv"
)

// Client is a thin wrapper that talks to an Aether cluster.
// In the full version this would discover the current leader and retry on redirect.
type Client struct {
	leaderAddr string
}

func New(leaderAddr string) *Client {
	return &Client{leaderAddr: leaderAddr}
}

func (c *Client) Put(key, value string) error {
	cmd := kv.Command{Op: "put", Key: key, Value: value}
	_, err := json.Marshal(cmd)
	// In a complete implementation we would send this over the network to the leader.
	return err
}

func (c *Client) Get(key string) (string, error) {
	// Placeholder — real client would issue a read RPC.
	return "", nil
}

func (c *Client) Delete(key string) error {
	cmd := kv.Command{Op: "delete", Key: key}
	_, err := json.Marshal(cmd)
	return err
}
