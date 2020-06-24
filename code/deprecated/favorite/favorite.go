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

type RecommendServer struct {
	FavorHandler     grpc_transport.Handler
	RecommendHandler grpc_transport.Handler
}

func (rs *RecommendServer) GetFavorite(ctx context.Context, in *bangumi.FavoriteReq) (*bangumi.FavoriteRsp, error) {
	_, rsp, err := rs.FavorHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return rsp.(*bangumi.FavoriteRsp), err
}

func makeGetFavorEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*bangumi.FavoriteReq)
		fmt.Println("favorite req:", req.Id)
		rsp := &bangumi.FavoriteRsp{
			FavoritesList: []*bangumi.BangumiInfo{
				&bangumi.BangumiInfo{
					BangumiId:    1,
					Name:         "B1",
					CoverUrl:     "www.example.com/b1",
					Tags:         "no tags",
					BangumiScore: 1,
				},
				&bangumi.BangumiInfo{
					BangumiId:    2,
					Name:         "B2",
					CoverUrl:     "www.example.com/b2",
					Tags:         "no tags",
					BangumiScore: 2,
				},
			},
			Msg: "OK",
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
	serviceAddress := ":50053"

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	rs := new(RecommendServer)

	var ep endpoint.Endpoint
	ep = loggings.LoggingMiddleware(log.With(logger, "method", "Favor"))(makeGetFavorEndpoint())
	ep = bangumijwt.AuthMiddleware()(ep)
	rs.FavorHandler = grpc_transport.NewServer(
		ep,
		decodeRequest,
		encodeResponse,
	)

	ls, _ := net.Listen("tcp", serviceAddress)
	gs := grpc.NewServer()
	bangumi.RegisterFavoriteServiceServer(gs, rs)
	gs.Serve(ls)
}
