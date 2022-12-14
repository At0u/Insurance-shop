// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// InfoClient is the client API for Info service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InfoClient interface {
	CreateProductSellInfo(ctx context.Context, in *SellInfoReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CreateCategorySellInfo(ctx context.Context, in *SellInfoReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type infoClient struct {
	cc grpc.ClientConnInterface
}

func NewInfoClient(cc grpc.ClientConnInterface) InfoClient {
	return &infoClient{cc}
}

func (c *infoClient) CreateProductSellInfo(ctx context.Context, in *SellInfoReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Info/CreateProductSellInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *infoClient) CreateCategorySellInfo(ctx context.Context, in *SellInfoReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Info/CreateCategorySellInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InfoServer is the server API for Info service.
// All implementations must embed UnimplementedInfoServer
// for forward compatibility
type InfoServer interface {
	CreateProductSellInfo(context.Context, *SellInfoReq) (*emptypb.Empty, error)
	CreateCategorySellInfo(context.Context, *SellInfoReq) (*emptypb.Empty, error)
	mustEmbedUnimplementedInfoServer()
}

// UnimplementedInfoServer must be embedded to have forward compatible implementations.
type UnimplementedInfoServer struct {
}

func (UnimplementedInfoServer) CreateProductSellInfo(context.Context, *SellInfoReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProductSellInfo not implemented")
}
func (UnimplementedInfoServer) CreateCategorySellInfo(context.Context, *SellInfoReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCategorySellInfo not implemented")
}
func (UnimplementedInfoServer) mustEmbedUnimplementedInfoServer() {}

// UnsafeInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InfoServer will
// result in compilation errors.
type UnsafeInfoServer interface {
	mustEmbedUnimplementedInfoServer()
}

func RegisterInfoServer(s grpc.ServiceRegistrar, srv InfoServer) {
	s.RegisterService(&Info_ServiceDesc, srv)
}

func _Info_CreateProductSellInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfoServer).CreateProductSellInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Info/CreateProductSellInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfoServer).CreateProductSellInfo(ctx, req.(*SellInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Info_CreateCategorySellInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfoServer).CreateCategorySellInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Info/CreateCategorySellInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfoServer).CreateCategorySellInfo(ctx, req.(*SellInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Info_ServiceDesc is the grpc.ServiceDesc for Info service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Info_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Info",
	HandlerType: (*InfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProductSellInfo",
			Handler:    _Info_CreateProductSellInfo_Handler,
		},
		{
			MethodName: "CreateCategorySellInfo",
			Handler:    _Info_CreateCategorySellInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "info.proto",
}
