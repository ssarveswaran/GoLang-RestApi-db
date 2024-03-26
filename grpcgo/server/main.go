package main

import (
	// "fmt"
	//
	"io"
	"log"
	pb "main/proto"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	port = ":2000"
)

type server struct {
	pb.GreetServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Println("Failed to access port")
	}
	gserver := grpc.NewServer()
	pb.RegisterGreetServer(gserver, &server{})
	log.Println("Server Started At", lis.Addr())
	if err := gserver.Serve(lis); err != nil {
		log.Println("Error in service", err)
	}

}

// // func (s *server) SayHello(ctx context.Context,request *pb.HelloServer)(*pb.HelloResponse,error){
// // 	return &pb.HelloResponse{
// // 		Message:"Hello Paytabs",
// // 	},nil
// }
// func (s *server) Server(req *pb.HelloServer, stream pb.Greet_ServerServer) error {
// 	name := req.Message

// 	for i := 0; i < len(name); i++ {
// 		if err := stream.Send(&pb.HelloResponse{
// 			Message: name,
// 		}); err != nil {
// 			return err
// 		}
//         time.Sleep(2*time.Second)
// 	}
// 	return nil
// }//	Client(Greet_ClientServer) error
func (s *server) Client(stream pb.Greet_ClientServer)error{
 var message string;
 for {
	req,err:=stream.Recv();
	if err!=nil{
	    log.Println("something went wrong");
	}
	if err==io.EOF{
		stream.SendAndClose(&pb.HelloResponse{Message:message})

	}
	
 }
}