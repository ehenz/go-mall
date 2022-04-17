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

// StockClient is the client API for Stock service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StockClient interface {
	SetStock(ctx context.Context, in *StockInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CheckStock(ctx context.Context, in *StockInfo, opts ...grpc.CallOption) (*StockInfo, error)
	PreSell(ctx context.Context, in *SellInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CancelOrder(ctx context.Context, in *SellInfo, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type stockClient struct {
	cc grpc.ClientConnInterface
}

func NewStockClient(cc grpc.ClientConnInterface) StockClient {
	return &stockClient{cc}
}

func (c *stockClient) SetStock(ctx context.Context, in *StockInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Stock/SetStock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) CheckStock(ctx context.Context, in *StockInfo, opts ...grpc.CallOption) (*StockInfo, error) {
	out := new(StockInfo)
	err := c.cc.Invoke(ctx, "/Stock/CheckStock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) PreSell(ctx context.Context, in *SellInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Stock/PreSell", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) CancelOrder(ctx context.Context, in *SellInfo, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Stock/CancelOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StockServer is the server API for Stock service.
// All implementations must embed UnimplementedStockServer
// for forward compatibility
type StockServer interface {
	SetStock(context.Context, *StockInfo) (*emptypb.Empty, error)
	CheckStock(context.Context, *StockInfo) (*StockInfo, error)
	PreSell(context.Context, *SellInfo) (*emptypb.Empty, error)
	CancelOrder(context.Context, *SellInfo) (*emptypb.Empty, error)
	mustEmbedUnimplementedStockServer()
}

// UnimplementedStockServer must be embedded to have forward compatible implementations.
type UnimplementedStockServer struct {
}

func (UnimplementedStockServer) SetStock(context.Context, *StockInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStock not implemented")
}
func (UnimplementedStockServer) CheckStock(context.Context, *StockInfo) (*StockInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckStock not implemented")
}
func (UnimplementedStockServer) PreSell(context.Context, *SellInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreSell not implemented")
}
func (UnimplementedStockServer) CancelOrder(context.Context, *SellInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
func (UnimplementedStockServer) mustEmbedUnimplementedStockServer() {}

// UnsafeStockServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StockServer will
// result in compilation errors.
type UnsafeStockServer interface {
	mustEmbedUnimplementedStockServer()
}

func RegisterStockServer(s grpc.ServiceRegistrar, srv StockServer) {
	s.RegisterService(&Stock_ServiceDesc, srv)
}

func _Stock_SetStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).SetStock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Stock/SetStock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).SetStock(ctx, req.(*StockInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_CheckStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StockInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).CheckStock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Stock/CheckStock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).CheckStock(ctx, req.(*StockInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_PreSell_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).PreSell(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Stock/PreSell",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).PreSell(ctx, req.(*SellInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_CancelOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).CancelOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Stock/CancelOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).CancelOrder(ctx, req.(*SellInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Stock_ServiceDesc is the grpc.ServiceDesc for Stock service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stock_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Stock",
	HandlerType: (*StockServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetStock",
			Handler:    _Stock_SetStock_Handler,
		},
		{
			MethodName: "CheckStock",
			Handler:    _Stock_CheckStock_Handler,
		},
		{
			MethodName: "PreSell",
			Handler:    _Stock_PreSell_Handler,
		},
		{
			MethodName: "CancelOrder",
			Handler:    _Stock_CancelOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}
