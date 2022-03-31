package grpcServer

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	courierProto "stlab.itechart-group.com/go/food_delivery/courier_service/GRPC"
	"stlab.itechart-group.com/go/food_delivery/courier_service/service"
)

type GRPCServer struct {
	service *service.Service
	courierProto.UnimplementedCourierServerServer
}

func NewGRPCServer(service *service.Service) {
	s := grpc.NewServer()
	str := &GRPCServer{service: service}
	courierProto.RegisterCourierServerServer(s, str)
	lis, err := net.Listen("tcp", ":8091")
	if err != nil {
		log.Fatalf("NewGRPCServer, Listen:%s", err)
	}
	reflection.Register(s)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("NewGRPCServer, Serve:%s", err)
	}

}

func (g *GRPCServer) CreateOrder(ctx context.Context, order *courierProto.OrderCourierServer) (*emptypb.Empty, error) {
	return g.service.CreateOrder(order)
}

func (g *GRPCServer) GetDeliveryServicesList(ctx context.Context, in *emptypb.Empty) (*courierProto.ServicesResponse, error) {
	res, err := g.service.GetServices(in)
	if err != nil {
		log.Fatalf("GetServices:%s", err)
		return nil, fmt.Errorf("GetServices:%w", err)
	}
	return res, nil
}
