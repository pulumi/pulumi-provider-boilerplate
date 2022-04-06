package provider

import (
	"fmt"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type xyzRandomResource struct {
	config xyzConfig
}

type xyzRandomInput struct {
	Length int
}

func (i *xyzRandomInput) ToPropertyMap() resource.PropertyMap {
	pm := resource.PropertyMap{}
	pm["length"] = resource.NewPropertyValue(i.Length)
	return pm
}

func (r *xyzRandomResource) ToRandomResourceInput(inputMap resource.PropertyMap) xyzRandomInput {
	input := xyzRandomInput{}

	if inputMap["length"].HasValue() && inputMap["length"].IsNumber() {
		input.Length = int(inputMap["length"].NumberValue())
	}

	return input
}

func (r *xyzRandomResource) Configure(config xyzConfig) {
	r.config = config
}

func (r *xyzRandomResource) Name() string {
	return "xyz:index:Random"
}

func (r *xyzRandomResource) Diff(req *pulumirpc.DiffRequest) (*pulumirpc.DiffResponse, error) {
	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	d := olds["__inputs"].ObjectValue().Diff(news)
	changes := pulumirpc.DiffResponse_DIFF_NONE
	if d.Changed("length") {
		changes = pulumirpc.DiffResponse_DIFF_SOME
	}

	return &pulumirpc.DiffResponse{
		Changes:  changes,
		Replaces: []string{"length"},
	}, nil
}

// Create allocates a new instance of the provided resource and returns its unique ID afterwards.
func (r *xyzRandomResource) Create(req *pulumirpc.CreateRequest) (*pulumirpc.CreateResponse, error) {
	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}

	resourceInput := r.ToRandomResourceInput(inputs)

	if resourceInput.Length <= 0 {
		return nil, fmt.Errorf("Expected input property 'length' of type 'number' but got '%s", inputs["length"].TypeString())
	}

	// Actually "create" the random number
	result := makeRandom(resourceInput.Length)

	outputs := map[string]interface{}{
		"length": resourceInput.Length,
		"result": result,
	}

	outputProperties, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputs),
		plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}
	return &pulumirpc.CreateResponse{
		Id:         result,
		Properties: outputProperties,
	}, nil
}

// Check validates that the given property bag is valid for a resource of the given type and returns
// the inputs that should be passed to successive calls to Diff, Create, or Update for this
// resource. As a rule, the provider inputs returned by a call to Check should preserve the original
// representation of the properties as present in the program inputs. Though this rule is not
// required for correctness, violations thereof can negatively impact the end-user experience, as
// the provider inputs are using for detecting and rendering diffs.
func (k *xyzRandomResource) Check(req *pulumirpc.CheckRequest) (*pulumirpc.CheckResponse, error) {
	return &pulumirpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

// Read the current live state associated with a resource.
func (k *xyzRandomResource) Read(req *pulumirpc.ReadRequest) (*pulumirpc.ReadResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Read is not yet implemented for 'xyz:index:Random'")
}

// Update updates an existing resource with new values.
func (k *xyzRandomResource) Update(req *pulumirpc.UpdateRequest) (*pulumirpc.UpdateResponse, error) {
	// Our Random resource will never be updated - if there is a diff, it will be a replacement.
	return nil, status.Error(codes.Unimplemented, "Update is not yet implemented for 'xyz:index:Random'")
}

// Delete tears down an existing resource with the given ID.  If it fails, the resource is assumed
// to still exist.
func (k *xyzRandomResource) Delete(req *pulumirpc.DeleteRequest) (*pbempty.Empty, error) {
	// Note that for our Random resource, we don't have to do anything on Delete.
	return &pbempty.Empty{}, nil
}
