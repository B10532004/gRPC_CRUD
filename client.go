package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
	pb "google.golang.org/grpc/examples/route_guide/routeguide"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

func printCreateResult(client CRUDClient, user *CreateUserRequest) {
	log.Println("Getting userid after create user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := client.CreateUser(ctx, user)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)
}

func printUsers(client CRUDClient, in *GetUsersRequest) {
	log.Println("Getting users from userlist")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userlist, err := client.GetUsers(ctx, in)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(userlist)
}

func printUpdateResult() {

}

func getDeleteResult() {

}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = data.Path("x509/ca_cert.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := NewCRUDClient(conn)

	printCreateResult(clinet, &CreateUserRequest{User{Id: 1, Username: "111", Password: "111", Phone: "111"}})

	printUsers(client, &GetUsersRequest{})
	// Looking for a valid feature
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
}
