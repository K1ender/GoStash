package handler

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/k1ender/go-stash/internal/store"
)

func startTestServer(tb testing.TB) (addr string, stop func()) {
	tb.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatalf("failed to listen: %v", err)
	}

	store := store.NewShardedStore(0)
	h := NewHandler(store)

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
	
	for b.Loop() {
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

	
	for b.Loop() {
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

	setCmd := []byte("SET\x00" + strconv.Itoa(len("foo")) + "\x00foo\x00" + strconv.Itoa(len("0")) + "\x000\r\n")

	conn.Write(setCmd)
	buf := make([]byte, 128)
	conn.Read(buf)

	cmd := []byte("INC\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")

	
	for b.Loop() {
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

	setCmd := []byte("SET\x00" + strconv.Itoa(len("foo")) + "\x00foo\x00" + strconv.Itoa(len("0")) + "\x000\r\n")

	conn.Write(setCmd)
	buf := make([]byte, 128)
	conn.Read(buf)

	cmd := []byte("DEC\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")

	
	for b.Loop() {
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

	
	for b.Loop() {
		conn.Write(setCmd)
		conn.Read(buf)
		conn.Write(cmd)
		conn.Read(buf)
	}
}

func BenchmarkSocketRandomKeyInserts(b *testing.B) {
	addr, stop := startTestServer(b)
	defer stop()

	time.Sleep(50 * time.Millisecond)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		b.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 128)

	
	for i := 0; b.Loop(); i++ {
		key := fmt.Sprintf("key_%d_%d", i, rand.Intn(10000))
		value := fmt.Sprintf("value_%d", rand.Intn(1000))

		cmd := []byte("SET\x00" + strconv.Itoa(len(key)) + "\x00" + key + "\x00" + strconv.Itoa(len(value)) + "\x00" + value + "\r\n")

		conn.Write(cmd)
		conn.Read(buf)
	}
}
