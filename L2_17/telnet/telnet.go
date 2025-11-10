package telnet

import (
	"mytelnet/config"
	"net"
	"sync"
)

// Client представляет telnet-клиент
type Client struct {
	conn   net.Conn
	config *config.Config
	done   chan struct{}
	wg     sync.WaitGroup
}

func NewTelnet(config *config.Config) *Client {

}
