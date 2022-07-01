package server

import (
	"fmt"

	"google.golang.org/grpc"
)

type Server interface {
	ListPosts()
}

type server struct {
	conn *grpc.ClientConn
}

func NewServer(conn *grpc.ClientConn) Server {
	return &server{
		conn: conn,
	}
}

func (srv server) ListPosts() {
	fmt.Println("1. 通过接口 Service 暴露对外的 ListPosts 方法；")
	fmt.Println("2. 使用 NewService 函数初始化 Service 接口的实现并通过私有的接口体 service 持有 grpc 连接；")
	fmt.Println("3. ListPosts 不再依赖全局变量，而是依赖接口体 service 持有的连接；")
}
