package main

import (
	bangumi "bangumi/pb"
	bangumijwt "bangumi/utils/bangumi-jwt"
	loggings "bangumi/utils/log"
	"fmt"
	"strconv"

	"context"
	"net"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type InfoServer struct {
	bangumiListHandler  grpc_transport.Handler
	bangumiDetailHander grpc_transport.Handler
}

func (ls *InfoServer) GetBangumiDetail(ctx context.Context, in *bangumi.BangumiDetailReq) (*bangumi.BangumiDetail, error) {
	_, rsp, err := ls.bangumiDetailHander.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return rsp.(*bangumi.BangumiDetail), err
}

func (ls *InfoServer) GetBangumiList(ctx context.Context, in *bangumi.BangumiListReq) (*bangumi.BangumiList, error) {
	_, rsp, err := ls.bangumiListHandler.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return rsp.(*bangumi.BangumiList), err
}

func makeGetBangumiListEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*bangumi.BangumiListReq)
		fmt.Println("Bangumi list req:", req)

		rsp := new(bangumi.BangumiList)
		rsp.BangumiList = []*bangumi.BangumiInfo{
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
		}
		return rsp, nil
	}
}

func makeGetBangumiDetailEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*bangumi.BangumiDetailReq)
		fmt.Println("bangumi info req:", req)

		rsp := &bangumi.BangumiDetail{
			BangumiId:    req.BangumiId,
			Name:         "B" + strconv.FormatInt(int64(req.BangumiId), 10),
			CoverUrl:     "www.example.com",
			BangumiScore: 1,
			VoteNum:      10,
			EpisodeNum:   10,
			Tags:         "no tags",
			Desc:         "no descriptions",
			StaffList:    "no staff",
			CvList:       "no cv",
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
	serviceAddress := ":50052"

	is := new(InfoServer)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var ep endpoint.Endpoint
	ep = loggings.LoggingMiddleware(log.With(logger, "method", "All"))(makeGetBangumiListEndpoint())
	ep = bangumijwt.AuthMiddleware()(ep)
	is.bangumiListHandler = grpc_transport.NewServer(
		ep,
		decodeRequest,
		encodeResponse,
	)

	ep = loggings.LoggingMiddleware(log.With(logger, "method", "Detail"))(makeGetBangumiDetailEndpoint())
	ep = bangumijwt.AuthMiddleware()(ep)
	is.bangumiDetailHander = grpc_transport.NewServer(
		ep,
		decodeRequest,
		encodeResponse,
	)

	ls, _ := net.Listen("tcp", serviceAddress)
	gs := grpc.NewServer()
	bangumi.RegisterInfoServiceServer(gs, is)
	gs.Serve(ls)
}
