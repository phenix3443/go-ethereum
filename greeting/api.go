package greeting

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type GreetingAPI struct{}

// Hello is a function that returns a greeting message.
func (a *GreetingAPI) Hello(name string) string {
	return "hello" + name
}

// Clock returns a subscription to receive clock updates.
func (a *GreetingAPI) Clock(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		tick := time.NewTicker(time.Second)

		for {
			select {
			case <-tick.C:
				notifier.Notify(rpcSub.ID, time.Now().Local().Format("2006-01-02 15:04:05"))
			case <-rpcSub.Err():
				fmt.Println("client unsubscribed")
				return
			case <-notifier.Closed():
				fmt.Println("notifier closed")
				return
			}
		}
	}()
	return rpcSub, nil
}

// APIs return the collection of RPC services the greeting package offers.
func APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "greeting",
			Service:   &GreetingAPI{},
		},
	}
}
