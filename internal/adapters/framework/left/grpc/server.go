package rpc

import (
	"context"
	"log"
	"net"

	"github.com/felipewom/go-hexagonal/internal/adapters/framework/left/grpc/pb"
	"github.com/felipewom/go-hexagonal/internal/ports"
	"google.golang.org/grpc"
	// Consider adding codes and status for rich gRPC errors if needed later
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// Adapter implements the GRPCPort interface and pb.ArithmeticServiceServer
type Adapter struct {
	pb.UnimplementedArithmeticServiceServer // Embed for forward compatibility
	api ports.APIPort
}

// NewAdapter creates a new Adapter
func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

// Run registers the ArithmeticServiceServer to a grpcServer and serves on
// the specified port
func (grpca *Adapter) Run() { // Changed to pointer receiver to match gRPC style, though value receiver works
	var err error

	listen, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen on port 9000: %v", err)
	}

	grpcServer := grpc.NewServer()
	// The adapter itself (grpca) implements ArithmeticServiceServer
	pb.RegisterArithmeticServiceServer(grpcServer, grpca)

	log.Printf("gRPC server listening on port 9000")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve gRPC server over port 9000: %v", err)
	}
}

// GetAddition implements the gRPC service method for addition
func (grpca *Adapter) GetAddition(ctx context.Context, req *pb.OperationParameters) (*pb.Answer, error) {
	// The 'ctx' parameter can be used for deadlines, cancellation, etc.
	// For now, it's not directly used in the call to the application port.
	res, err := grpca.api.GetAddition(req.GetA(), req.GetB())
	if err != nil {
		// TODO: Consider translating application errors to gRPC status codes
		return nil, err
	}
	return &pb.Answer{Value: res}, nil
}

// GetSubtraction implements the gRPC service method for subtraction
func (grpca *Adapter) GetSubtraction(ctx context.Context, req *pb.OperationParameters) (*pb.Answer, error) {
	res, err := grpca.api.GetSubtraction(req.GetA(), req.GetB())
	if err != nil {
		return nil, err
	}
	return &pb.Answer{Value: res}, nil
}

// GetMultiplication implements the gRPC service method for multiplication
func (grpca *Adapter) GetMultiplication(ctx context.Context, req *pb.OperationParameters) (*pb.Answer, error) {
	res, err := grpca.api.GetMultiplication(req.GetA(), req.GetB())
	if err != nil {
		return nil, err
	}
	return &pb.Answer{Value: res}, nil
}

// GetDivision implements the gRPC service method for division
func (grpca *Adapter) GetDivision(ctx context.Context, req *pb.OperationParameters) (*pb.Answer, error) {
	res, err := grpca.api.GetDivision(req.GetA(), req.GetB())
	if err != nil {
		// Example: If division by zero returns a specific error type from API,
		// you could map it to codes.InvalidArgument
		// if errors.Is(err, core.ErrDivisionByZero) { // Assuming core.ErrDivisionByZero exists
		// 	return nil, status.Errorf(codes.InvalidArgument, "division by zero")
		// }
		return nil, err
	}
	return &pb.Answer{Value: res}, nil
	}
}
