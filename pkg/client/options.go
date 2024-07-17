package client

import "net"

const(
  TCP = "tcp"
  UDP = "udp"
)

type ClientOptionsFunc func(*ClientOptions)

type ClientOptions struct{
  ConnectionAddr  string
  Protocol        string
}

type Client struct{
  ClientOptions
  Conn            net.Conn
}

func defaultClientOptions() ClientOptions{
  return ClientOptions{
    ConnectionAddr: "localhost:3306",
    Protocol: TCP,
  }
}

func WithConnectionAddr(connectionAddr string) ClientOptionsFunc{
  return func(opts *ClientOptions){
    opts.ConnectionAddr = connectionAddr
  }
}

func WithProtocol(protocol string) ClientOptionsFunc{
  return func(opts *ClientOptions){
    opts.Protocol = protocol
  }
}
