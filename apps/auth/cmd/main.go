package main

import (
	"net"

	au "wave/apps/auth"
	iu "wave/apps/auth/iauth"
	sv "wave/internal/service"

	"google.golang.org/grpc"
)

func main() {
	var (
		c = sv.LoadConfig("configs/ms_auth.json")
		s = au.NewAuthService(c)
		p = s.GetPort()
	)
	con, err := net.Listen("tcp", p)
	if err != nil {
		s.Logger().Errorf("tcp port %s cannot be rised: %s", p, err.Error())
		sv.Panic(err)
	}

	server := grpc.NewServer()
	iu.RegisterIAuthServer(server, s)
	s.Logger().Infof("Auth server has been started on a port %s", p)
	server.Serve(con)
}
