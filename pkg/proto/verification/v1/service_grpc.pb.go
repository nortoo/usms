// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.31.1
// source: pkg/proto/verification/v1/service.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Service_ListVerificationMethods_FullMethodName = "/nortoo.usms.verification.v1.Service/ListVerificationMethods"
	Service_GetVerificationTarget_FullMethodName   = "/nortoo.usms.verification.v1.Service/GetVerificationTarget"
)

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	// ListVerificationMethods lists the supported verification approaches.
	ListVerificationMethods(ctx context.Context, in *ListVerificationMethodsReq, opts ...grpc.CallOption) (*ListVerificationMethodsResp, error)
	// GetVerificationTarget queries the email address or mobile number of a user through the identifier.
	GetVerificationTarget(ctx context.Context, in *GetVerificationTargetReq, opts ...grpc.CallOption) (*GetVerificationTargetResp, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) ListVerificationMethods(ctx context.Context, in *ListVerificationMethodsReq, opts ...grpc.CallOption) (*ListVerificationMethodsResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListVerificationMethodsResp)
	err := c.cc.Invoke(ctx, Service_ListVerificationMethods_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) GetVerificationTarget(ctx context.Context, in *GetVerificationTargetReq, opts ...grpc.CallOption) (*GetVerificationTargetResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetVerificationTargetResp)
	err := c.cc.Invoke(ctx, Service_GetVerificationTarget_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility.
type ServiceServer interface {
	// ListVerificationMethods lists the supported verification approaches.
	ListVerificationMethods(context.Context, *ListVerificationMethodsReq) (*ListVerificationMethodsResp, error)
	// GetVerificationTarget queries the email address or mobile number of a user through the identifier.
	GetVerificationTarget(context.Context, *GetVerificationTargetReq) (*GetVerificationTargetResp, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedServiceServer struct{}

func (UnimplementedServiceServer) ListVerificationMethods(context.Context, *ListVerificationMethodsReq) (*ListVerificationMethodsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListVerificationMethods not implemented")
}
func (UnimplementedServiceServer) GetVerificationTarget(context.Context, *GetVerificationTargetReq) (*GetVerificationTargetResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVerificationTarget not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}
func (UnimplementedServiceServer) testEmbeddedByValue()                 {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	// If the following call pancis, it indicates UnimplementedServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_ListVerificationMethods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListVerificationMethodsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).ListVerificationMethods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_ListVerificationMethods_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).ListVerificationMethods(ctx, req.(*ListVerificationMethodsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_GetVerificationTarget_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVerificationTargetReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).GetVerificationTarget(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Service_GetVerificationTarget_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).GetVerificationTarget(ctx, req.(*GetVerificationTargetReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nortoo.usms.verification.v1.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListVerificationMethods",
			Handler:    _Service_ListVerificationMethods_Handler,
		},
		{
			MethodName: "GetVerificationTarget",
			Handler:    _Service_GetVerificationTarget_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/proto/verification/v1/service.proto",
}
