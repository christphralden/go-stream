package client

import (
	"log"
	"net"
	"time"

	"github.com/christpheralden/go-stream/types"
)

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


func NewClient(opts ...ClientOptionsFunc) *Client {
  o := defaultClientOptions()

  for _, fn := range opts{
    fn(&o)
  }

  return &Client{
    ClientOptions: o,
  }
}

func (c *Client) Dial() error {
  timeout := time.Second * 5
  conn, err := net.DialTimeout(c.Protocol, c.ConnectionAddr, timeout)  
  c.Conn = conn
  

  if err != nil{
    return err
  }

  return nil
}


func (c *Client) Stop() error {
  if c.Conn != nil{
    c.Conn.Close()
  }

  return nil
}

func (c *Client) SendMessage(msg string) error {
  payload := types.Binary(msg)

  _, err := payload.WriteTo(c.Conn)

  if err != nil {
    log.Println("Error sending message: ", err)
    return err
  }

  reply, err := types.Decode(c.Conn)


  if err != nil {
    log.Println("Error decoding reply: ", err)
    return err
  }

  log.Println("Client recieved: ", string(reply.Bytes()))

  return nil
}
