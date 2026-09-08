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

	t.wg.Add(1)
	//как это все завершается? ф-я видит EOF, прекращает читать in и спокойно завершает свою работу
	go func() {
		defer t.wg.Done()
		io.Copy(conn, t.in) //in -> conn
	}()

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		io.Copy(t.out, conn) //conn -> out
	}()

	//если тут написать wg.Wait -> все заблокируется
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
	return nil
}

func (t *telnetClient) Receive() error {
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
