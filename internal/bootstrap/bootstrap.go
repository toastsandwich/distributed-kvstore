package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/toastsandwich/kvstore/pkg/proto"
	"github.com/toastsandwich/kvstore/pkg/retrymanager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type (
	Nodes []proto.ServiceKVStoreClient
	node  []*grpc.ClientConn
)

var cc node

func withSignalCancel(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	return ctx, cancel
}

func Init(hosts []string) Nodes {
	cc = make(node, 0)

	nodes := make([]proto.ServiceKVStoreClient, 0)
	rm := retrymanager.New(retrymanager.Opts{
		RetryAfter: 1 * time.Second,
	})
	ctx, cancelTO := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelTO()

	ctx, cancelSig := withSignalCancel(ctx)
	defer cancelSig()

	// online will account which node have success fully been connected to
	online := make(map[string]struct{})

	rm.DoWithCtx(ctx, func() bool {
		body := make(map[string]any)
		for _, n := range hosts {
			var (
				err           error
				clientErr     string
				clientNode    string
				node          proto.ServiceKVStoreClient
				ok            bool
				gc            *grpc.ClientConn
				insecurecreds = insecure.NewCredentials()
			)

			if _, ok := online[n]; ok {
				continue
			}
			c := fiber.AcquireClient()

			url := n + "/endpoint"
			a := c.Get(url)
			code, b, errs := a.Bytes()
			if len(errs) > 0 || code != fiber.StatusFound {
				for _, e := range errs {
					fmt.Println(code, e)
				}
				goto END
			}

			fiber.ReleaseAgent(a)
			fiber.ReleaseClient(c)

			if err := json.Unmarshal(b, &body); err != nil {
				fmt.Println("Marshal error:", err)
				goto END
			}

			clientErr, ok = body["error"].(string)
			if ok {
				fmt.Println(clientErr)
				goto END
			}

			clientNode, ok = body["endpoint"].(string)
			if !ok || clientNode == "" {
				fmt.Println("problem with endpoint!!")
				goto END
			}

			gc, err = grpc.NewClient(clientNode, grpc.WithTransportCredentials(insecurecreds))
			if err != nil {
				fmt.Println(err)
				goto END
			}

			cc = append(cc, gc)
			node = proto.NewServiceKVStoreClient(gc)
			nodes = append(nodes, node)

			online[n] = struct{}{}

		END:
			fmt.Printf("online %d/%d\n", len(nodes), len(online))
		}
		return len(online) == len(hosts)
	})

	return nodes
}

func Close() {
	if cc == nil {
		return
	}
	for _, conn := range cc {
		conn.Close()
	}
}
