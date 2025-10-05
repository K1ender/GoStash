package handler

import (
	"strconv"
	"testing"

	"github.com/k1ender/go-stash/internal/store"
)

var s = store.NewShardedStore(0)

func BenchmarkGetHandler(b *testing.B) {
	s.Set("foo", "bar")

	h := NewGetHandler(s)
	cmd := []byte("GET\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")

	for b.Loop() {
		h.Handle(cmd)
	}
}

func BenchmarkSetHandler(b *testing.B) {
	h := NewSetHandler(s)
	cmd := []byte("SET\x00" + strconv.Itoa(len("foo")) + "\x00foo\x00" + strconv.Itoa(len("bar")) + "\x00bar\r\n")

	for b.Loop() {
		h.Handle(cmd)
	}
}

func BenchmarkIncrHandler(b *testing.B) {
	h := NewIncrHandler(s)
	s.Set("foo", "0")
	cmd := []byte("INC\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")

	for b.Loop() {
		h.Handle(cmd)
	}
}

func BenchmarkDecrHandler(b *testing.B) {
	h := NewDecrHandler(s)
	s.Set("foo", "0")
	cmd := []byte("DEC\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")

	for b.Loop() {
		h.Handle(cmd)
	}
}

func BenchmarkDelHandler(b *testing.B) {
	h := NewDelHandler(s)
	s.Set("foo", "bar")
	cmd := []byte("DEL\x00" + strconv.Itoa(len("foo")) + "\x00foo\r\n")

	for b.Loop() {
		h.Handle(cmd)
		s.Set("foo", "bar")
	}
}
