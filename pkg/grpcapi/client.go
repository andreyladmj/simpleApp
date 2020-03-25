package grpcapi

import (
	"andreyladmj/analytics/pkg/grpcapi/userpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)


type User struct {
	Name           string
	Email          string
	Picture        string
	Gender         string
	Locale         string
	Created        time.Time
}

type GRPCClient struct {
	client userpb.AuthServiceClient
}

func New(cc *grpc.ClientConn) *GRPCClient {

	c := userpb.NewAuthServiceClient(cc)

	return &GRPCClient{client:c}
}

func (cc *GRPCClient) GetUser(token string) *User {
	req := &userpb.AuthRequest{Token:token}

	res, err := cc.client.GetUser(context.Background(), req)
	if err != nil {
		log.Printf("error while calling GetUser RPC: %v", err)
		return nil
	}

	if res.User == nil {
		log.Printf("error while calling GetUser RPC: %v, %v", res.Error, res.Status)
		return nil
	}

	fmt.Println("res", res)
	created, err := time.Parse("2006-01-02T15:04:05", res.User.Created)

	if err != nil {
		log.Printf("error while parsing time: %v", err)
		return nil
	}

	user := &User{
		Name:    res.User.Name,
		Email:   res.User.Email,
		Picture: res.User.Picture,
		Gender:  res.User.Gender,
		Locale:  res.User.Locale,
		Created: created,
	}

	return user
}