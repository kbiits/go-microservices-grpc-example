package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/kbiits/learn-microservice-grpc-golang/user/userdb"
	"github.com/kbiits/learn-microservice-grpc-golang/user/userpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	timeout = time.Second
)

type userServer struct {
	userpb.UnimplementedUserServiceServer
}

func (*userServer) CreateUpdateUser(ctx context.Context, req *userpb.CreateUpdateUserRequest) (*userpb.CreateUpdateUserResponse, error) {
	if req.Operation == userpb.Operation_CREATE {
		log.Println("Create user")
	} else {
		log.Println("Update user")
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var uid primitive.ObjectID
	if req.Operation == userpb.Operation_CREATE {
		uid = primitive.NewObjectID()
	} else {
		var err error
		uid, err = primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return nil, error_response(err)
		}
	}

	err := userdb.UpsertOne(ctx, &userdb.User{
		Id:       uid,
		Name:     req.Name,
		Age:      req.Age,
		Greeting: req.Greeting,
		Salary:   req.Salary,
		Power:    req.Power,
	})

	if err != nil {
		return nil, error_response(err)
	}

	return &userpb.CreateUpdateUserResponse{}, nil
}

func (*userServer) GetUserDetails(ctx context.Context, req *userpb.GetUserDetailsRequest) (*userpb.GetUserDetailsResponse, error) {
	log.Println("Called GetUserDetails, Id", req.Id)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	uid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, error_response(err)
	}

	result, err := userdb.FindOne(c, uid)
	if err != nil {
		return nil, error_response(err)
	}

	return &userpb.GetUserDetailsResponse{Salary: result.Salary, Power: result.Power}, nil
}

func (*userServer) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	log.Println("Called GetUsers")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := userdb.Find(c)

	if err != nil {
		return nil, error_response(err)
	}

	var res userpb.GetUsersResponse
	for _, d := range *data {
		res.Users = append(res.Users, &userpb.GetUserResponse{Id: d.Id.Hex(), Name: d.Name, Age: d.Age, Greeting: d.Greeting})
	}

	return &res, nil
}

func main() {
	log.Println("User Service")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

	userdb.Mongo_Client, err = userdb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer userdb.Mongo_Client.Disconnect(context.Background())

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, &userServer{})
	reflection.Register(s)

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
}

func error_response(err error) error {
	log.Println("ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}
