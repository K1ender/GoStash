package handler

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func startTestServer(tb testing.TB) (addr string, stop func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatalf("failed to listen: %v", err)
	}
	h := NewHandler()
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				h.Handle(c)
			}(conn)
		}
	}()
	return ln.Addr().String(), func() { ln.Close() }
}

func BenchmarkSocketGetHandler(b *testing.B) {
	addr, stop := startTestServer(b)
	defer stop()
	time.Sleep(50 * time.Millisecond)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		b.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	cmd := []byte("GET\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")
	for i := 0; i < b.N; i++ {
		conn.Write(cmd)
		buf := make([]byte, 128)
		conn.Read(buf)
	}
}

func BenchmarkSocketSetHandler(b *testing.B) {
	addr, stop := startTestServer(b)
	defer stop()
	time.Sleep(50 * time.Millisecond)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		b.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	cmd := []byte("SET\x00" + strconv.Itoa(len("foo")) + "\x00foo\x00" + strconv.Itoa(len("bar")) + "\x00bar\r\n")
	for i := 0; i < b.N; i++ {
		conn.Write(cmd)
		buf := make([]byte, 128)
		conn.Read(buf)
	}
}

func BenchmarkSocketIncrHandler(b *testing.B) {
	addr, stop := startTestServer(b)
	defer stop()
	time.Sleep(50 * time.Millisecond)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		b.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	// сначала установим foo=0
	setCmd := []byte("SET\x00" + strconv.Itoa(len("foo")) + "\x00foo\x00" + strconv.Itoa(len("0")) + "\x000\r\n")
	conn.Write(setCmd)
	buf := make([]byte, 128)
	conn.Read(buf)
	cmd := []byte("INC\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")
	for i := 0; i < b.N; i++ {
		conn.Write(cmd)
		conn.Read(buf)
	}
}

func BenchmarkSocketDecrHandler(b *testing.B) {
	addr, stop := startTestServer(b)
	defer stop()
	time.Sleep(50 * time.Millisecond)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		b.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	// сначала установим foo=0
	setCmd := []byte("SET\x00" + strconv.Itoa(len("foo")) + "\x00foo\x00" + strconv.Itoa(len("0")) + "\x000\r\n")
	conn.Write(setCmd)
	buf := make([]byte, 128)
	conn.Read(buf)
	cmd := []byte("DEC\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")
	for i := 0; i < b.N; i++ {
		conn.Write(cmd)
		conn.Read(buf)
	}
}

func BenchmarkSocketDelHandler(b *testing.B) {
	addr, stop := startTestServer(b)
	defer stop()
	time.Sleep(50 * time.Millisecond)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		b.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	setCmd := []byte("SET\x00" + strconv.Itoa(len("foo")) + "\x00foo\x00" + strconv.Itoa(len("bar")) + "\x00bar\r\n")
	cmd := []byte("DEL\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")
	buf := make([]byte, 128)
	for i := 0; i < b.N; i++ {
		conn.Write(setCmd)
		conn.Read(buf)
		conn.Write(cmd)
		conn.Read(buf)
	}
}
