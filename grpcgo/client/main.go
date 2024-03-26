package main

import (
	"context"
	"io"
	"log"
	pb "main/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":2000"
)

func main() {
	li, err := grpc.Dial("localhost"+port,grpc.WithTransportCredentials(insecure.NewCredentials()));
	if err!=nil{
		log.Println("Error in grpc client");
	}
	defer li.Close();
	client:=pb.NewGreetClient(li);
//    result,err:=client.SayHello(context.Background(),&pb.HelloServer{});
//    if err!=nil{
// 	  log.Fatal("Error while calling");
//    }
//    log.Println(result.Message);
//        name := &pb.HelloServer{
// 		Message: "sarveshwaran",
// 	   }
//     result,err:=client.Server(context.Background(),name)
	
// 	for {
// 		message, err := result.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Println("Error while receiving")
// 		} 
// 		log.Println("message",message)
// 	}
// }