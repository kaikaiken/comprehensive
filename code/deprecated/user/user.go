package main

import (
	bangumi "bangumi/pb"
	bangumijwt "bangumi/utils/bangumi-jwt"
	loggings "bangumi/utils/log"
	"context"
	"fmt"
	"net"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type UserServer struct {
	SessionsHandler grpc_transport.Handler
	UserHandler     grpc_transport.Handler
}

func (us *UserServer) Sessions(ctx context.Context, in *bangumi.SessionReq) (*bangumi.SessionRsp, error) {
	_, rsp, err := us.SessionsHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return rsp.(*bangumi.SessionRsp), err
}

func (us *UserServer) Register(ctx context.Context, in *bangumi.RegisterReq) (*bangumi.RegisterRsp, error) {
	_, rsp, err := us.UserHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return rsp.(*bangumi.RegisterRsp), err
}

func makeSessionEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*bangumi.SessionReq)
		fmt.Println("session req:", req)

		rsp := &bangumi.SessionRsp{
			Id:             12345,
			Username:       "hhh",
			Email:          "12345@qq.com",
			FavBangumiList: "1,2,3,4,5",
			Msg:            "login succeed",
		}

		if req.Id == 12345 && req.Password == "hello12345" {
			token, err := bangumijwt.Sign(rsp.Username, req.Id)
			if err != nil {
				rsp.Msg = "Sign error"
				return rsp, nil
			}
			rsp.Jwt = token
		}
		return rsp, nil
	}
}

func makeRegisterEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*bangumi.RegisterReq)
		fmt.Println("register req: ", req)
		rsp := &bangumi.RegisterRsp{
			Msg: "register failed: not open yet",
		}
		return rsp, nil
	}
}

func decodeRequest(c context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodeResponse(_ context.Context, rsp interface{}) (interface{}, error) {
	return rsp, nil
}

func main() {
	serviceAddress := ":50054"

	us := new(UserServer)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	var ep endpoint.Endpoint
	ep = loggings.LoggingMiddleware(log.With(logger, "method", "Session"))(makeSessionEndpoint())
	us.SessionsHandler = grpc_transport.NewServer(
		ep,
		decodeRequest,
		encodeResponse,
	)

	ep = loggings.LoggingMiddleware(log.With(logger, "method", "Register"))(makeRegisterEndpoint())
	ep = bangumijwt.AuthMiddleware()(ep)
	us.UserHandler = grpc_transport.NewServer(
		ep,
		decodeRequest,
		encodeResponse,
	)

	ls, _ := net.Listen("tcp", serviceAddress)
	gs := grpc.NewServer()
	bangumi.RegisterUserServiceServer(gs, us)
	gs.Serve(ls)
}
