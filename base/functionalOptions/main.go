package main

import (
	"crypto/tls"
	"time"
)

// 1. 配置选项问题
// 针对是否使用默认值的不同组合，需要有多种不同的创建不同配置的函数
// 因为 Go 不支持重载函数，所有需要使用不同的函数名来应对不同的配置选项
type Server struct {
	Addr           string
	Port           int
	Protocol       string
	Timeout        time.Duration
	MaxConnections int
	TLS            *tls.Config
}

func NewDefaultServer(addr string, port int) (*Server, error) {
	return &Server{addr, port, "tcp", 30 * time.Second, 100, nil}, nil
}
func NewTLSServer(addr string, port int, tls *tls.Config) (*Server, error) {
	return &Server{addr, port, "tcp", 30 * time.Second, 100, tls}, nil
}
func NewServerWithTimeout(addr string, port int, timeout time.Duration) (*Server, error) {
	return &Server{addr, port, "tcp", timeout, 100, nil}, nil
}
func NewTLSServerWithMaxConnAndTimeout(addr string, port int, maxconns int, timeout time.Duration, tls *tls.Config) (*Server, error) {
	return &Server{addr, port, "tcp", 30 * time.Second, maxconns, tls}, nil
}

// 2. 配置对象方案
// 将非必填项移到一个结构体里
type Config struct {
	Protocol       string
	Timeout        time.Duration
	MaxConnections int
	TLS            *tls.Config
}

type Server2 struct {
	Addr string
	Port int
	Conf *Config
}

func NewServer(addr string, port int, conf *Config) (*Server2, error) {
	return &Server2{addr, port, conf}, nil
}

// default config
//srv1, _ := NewServer("localhost", 9000, nil)

// customized config
//conf := ServerConfig{Protocol:"tcp", Timeout: 60*time.Duration}
//srv2, _ := NewServer("locahost", 9000, &conf)

// 3. Builder 模式
type ServerBuilder struct {
	Server
}

func (sb *ServerBuilder) Create(addr string, port int) *ServerBuilder {
	sb.Server.Addr = addr
	sb.Server.Port = port
	//其它代码设置其它成员的默认值
	return sb
}
func (sb *ServerBuilder) WithProtocol(protocol string) *ServerBuilder {
	sb.Server.Protocol = protocol
	return sb
}
func (sb *ServerBuilder) WithMaxConn(maxconn int) *ServerBuilder {
	sb.Server.MaxConnections = maxconn
	return sb
}
func (sb *ServerBuilder) WithTimeOut(timeout time.Duration) *ServerBuilder {
	sb.Server.Timeout = timeout
	return sb
}
func (sb *ServerBuilder) WithTLS(tls *tls.Config) *ServerBuilder {
	sb.Server.TLS = tls
	return sb
}
func (sb *ServerBuilder) Build() Server {
	return sb.Server
}

// 通过 Builder 构建
//sb := ServerBuilder{}
//server := sb.Create("127.0.0.1", 8080).
//WithProtocol("udp").
//WithMaxConn(1024).
//WithTimeOut(30*time.Second).
//Build()

// 4. Functional Options
type Option func(s *Server2)

func Protocol(p string) Option {
	return func(s *Server2) {
		s.Conf.Protocol = p
	}
}
func Timeout(d time.Duration) Option {
	return func(s *Server2) {
		s.Conf.Timeout = d
	}
}
func MaxConnections(n int) Option {
	return func(s *Server2) {
		s.Conf.MaxConnections = n
	}
}
func TLS(tls *tls.Config) Option {
	return func(s *Server2) {
		s.Conf.TLS = tls
	}
}

func NewServer4(addr string, port int, options ...func(*Server2)) (*Server2, error) {
	server := Server2{
		Addr: addr,
		Port: port,
		Conf: &Config{
			Protocol:       "tcp",
			Timeout:        30 * time.Second,
			MaxConnections: 1000,
			TLS:            nil,
		},
	}
	for _, option := range options {
		option(&server)
	}
	return &server, nil
}

//s1, _ := NewServer("localhost", 1024)
//s2, _ := NewServer("localhost", 2048, Protocol("udp"))
//s3, _ := NewServer("0.0.0.0", 8080, Timeout(300*time.Second), MaxConns(1000))
