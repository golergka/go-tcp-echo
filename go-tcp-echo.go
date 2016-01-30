package main

import (
    "net"
    "log"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

func main() {
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        log.Panicln(err)
    }
    defer l.Close()

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Panicln(err)
        }

        go handleRequest(conn)
    }
}

func handleRequest(conn net.Conn) {
    defer conn.Close()

    for {
        buf := make([]byte, 1024)
        _, err := conn.Read(buf)
        if err != nil {
            return;
        }

        conn.Write([]byte("Message received."))
    }
}
