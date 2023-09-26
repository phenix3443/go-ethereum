package greeting

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

func TestRPC(t *testing.T) {
	cli, err := rpc.Dial("http://localhost:8545")
	assert.NoError(t, err)
	var msg string
	assert.NoError(t, cli.Call(&msg, "greeting_hello", "world"))
	t.Log(msg)
}

func TestSubscribe(t *testing.T) {
	cli, err := rpc.Dial("ws://localhost:8546")
	assert.NoError(t, err)
	ctx := context.Background()
	ch := make(chan string)
	sub, err := cli.Subscribe(ctx, "greeting", ch, "clock")
	assert.NoError(t, err)
	go func() {
		time.Sleep(5 * time.Second)
		sub.Unsubscribe()
	}()
	for {
		select {
		case v := <-ch:
			fmt.Println(v)
		case err := <-sub.Err():
			fmt.Println(err)
			return
		}
	}
}
