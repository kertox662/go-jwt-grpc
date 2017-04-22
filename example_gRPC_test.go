package bearerware_test

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ckaznocha/go-JWTBearerware"
	"github.com/dgrijalva/jwt-go"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	certFile  = "./test_cert/server.pem"
	keyFile   = "./test_cert/server.key"
	host      = "127.0.0.1"
	port      = "50051"
	netString = "tcp"
)

var (
	jwtKey        = []byte("MySecret")
	signingMethod = jwt.SigningMethodHS256
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}

// SayHello implements helloworld.GreeterServer it requires a valid JWT
func (s *server) SayHello(
	ctx context.Context,
	in *pb.HelloRequest,
) (*pb.HelloReply, error) {
	//Validate and extract the JWT from the context using
	//bearerware.JWTFromContext
	token, err := bearerware.JWTFromContext(ctx, jwtKeyFunc, signingMethod)
	if err != nil {
		return nil, err
	}
	return &pb.HelloReply{
		Message: fmt.Sprintf(
			"Hello %s! Token signed using %s",
			in.Name,
			token.Method.Alg(),
		),
	}, nil
}

func Example_gRPC() {
	//The server needs to be started using TLS
	var (
		cert, _ = tls.LoadX509KeyPair(certFile, keyFile)
		opts    = []grpc.ServerOption{
			grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		}
	)
	lis, err := net.Listen(netString, net.JoinHostPort(host, port))
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	//Start the server
	s := grpc.NewServer(opts...)
	pb.RegisterGreeterServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Print(err)
		}
	}()
	defer s.Stop()

	// Set up a connection to the server using TLS and a JWT
	var (
		tlsCreds, _ = credentials.NewClientTLSFromFile(certFile, "localhost")
		//Create a JWT for  the example
		tokenString, _ = jwt.New(signingMethod).SignedString(jwtKey)
		jwtCreds, _    = bearerware.NewJWTAccessFromJWT(tokenString)
		dialOpts       = []grpc.DialOption{
			grpc.WithTransportCredentials(tlsCreds),
			//Pass our jwtCreds to grpc.WithPerRPCCredentials to have it
			//included in every request.
			grpc.WithPerRPCCredentials(jwtCreds),
			grpc.WithTimeout(5 * time.Second),
			grpc.WithBlock(),
		}
	)
	conn, err := grpc.Dial(net.JoinHostPort(host, port), dialOpts...)
	if err != nil {
		panic(fmt.Sprintf("did not connect: %v", err))
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Print(err)
		}
	}()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	// Our JWT is included in every request; no extra steps needed.
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "World"})
	if err != nil {
		panic(fmt.Sprintf("could not greet: %v", err))
	}
	fmt.Printf("Greeting: %s", r.Message)
	// Output: Greeting: Hello World! Token signed using HS256
}
