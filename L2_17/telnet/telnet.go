package telnet

import (
	"bufio"
	"fmt"
	"mytelnet/config"
	"net"
	"os"
	"sync"
	"time"
)

// Client представляет telnet-клиент
type Client struct {
	conn   net.Conn
	config *config.Config
	done   chan struct{}
	wg     *sync.WaitGroup
}

// New создает новый клиент
func New(config *config.Config) *Client {
	return &Client{
		config: config,
		done:   make(chan struct{}),
		wg:     &sync.WaitGroup{},
	}
}

// Connect устанавливает соединение с сервером
func (tc *Client) Connect() error {
	address := net.JoinHostPort(tc.config.Host, tc.config.Port)

	conn, err := net.DialTimeout("tcp", address, tc.config.Timeout)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	tc.conn = conn
	fmt.Printf("Connected to %s\n", address)

	return nil
}

// Start запускает обработку ввода/вывода
func (tc *Client) Start() {

	tc.wg.Add(1)
	go tc.readFromSocket()

	tc.wg.Add(1)
	go tc.writeToSocket()

	tc.wg.Wait()
	tc.Cleanup()
}

// readFromSocket читает данные из сокета и выводит в STDOUT
func (tc *Client) readFromSocket() {
	defer tc.wg.Done()

	scanner := bufio.NewScanner(tc.conn)

	for {
		select {
		case <-tc.done:
			fmt.Println("Reader: received done signal, shutting down...")
			return
		default:
			tc.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					if ne, ok := err.(net.Error); ok && ne.Timeout() {
						continue
					}
				}
				fmt.Println("Reader: connection closed by server")
				close(tc.done)
				return
			}
			text := scanner.Text()
			fmt.Println(text)
		}
	}
}

// writeToSocket читает данные из STDIN и отправляет в сокет
func (tc *Client) writeToSocket() {
	defer tc.wg.Done()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-tc.done:
			return
		default:
			if !scanner.Scan() {
				fmt.Println("Writer: closing connection...")
				close(tc.done)
				return
			}

			text := scanner.Text() + "\n"
			_, err := tc.conn.Write([]byte(text))
			if err != nil {
				fmt.Printf("Error writing to socket: %v\n", err)
				close(tc.done)
				return
			}
		}
	}
}

// Cleanup закрывает соединение и освобождает ресурсы
func (tc *Client) Cleanup() {
	if tc.conn != nil {
		tc.conn.Close()
	}
}
