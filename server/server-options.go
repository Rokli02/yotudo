package server

import (
	"fmt"
)

type ServerOptions struct {
	addr         string
	port         int
	listenerIp   string
	hostFrontend bool
}

type ServerOptionsFunc func(so *ServerOptions)

func (so *ServerOptions) Host(host bool) *ServerOptions {
	if host {
		so.listenerIp = "0.0.0.0"
	} else {
		so.listenerIp = "localhost"
	}

	return so
}

func (so *ServerOptions) String() string {
	return fmt.Sprintf(
		"ServerOptions(addr=%s, port=%d, listenerIp=%s, hostFrontend=%t)",
		so.addr, so.port, so.listenerIp, so.hostFrontend,
	)
}

func (so *ServerOptions) HostFrontend(host bool) *ServerOptions {
	so.hostFrontend = host

	return so
}

func (so *ServerOptions) SetPort(port int) *ServerOptions {
	if port < 0 || port > 65535 {
		so.port = 0
	}

	so.port = port

	return so
}
