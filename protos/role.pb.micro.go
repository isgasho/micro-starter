// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: role.proto

package protos

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "github.com/micro/micro/v3/service/api"
	client "github.com/micro/micro/v3/service/client"
	server "github.com/micro/micro/v3/service/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Role service

func NewRoleEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Role service

type RoleService interface {
	//  Role
	GetRole(ctx context.Context, in *GetRoleRequest, opts ...client.CallOption) (*GetRoleResponse, error)
	InsertRole(ctx context.Context, in *InsertRoleRequest, opts ...client.CallOption) (*InsertRoleResponse, error)
	DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...client.CallOption) (*DeleteRoleResponse, error)
	UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...client.CallOption) (*UpdateRoleResponse, error)
}

type roleService struct {
	c    client.Client
	name string
}

func NewRoleService(name string, c client.Client) RoleService {
	return &roleService{
		c:    c,
		name: name,
	}
}

func (c *roleService) GetRole(ctx context.Context, in *GetRoleRequest, opts ...client.CallOption) (*GetRoleResponse, error) {
	req := c.c.NewRequest(c.name, "Role.GetRole", in)
	out := new(GetRoleResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleService) InsertRole(ctx context.Context, in *InsertRoleRequest, opts ...client.CallOption) (*InsertRoleResponse, error) {
	req := c.c.NewRequest(c.name, "Role.InsertRole", in)
	out := new(InsertRoleResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleService) DeleteRole(ctx context.Context, in *DeleteRoleRequest, opts ...client.CallOption) (*DeleteRoleResponse, error) {
	req := c.c.NewRequest(c.name, "Role.DeleteRole", in)
	out := new(DeleteRoleResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roleService) UpdateRole(ctx context.Context, in *UpdateRoleRequest, opts ...client.CallOption) (*UpdateRoleResponse, error) {
	req := c.c.NewRequest(c.name, "Role.UpdateRole", in)
	out := new(UpdateRoleResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Role service

type RoleHandler interface {
	//  Role
	GetRole(context.Context, *GetRoleRequest, *GetRoleResponse) error
	InsertRole(context.Context, *InsertRoleRequest, *InsertRoleResponse) error
	DeleteRole(context.Context, *DeleteRoleRequest, *DeleteRoleResponse) error
	UpdateRole(context.Context, *UpdateRoleRequest, *UpdateRoleResponse) error
}

func RegisterRoleHandler(s server.Server, hdlr RoleHandler, opts ...server.HandlerOption) error {
	type role interface {
		GetRole(ctx context.Context, in *GetRoleRequest, out *GetRoleResponse) error
		InsertRole(ctx context.Context, in *InsertRoleRequest, out *InsertRoleResponse) error
		DeleteRole(ctx context.Context, in *DeleteRoleRequest, out *DeleteRoleResponse) error
		UpdateRole(ctx context.Context, in *UpdateRoleRequest, out *UpdateRoleResponse) error
	}
	type Role struct {
		role
	}
	h := &roleHandler{hdlr}
	return s.Handle(s.NewHandler(&Role{h}, opts...))
}

type roleHandler struct {
	RoleHandler
}

func (h *roleHandler) GetRole(ctx context.Context, in *GetRoleRequest, out *GetRoleResponse) error {
	return h.RoleHandler.GetRole(ctx, in, out)
}

func (h *roleHandler) InsertRole(ctx context.Context, in *InsertRoleRequest, out *InsertRoleResponse) error {
	return h.RoleHandler.InsertRole(ctx, in, out)
}

func (h *roleHandler) DeleteRole(ctx context.Context, in *DeleteRoleRequest, out *DeleteRoleResponse) error {
	return h.RoleHandler.DeleteRole(ctx, in, out)
}

func (h *roleHandler) UpdateRole(ctx context.Context, in *UpdateRoleRequest, out *UpdateRoleResponse) error {
	return h.RoleHandler.UpdateRole(ctx, in, out)
}
