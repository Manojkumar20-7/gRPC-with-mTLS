package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	rpc "grpc/https/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main(){
	clientCert,err:=tls.LoadX509KeyPair("certificates/client-certificates/client1-cert.pem","certificates/client-certificates/client1-key.pem")
	if err!=nil{
		log.Fatalln("Error in loading client certificate",err)
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
	tlsConfig:=&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs: certPool,
		InsecureSkipVerify: false,
	}

	cred:=credentials.NewTLS(tlsConfig)
	conn,err:=grpc.NewClient("localhost:8000",grpc.WithTransportCredentials(cred))
	if err!=nil{
		log.Fatalln("Error in creating connection",err)
	}
	defer conn.Close()
	client:=rpc.NewGreeterClient(conn)

	request:=&rpc.HelloRequest{
		Name:"Manoj",
	}
	response,err:=client.SayHello(context.Background(),request)
	if err!=nil{
		log.Fatalln("Failed to receive response",err)
	}
	log.Println(response)
}