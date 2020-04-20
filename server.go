package main

import (
	"fmt"
	"go-test/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	lis, err := net.Listen("tcp", ":7777")

	if err != nil {
		log.Fatalln("Can't listen port", err)
	}

	server := grpc.NewServer()

	err = godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	service.RegisterNasaEpicServiceServer(server, service.NewDownloadService(os.Getenv("NASA_API_KEY"), 2))

	fmt.Println("Starting server at 7777 port")

	err = server.Serve(lis)

	if err != nil {
		log.Fatalln("Error serving listener", err)
	}
}
