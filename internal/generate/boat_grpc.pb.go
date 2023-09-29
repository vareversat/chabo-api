// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: grpc/v1/boat.proto

package generate

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BoatService_GetBoat_FullMethodName = "/grpc.v1.BoatService/GetBoat"
)

// BoatServiceClient is the client API for BoatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BoatServiceClient interface {
	GetBoat(ctx context.Context, in *BoatNameRequest, opts ...grpc.CallOption) (*Boat, error)
}

type boatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBoatServiceClient(cc grpc.ClientConnInterface) BoatServiceClient {
	return &boatServiceClient{cc}
}

func (c *boatServiceClient) GetBoat(ctx context.Context, in *BoatNameRequest, opts ...grpc.CallOption) (*Boat, error) {
	out := new(Boat)
	err := c.cc.Invoke(ctx, BoatService_GetBoat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BoatServiceServer is the server API for BoatService service.
// All implementations must embed UnimplementedBoatServiceServer
// for forward compatibility
type BoatServiceServer interface {
	GetBoat(context.Context, *BoatNameRequest) (*Boat, error)
	mustEmbedUnimplementedBoatServiceServer()
}

// UnimplementedBoatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBoatServiceServer struct {
}

func (UnimplementedBoatServiceServer) GetBoat(context.Context, *BoatNameRequest) (*Boat, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBoat not implemented")
}
func (UnimplementedBoatServiceServer) mustEmbedUnimplementedBoatServiceServer() {}

// UnsafeBoatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BoatServiceServer will
// result in compilation errors.
type UnsafeBoatServiceServer interface {
	mustEmbedUnimplementedBoatServiceServer()
}

func RegisterBoatServiceServer(s grpc.ServiceRegistrar, srv BoatServiceServer) {
	s.RegisterService(&BoatService_ServiceDesc, srv)
}

func _BoatService_GetBoat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BoatNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoatServiceServer).GetBoat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BoatService_GetBoat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoatServiceServer).GetBoat(ctx, req.(*BoatNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BoatService_ServiceDesc is the grpc.ServiceDesc for BoatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BoatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.v1.BoatService",
	HandlerType: (*BoatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBoat",
			Handler:    _BoatService_GetBoat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/v1/boat.proto",
}
