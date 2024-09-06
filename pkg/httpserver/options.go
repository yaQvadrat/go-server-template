package httpserver

import (
	"net"
	"time"
)

type Option func(*HttpServer)

func ReadTimeout(timeout time.Duration) Option {
	return func(hs *HttpServer) {
		hs.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(hs *HttpServer) {
		hs.server.WriteTimeout = timeout
	}
}

func Port(port string) Option {
	return func(hs *HttpServer) {
		hs.server.Addr = net.JoinHostPort("", port)
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(hs *HttpServer) {
		hs.shutdownTimeout = timeout
	}
}
