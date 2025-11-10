package telnet

import (
	"bufio"
	"fmt"
	"io"
	"mytelnet/config"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	//tc.setupSignalHandler()

	tc.wg.Add(1)
	go tc.readFromSocket()

	tc.wg.Add(1)
	go tc.writeToSocket()

	tc.wg.Wait()
	tc.Cleanup()
}

// setupSignalHandler настраивает обработчик сигналов
func (tc *Client) setupSignalHandler() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nReceived interrupt signal, closing connection...")
		close(tc.done)
	}()
}

// readFromSocket читает данные из сокета и выводит в STDOUT
func (tc *Client) readFromSocket() {
	defer tc.wg.Done()

	lines := make(chan string)
	errors := make(chan error)
	scanner := bufio.NewScanner(tc.conn)

	go func() {

		for scanner.Scan() {
			lines <- scanner.Text() + "\n"
		}
		if err := scanner.Err(); err != nil {
			errors <- err
		} else {
			errors <- io.EOF
		}
	}()

	for {
		select {
		case <-tc.done:
			return
		default:
			if !scanner.Scan() {
				fmt.Println("\nConnection closed by server")
				close(tc.done)
				return
			}
			text := scanner.Text()
			fmt.Println(text)
		}
	}
}

//// readFromSocket читает данные из сокета и выводит в STDOUT
//func (tc *Client) readFromSocket() {
//	defer tc.wg.Done()
//
//	scanner := bufio.NewScanner(tc.conn)
//	for {
//		select {
//		case <-tc.done:
//			return
//		default:
//			if !scanner.Scan() {
//				fmt.Println("\nConnection closed by server")
//				close(tc.done)
//				return
//			}
//			text := scanner.Text()
//			fmt.Println(text)
//		}
//	}
//}

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
				fmt.Println("\nClosing connection...")
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
