package main

import (
	"context"
	"fmt"
	"go-test/service"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

var downloadService service.NasaEpicServiceClient

func main() {
	grpcConn, err := grpc.Dial(
		os.Getenv("NASA_API_KEY") + ":7777",
		grpc.WithInsecure(),
		)

	if err != nil {
		log.Fatalf("Can't connect to rpc.")
	}
	defer func() {
		err = grpcConn.Close()
		if err != nil {
			log.Fatalln("Grpc connection close error", err)
		}
	}()

	downloadService = service.NewNasaEpicServiceClient(grpcConn)

	http.HandleFunc("/", mainHandler)
	fs := http.FileServer(http.Dir("./downloads"))
	http.Handle("/downloads/", http.StripPrefix("/downloads/", fs))

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		http.ServeFile(w, r, "./html/index.html")
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		date := r.FormValue("date")

		reqStr := &service.RequestString{
			Date: date,
		}

		ctx := context.Background()
		json, err := downloadService.DownloadNatural(ctx, reqStr)

		if err != nil {
			log.Fatalln(err)
		}

		fmt.Fprintf(w, "EPIC Download microservice response: %v", json)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
