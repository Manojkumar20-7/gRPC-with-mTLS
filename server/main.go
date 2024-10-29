package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	rpc "grpc/https/proto"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct{
	rpc.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, in *rpc.HelloRequest) (*rpc.HelloReply,error){
	log.Printf("Received: %s\n",in.GetName())
	return &rpc.HelloReply{Message: "Hello "+in.GetName()},nil
}

func main(){
	serverCert,err:=tls.LoadX509KeyPair("certificates/server-certificates/server-cert.pem","certificates/server-certificates/server-key.pem")
	if err!=nil{
		log.Fatalln("Error in loading Server certificate and key",err)
	}
	trustedCert,err:=os.ReadFile("certificates/root-ca-certificates/root-ca-cert.pem")
	if err!=nil{
		log.Fatalln("Error in reading CA certificate")
	}
	certPool:=x509.NewCertPool()
	certPool.AppendCertsFromPEM(trustedCert)
	trustedCert,err=os.ReadFile("certificates/root-ca-certificates/root-cert.pem")
	if err!=nil{
		log.Fatalln("Error in reading CA certificate")
	}
	certPool.AppendCertsFromPEM(trustedCert)
	tlsConfig:=tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs: certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	cred:=credentials.NewTLS(&tlsConfig)

	listen,err:=net.Listen("tcp","localhost:8000")
	if err!=nil{
		log.Fatalln("Error in listening")
	}
	defer listen.Close()
	grpcServer:=grpc.NewServer(grpc.Creds(cred))
	rpc.RegisterGreeterServer(grpcServer,&server{})
	log.Println("Server is listening at https://localhost:8000")
	err=grpcServer.Serve(listen)
	if err!=nil{
		log.Fatalln("Failed to start the server", err)
	}
}