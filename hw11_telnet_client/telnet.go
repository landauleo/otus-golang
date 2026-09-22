package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

/*
Пользователь нажал Ctrl+D
 1. io.Copy(conn, in) завершается
 2. Вызывается defer conn.Close()  <── (Закрываем сокет!)
 3. Вторая горутина, сидевшая в io.Copy(out, conn), пытается прочитать из сокета, но видит, что он ЗАКРЫТ.
 4. Вторая io.Copy моментально аварийно завершается!
 5. Срабатывают оба wg.Done(), wg.Wait() отпускает главный поток.
*/

// интерфейс-договор
type TelnetClient interface {
	Connect() error
	io.Closer //Interface Embedding
	Send() error
	Receive() error
}

// реализация
type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
	wg      sync.WaitGroup
}

func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout) //получаем сокет, созданный ОС, в GO он оборачивается в conn
	if err != nil {
		return fmt.Errorf("failed to create connection: %w", err)
	}

	t.conn = conn //!!!
	return nil
}

func (t *telnetClient) Close() error {
	if t.conn != nil {
		conErr := t.conn.Close()
		if conErr != nil {
			return fmt.Errorf("failed to close connection: %w", conErr)
		}
	}

	if t.in != nil {
		inErr := t.in.Close()
		if inErr != nil {
			return fmt.Errorf("failed to close input stream: %w", inErr)
		}
	}

	return nil
}

func (t *telnetClient) Send() error {
	if t.conn == nil {
		return fmt.Errorf("connection is not established")
	}

	_, err := io.Copy(t.conn, t.in) //в Copy прячется бесконечный цикл
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to send: %w", err)
	}
	return nil
}

func (t *telnetClient) Receive() error {
	if t.conn == nil {
		return fmt.Errorf("connection is not established")
	}

	_, err := io.Copy(t.out, t.conn)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to receive: %w", err)
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
