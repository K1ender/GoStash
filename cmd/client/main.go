package main

import "net"

func main() {
	conn, err := net.Dial("tcp", "localhost:19201")
	if err != nil {
		panic(err)
	}

	conn.Write([]byte("SET\0003\000key\0005\000value\r\n"))
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	println(string(buf[:n]))

	conn.Write([]byte("GET\0003\000key\r\n"))
	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}
	println(string(buf[:n]))

	conn.Write([]byte("SET\0003\000key\0001\0001\r\n"))
	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}
	println(string(buf[:n]))

	conn.Write([]byte("INC\0003\000key\r\n"))
	buf = make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}
	println(string(buf[:n]))

	conn.Close()
}
