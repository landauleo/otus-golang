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
	var err error
	if t.conn != nil {
		conErr := t.conn.Close()
		if conErr != nil {
			err = fmt.Errorf("failed to close connection: %w", conErr)
		}
	}

	if t.in != nil {
		inErr := t.in.Close()
		if inErr != nil {
			err = fmt.Errorf("failed to close input stream: %w", inErr)
		}
	}

	t.wg.Wait() //чтобы горутины подчистили ресурсы

	return err
}

func (t *telnetClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("failed to send: %w", err)
	} //in -> conn
	return nil
}

func (t *telnetClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("failed to receive: %w", err)
	} //conn -> out
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
