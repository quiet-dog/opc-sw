package global

import (
	"context"
	"fmt"
	"os"
	"sw/config"
	"sw/opc"
	"sync"

	"github.com/lxzan/gws"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DEVICEDATA = 0
	YUZHI      = 1
	BAOJING    = 2
	ANIMATION  = 3
)

type RecHandler struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

var (
	DB         *gorm.DB
	Config     config.Config
	OpcGateway = opc.New()
	Ctx        = context.Background()
	Redis      *redis.Client
	Upgrader   *gws.Upgrader
	Session    = sync.Map{}
	RecChanel  = make(chan RecHandler, 5)
	Handler    = &HandlerChanel{
		clients: sync.Map{},
	}
)

type HandlerChanel struct {
	clients sync.Map // 存储客户端连接
}

// 注册客户端通道
func (h *HandlerChanel) RegisterClientChannel() chan RecHandler {
	ch := make(chan RecHandler)
	h.clients.Store(ch, nil) // 将通道存储到 Handler 的 clients map 中
	return ch
}

func (h *HandlerChanel) UnregisterClientChannel(ch chan RecHandler) {
	if _, ok := h.clients.Load(ch); ok {
		h.clients.Delete(ch) // 从 Handler 的 clients map 中删除通道
		close(ch)            // 关闭通道
	}
}

func (h *HandlerChanel) Start() {
	// 启动一个 goroutine 来处理 RecChanel 中的消息
	go func() {
		ctx := context.Background()
		for {
			select {
			case <-ctx.Done():
				return
			case rec := <-RecChanel:
				f, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				defer f.Close()
				h.clients.Range(func(key, value interface{}) bool {
					ch := key.(chan RecHandler)
					select {
					case ch <- rec: // 将消息发送到每个客户端通道
						fmt.Println("senddata3333:", rec)
						f.Write([]byte(fmt.Sprintf("send data: %v\n", rec)))
					default:
						// 如果通道满了，可以选择忽略或处理
						fmt.Println("senddata4444:")
					}
					return true // 继续遍历所有客户端通道
				})
			}
		}
	}()
}
