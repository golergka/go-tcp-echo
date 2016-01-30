package main

import (
    "net"
    "log"
    "flag"
    "strconv"
)

const (
    CONN_TYPE = "tcp"
)

func main() {
    host := flag.String("host", "localhost", "The host to run echo serve on.")
    port := flag.Int("port", 3333, "Port to accept connections on.")
    flag.Parse();

    url := *host + ":" + strconv.Itoa(*port)
    l, err := net.Listen(CONN_TYPE, url)
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
        size, err := conn.Read(buf)
        if err != nil {
            return;
        }
        conn.Write(buf[:size])
    }
}
