package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/darshankapadiya19/grpc-gw-rest-proto/gen"
)

// server is used to implement hello.GreeterServer
type helloServer struct {
	pb.UnimplementedHelloServiceServer
}

// SayHello implements hello.GreeterServer
func (s *helloServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

// RegisterGatewayServer sets up the HTTP server and the gRPC-Gateway to handle HTTP requests
func RegisterGatewayServer(ctx context.Context, grpcEndpoint string) (*runtime.ServeMux, error) {
	//	mux := runtime.NewServeMux(
	//		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
	//			MarshalOptions: protojson.MarshalOptions{
	//				UseProtoNames:   true,
	//				EmitUnpopulated: true,
	//			},
	//		}),
	//	)
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption("application/protobuf", &runtime.ProtoMarshaller{}),
	)
	// Register the Greeter service with the HTTP server
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func main() {
	grpcEndpoint := ":50051"
	httpEndpoint := ":8080"

	// Listen on a TCP port
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the Greeter service with the gRPC server
	pb.RegisterHelloServiceServer(server, &helloServer{})

	// Enable server reflection (useful for debugging)
	reflection.Register(server)

	// Start the server
	log.Printf("Server is listening on port 50051...")
	go server.Serve(lis)

	// Create a context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Listen for HTTP/REST connections
	gatewayMux, err := RegisterGatewayServer(ctx, grpcEndpoint)
	if err != nil {
		log.Fatalf("Failed to register HTTP gateway: %v", err)
	}

	// Set up a Gorilla mux router to handle additional routes
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(gatewayMux) // Forward all requests to the gatewayMux

	log.Printf("HTTP server listening on %s", httpEndpoint)
	if err := http.ListenAndServe(httpEndpoint, r); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
