package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"
)

var (
	tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 10000, "The server port")
)

type crudServer struct {
	UnimplementedCRUDServer
}

func (s *crudServer) CreateUser(ctx context.Context, user *User) (int, error) {
	return int(user.Id), MysqlDB.Create(user).Error
}

func (s *crudServer) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	var err error
	if err == MysqlDB.Find(&users).Error {
		return users, err
	}
	return users, err
}

func (s *crudServer) UpdateUser(ctx context.Context, id int, newPassword string) (bool, error) {
	var err error
	if err = MysqlDB.Table("users").Where("ID = ?", id).Update("Password", newPassword).Error; err == nil {
		return true, err
	}
	return false, err
}

func (s *crudServer) DeleteUser(ctx context.Context, id int) (bool, error) {
	var err error
	if err = MysqlDB.Where("ID = ?", id).Delete(User{}).Error; err == nil {
		return true, err
	}
	return false, err
}

func newServer() *crudServer {
	// s := &crudServer{routeNotes: make(map[string][]*pb.RouteNote)}
	// s.loadFeatures(*jsonDBFile)
	return s
}

func main() {
	ConnectMysql()
	ConnectRedis()
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	RegisterCRUDServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
