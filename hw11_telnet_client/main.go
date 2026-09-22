package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	timeoutPointer := flag.Duration("timeout", 10*time.Second, "timeout for connection")
	flag.Parse() //!!!

	args := flag.Args()
	if len(args) < 2 {
		//под капотом os.Exit(1), используем только в консольных утилитах
		log.Fatalf("Usage: go-telnet [--timeout=10s] host port\n")
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeoutPointer, os.Stdin, os.Stdout)
	defer client.Close()

	if err := client.Connect(); err != nil {
		log.Fatalf("connection error: %v", err)
	}

	//обработка системных сигналов
	done := make(chan struct{}, 2)

	go func() {
		_ = client.Send()
		done <- struct{}{} //сигнал "финиша", от каждой горутины запишется ровно 1 раз
	}()

	go func() {
		_ = client.Receive()
		done <- struct{}{}
	}()
}
