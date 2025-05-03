package raydium

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/go-enols/gosolana"
	"github.com/go-enols/gosolana/ws"
)

type Client struct {
	RpcClient *rpc.Client // http rpc
	WsClient  *ws.Client  // websocket

	LogApusic []LogApusic // 异步中间件，在获取到日志后，异步处理
	// 价格中间件，在获取到价格之后处理

	LogQuque chan *ws.LogResult // 日志队列
	lock     sync.RWMutex       // 读写锁，确保读写安全
	ctx      context.Context
	cancel   context.CancelFunc
	option   []gosolana.Option
}

func NewClient(ctx context.Context, option ...gosolana.Option) *Client {
	_ctx, cancel := context.WithCancel(ctx)
	opt := gosolana.NewDefaultOption(_ctx, option...)

	c := &Client{
		RpcClient: opt.RpcClient,
		WsClient:  opt.WsClient,
		LogQuque:  make(chan *ws.LogResult, 1000),
		lock:      sync.RWMutex{},
		ctx:       ctx,
		cancel:    cancel,
		option:    option,
	}
	go c.logProcess(ctx)
	return c
}

// 重新链接到客户端
func (c *Client) Reconnect() {
	c.cancel()
	_ctx, cancel := context.WithCancel(c.ctx)
	opt := gosolana.NewDefaultOption(_ctx, c.option...)
	c.RpcClient = opt.RpcClient
	c.WsClient = opt.WsClient
	c.cancel = cancel
	log.Println("WebSocket 连接成功")
}

// 启动监听, 如果需要取消请直接结束ctx
func (c *Client) Start(ctx context.Context, pubKey solana.PublicKey, commit rpc.CommitmentType) (*ws.LogSubscription, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			c.monit(ctx, pubKey, commit)
			c.Reconnect()
		}
	}
}

// 添加一个日志中间件，在获取到日志后，异步处理
func (c *Client) UseLog(apusic LogApusic) {
	// 加锁确保写入安全
	c.lock.Lock()
	defer c.lock.Unlock()
	c.LogApusic = append(c.LogApusic, apusic)
}

// 订阅日志，传输给管道
func (c *Client) monit(ctx context.Context, pubKey solana.PublicKey, commit rpc.CommitmentType) {
	sub, err := c.WsClient.LogsSubscribeMentions(pubKey, commit)
	if err != nil {
		log.Println("订阅失败,3秒后重试 | ", err)
		time.Sleep(3 * time.Second)
		return
	}
	defer sub.Unsubscribe()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := sub.Recv(ctx)
			if err != nil {
				log.Println("WebSocket链接异常,重连中 | ", err)
				return
			}

			if msg.Value.Logs == nil || msg.Value.Err != nil {
				continue // Skip this message.
			}
			c.LogQuque <- msg
		}
	}
}

// 将管道中的日志进行分发
func (c *Client) logProcess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case log := <-c.LogQuque:
			for _, apusic := range c.LogApusic {
				go apusic(log) // 异步处理启动所有
			}
		}
	}
}
