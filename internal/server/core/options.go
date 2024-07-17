package core

const(
  TCP = "tcp"
  UDP = "udp"
)

type ServerOptionFunc func (*ServerOptions)

type ServerOptions struct{
  Id          string
  Protocol    string 
  ListenAddr  string
  MaxConn     int8
  Tls         bool
}

func defaultOptions() ServerOptions {
  return ServerOptions{
    Id: "default",
    Protocol: TCP,
    ListenAddr: "localhost:3306",
    MaxConn: 100,
    Tls: false,
  }
}

func WithTls(opts *ServerOptions){
  opts.Tls = true
}

func WithMaxConn(n int8) ServerOptionFunc {
  return func(opts *ServerOptions){
    opts.MaxConn = n
  }
}

func WithProtocol(protocol string) ServerOptionFunc {
  return func(opts *ServerOptions){
    opts.Protocol = protocol
  }
}


