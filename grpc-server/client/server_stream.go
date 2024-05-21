package main

import (
	"context"
	pb "grpcserver/proto"
	"io"
	"log"
)

func callSayHelloServerStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Streaming started")
	stream, err := client.SayHelloServerStreaming(context.Background(), names) // send req stream to server
	if err != nil {
		log.Fatalf("could not send names: %v", err)
	}

	for {
		message, err := stream.Recv() // received the stream response , serially
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming %v", err)
		}
		log.Println(message)
	}
	log.Printf("Streaming finished")
}
