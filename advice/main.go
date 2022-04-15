package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kbiits/learn-microservice-grpc-golang/advice/advicedb"
	"github.com/kbiits/learn-microservice-grpc-golang/advice/advicepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	timeout = time.Second
	pg_db   *pgxpool.Pool
)

type adviceServer struct {
	advicepb.UnimplementedAdviceServiceServer
}

func (*adviceServer) CreateUpdateAdvice(ctx context.Context, req *advicepb.CreateUpdateAdviceRequest) (*advicepb.CreateUpdateAdviceResponse, error) {
	log.Println("create update advice, operation:", req.Operation)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var err error
	if req.Operation == advicepb.Operation_CREATE {
		err = advicedb.CreateOne(pg_db, ctx, &advicedb.Advice{
			UserId: req.UserId,
			Advice: req.Advice,
		})
	} else {
		err = advicedb.UpdateOne(pg_db, ctx, &advicedb.Advice{
			UserId: req.UserId,
			Advice: req.Advice,
		})
	}

	if err != nil {
		return nil, error_response(err)
	}

	return &advicepb.CreateUpdateAdviceResponse{}, nil
}

func (*adviceServer) GetAdvice(ctx context.Context, req *advicepb.GetUserAdviceRequest) (*advicepb.GetUserAdviceResponse, error) {
	log.Println("Called GetAdvice for User Id", req.UserId)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result, err := advicedb.FindOne(pg_db, c, req.UserId)
	if err != nil {
		return nil, error_response(err)
	}

	return &advicepb.GetUserAdviceResponse{Advice: result.Advice, CreatedAt: timestamppb.New(result.CreatedAt)}, nil
}

func error_response(err error) error {
	log.Println("ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func main() {
	log.Println("ADVICE SERVICE")

	lis, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("something went wrong %v\n", err)
	}

	pg_db, err = advicedb.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to db %v\n", err)
	}
	defer pg_db.Close()

	s := grpc.NewServer()
	advicepb.RegisterAdviceServiceServer(s, &adviceServer{})
	reflection.Register(s)

	log.Printf("Service started at %v\n", lis.Addr().String())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}
