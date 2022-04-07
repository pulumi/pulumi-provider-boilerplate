package provider

import (
	"fmt"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

type xyzUnknownResource struct{}
type xyzUnknownFunction struct{}

func (c xyzUnknownResource) Name() string {
	return "xyz:index:Unknown"
}

func (u *xyzUnknownResource) Configure(config xyzConfig) {
}

func (c *xyzUnknownResource) Diff(req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	return nil, createUnknownResourceErrorFromRequest(req)
}

func (c *xyzUnknownResource) Delete(req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	return nil, createUnknownResourceErrorFromRequest(req)
}

func (c *xyzUnknownResource) Create(req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	return nil, createUnknownResourceErrorFromRequest(req)
}

func (k *xyzUnknownResource) Check(req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return nil, createUnknownResourceErrorFromRequest(req)
}

func (k *xyzUnknownResource) Update(req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	return nil, createUnknownResourceErrorFromRequest(req)
}

func (k *xyzUnknownResource) Read(req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	return nil, createUnknownResourceErrorFromRequest(req)
}

func createUnknownResourceErrorFromRequest(req ResourceBase) error {
	rn := getResourceNameFromRequest(req)
	return fmt.Errorf("unknown resource type '%s'", rn)
}

func (f *xyzUnknownResource) Invoke(s *xyzProvider, req *pulumirpc.InvokeRequest) (*pulumirpc.InvokeResponse, error) {
	return &rpc.InvokeResponse{Return: nil}, fmt.Errorf("unknown function '%s'", req.Tok)
}
