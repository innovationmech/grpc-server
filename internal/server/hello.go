package server

import (
	"context"
	"io"
	"log"

	"github.com/innovationmech/grpc-server/pb"
)

type helloServer struct {
	pb.UnimplementedHelloServiceServer
}

var server *helloServer = &helloServer{}

func (s *helloServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Reply: "Hello " + req.GetGreeting()}, nil
}

func (s *helloServer) LotsOfReplies(req *pb.HelloRequest, stream pb.HelloService_LotsOfRepliesServer) error {
	for i := 0; i < 10; i++ {
		if err := stream.Send(&pb.HelloResponse{Reply: "Hello " + req.GetGreeting()}); err != nil {
			return err
		}
	}
	return nil
}

func (s *helloServer) LotsOfGreetings(stream pb.HelloService_LotsOfGreetingsServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloResponse{Reply: "Done receiving greetings"})
		}
		if err != nil {
			return err
		}
		log.Printf("Received greeting: %v", req.GetGreeting())
	}
}

func (s *helloServer) BidiHello(stream pb.HelloService_BidiHelloServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		err = stream.Send(&pb.HelloResponse{Reply: "Hello " + req.GetGreeting()})
		if err != nil {
			return err
		}
	}
}

func HelloServer() pb.HelloServiceServer {
	return server
}
