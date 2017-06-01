// *** WARNING: this file was generated by the Lumi IDL Compiler (LUMIDL). ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package ec2

import (
    "errors"

    pbempty "github.com/golang/protobuf/ptypes/empty"
    pbstruct "github.com/golang/protobuf/ptypes/struct"
    "golang.org/x/net/context"

    "github.com/pulumi/lumi/pkg/resource"
    "github.com/pulumi/lumi/pkg/tokens"
    "github.com/pulumi/lumi/pkg/util/contract"
    "github.com/pulumi/lumi/pkg/util/mapper"
    "github.com/pulumi/lumi/sdk/go/pkg/lumirpc"
)

/* RPC stubs for Instance resource provider */

// InstanceToken is the type token corresponding to the Instance package type.
const InstanceToken = tokens.Type("aws:ec2/instance:Instance")

// InstanceProviderOps is a pluggable interface for Instance-related management functionality.
type InstanceProviderOps interface {
    Check(ctx context.Context, obj *Instance) ([]mapper.FieldError, error)
    Create(ctx context.Context, obj *Instance) (resource.ID, error)
    Get(ctx context.Context, id resource.ID) (*Instance, error)
    InspectChange(ctx context.Context,
        id resource.ID, old *Instance, new *Instance, diff *resource.ObjectDiff) ([]string, error)
    Update(ctx context.Context,
        id resource.ID, old *Instance, new *Instance, diff *resource.ObjectDiff) error
    Delete(ctx context.Context, id resource.ID) error
}

// InstanceProvider is a dynamic gRPC-based plugin for managing Instance resources.
type InstanceProvider struct {
    ops InstanceProviderOps
}

// NewInstanceProvider allocates a resource provider that delegates to a ops instance.
func NewInstanceProvider(ops InstanceProviderOps) lumirpc.ResourceProviderServer {
    contract.Assert(ops != nil)
    return &InstanceProvider{ops: ops}
}

func (p *InstanceProvider) Check(
    ctx context.Context, req *lumirpc.CheckRequest) (*lumirpc.CheckResponse, error) {
    contract.Assert(req.GetType() == string(InstanceToken))
    obj, _, decerr := p.Unmarshal(req.GetProperties())
    if decerr == nil || len(decerr.Failures()) == 0 {
        failures, err := p.ops.Check(ctx, obj)
        if err != nil {
            return nil, err
        }
        if len(failures) > 0 {
            decerr = mapper.NewDecodeErr(failures)
        }
    }
    return resource.NewCheckResponse(decerr), nil
}

func (p *InstanceProvider) Name(
    ctx context.Context, req *lumirpc.NameRequest) (*lumirpc.NameResponse, error) {
    contract.Assert(req.GetType() == string(InstanceToken))
    obj, _, decerr := p.Unmarshal(req.GetProperties())
    if decerr != nil {
        return nil, decerr
    }
    if obj.Name == nil || *obj.Name == "" {
        if req.Unknowns[Instance_Name] {
            return nil, errors.New("Name property cannot be computed from unknown outputs")
        }
        return nil, errors.New("Name property cannot be empty")
    }
    return &lumirpc.NameResponse{Name: *obj.Name}, nil
}

func (p *InstanceProvider) Create(
    ctx context.Context, req *lumirpc.CreateRequest) (*lumirpc.CreateResponse, error) {
    contract.Assert(req.GetType() == string(InstanceToken))
    obj, _, decerr := p.Unmarshal(req.GetProperties())
    if decerr != nil {
        return nil, decerr
    }
    id, err := p.ops.Create(ctx, obj)
    if err != nil {
        return nil, err
    }
    return &lumirpc.CreateResponse{Id: string(id)}, nil
}

func (p *InstanceProvider) Get(
    ctx context.Context, req *lumirpc.GetRequest) (*lumirpc.GetResponse, error) {
    contract.Assert(req.GetType() == string(InstanceToken))
    id := resource.ID(req.GetId())
    obj, err := p.ops.Get(ctx, id)
    if err != nil {
        return nil, err
    }
    return &lumirpc.GetResponse{
        Properties: resource.MarshalProperties(
            nil, resource.NewPropertyMap(obj), resource.MarshalOptions{}),
    }, nil
}

func (p *InstanceProvider) InspectChange(
    ctx context.Context, req *lumirpc.InspectChangeRequest) (*lumirpc.InspectChangeResponse, error) {
    contract.Assert(req.GetType() == string(InstanceToken))
    id := resource.ID(req.GetId())
    old, oldprops, decerr := p.Unmarshal(req.GetOlds())
    if decerr != nil {
        return nil, decerr
    }
    new, newprops, decerr := p.Unmarshal(req.GetNews())
    if decerr != nil {
        return nil, decerr
    }
    var replaces []string
    diff := oldprops.Diff(newprops)
    if diff != nil {
        if diff.Changed("name") {
            replaces = append(replaces, "name")
        }
    }
    more, err := p.ops.InspectChange(ctx, id, old, new, diff)
    if err != nil {
        return nil, err
    }
    return &lumirpc.InspectChangeResponse{
        Replaces: append(replaces, more...),
    }, err
}

func (p *InstanceProvider) Update(
    ctx context.Context, req *lumirpc.UpdateRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(InstanceToken))
    id := resource.ID(req.GetId())
    old, oldprops, err := p.Unmarshal(req.GetOlds())
    if err != nil {
        return nil, err
    }
    new, newprops, err := p.Unmarshal(req.GetNews())
    if err != nil {
        return nil, err
    }
    diff := oldprops.Diff(newprops)
    if err := p.ops.Update(ctx, id, old, new, diff); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *InstanceProvider) Delete(
    ctx context.Context, req *lumirpc.DeleteRequest) (*pbempty.Empty, error) {
    contract.Assert(req.GetType() == string(InstanceToken))
    id := resource.ID(req.GetId())
    if err := p.ops.Delete(ctx, id); err != nil {
        return nil, err
    }
    return &pbempty.Empty{}, nil
}

func (p *InstanceProvider) Unmarshal(
    v *pbstruct.Struct) (*Instance, resource.PropertyMap, mapper.DecodeError) {
    var obj Instance
    props := resource.UnmarshalProperties(nil, v, resource.MarshalOptions{RawResources: true})
    result := mapper.MapIU(props.Mappable(), &obj)
    return &obj, props, result
}

/* Marshalable Instance structure(s) */

// Instance is a marshalable representation of its corresponding IDL type.
type Instance struct {
    Name *string `json:"name,omitempty"`
    ImageID string `json:"imageId"`
    InstanceType *InstanceType `json:"instanceType,omitempty"`
    SecurityGroups *[]resource.ID `json:"securityGroups,omitempty"`
    KeyName *string `json:"keyName,omitempty"`
    AvailabilityZone string `json:"availabilityZone,omitempty"`
    PrivateDNSName *string `json:"privateDNSName,omitempty"`
    PublicDNSName *string `json:"publicDNSName,omitempty"`
    PrivateIP *string `json:"privateIP,omitempty"`
    PublicIP *string `json:"publicIP,omitempty"`
}

// Instance's properties have constants to make dealing with diffs and property bags easier.
const (
    Instance_Name = "name"
    Instance_ImageID = "imageId"
    Instance_InstanceType = "instanceType"
    Instance_SecurityGroups = "securityGroups"
    Instance_KeyName = "keyName"
    Instance_AvailabilityZone = "availabilityZone"
    Instance_PrivateDNSName = "privateDNSName"
    Instance_PublicDNSName = "publicDNSName"
    Instance_PrivateIP = "privateIP"
    Instance_PublicIP = "publicIP"
)

/* Typedefs */

type (
    InstanceType string
)

/* Constants */

const (
    C3Instance2XLarge InstanceType = "c3.2xlarge"
    C3Instance4XLarge InstanceType = "c3.4xlarge"
    C3Instance8XLarge InstanceType = "c3.8xlarge"
    C3InstanceLarge InstanceType = "c3.large"
    C3InstanceXLarge InstanceType = "c3.xlarge"
    C4Instance2XLarge InstanceType = "c4.2xlarge"
    C4Instance4XLarge InstanceType = "c4.4xlarge"
    C4Instance8XLarge InstanceType = "c4.8xlarge"
    C4InstanceLarge InstanceType = "c4.large"
    C4InstanceXLarge InstanceType = "c4.xlarge"
    D2Instance2XLarge InstanceType = "d2.2xlarge"
    D2Instance4XLarge InstanceType = "d2.4xlarge"
    D2Instance8XLarge InstanceType = "d2.8xlarge"
    D2InstanceXLarge InstanceType = "d2.xlarge"
    F1Instance16XLarge InstanceType = "f1.16xlarge"
    F1Instance2XLarge InstanceType = "f1.2xlarge"
    G2Instance2XLarge InstanceType = "g2.2xlarge"
    G2Instance8XLarge InstanceType = "g2.8xlarge"
    I3Instance16XLarge InstanceType = "i3.16xlarge"
    I3Instance2XLarge InstanceType = "i3.2xlarge"
    I3Instance4XLarge InstanceType = "i3.4xlarge"
    I3Instance8XLarge InstanceType = "i3.8xlarge"
    I3InstanceLarge InstanceType = "i3.large"
    I3InstanceXLarge InstanceType = "i3.xlarge"
    M3Instance2XLarge InstanceType = "m3.2xlarge"
    M3InstanceLarge InstanceType = "m3.large"
    M3InstanceMedium InstanceType = "m3.medium"
    M3InstanceXLarge InstanceType = "m3.xlarge"
    M4Instance10XLarge InstanceType = "m4.10xlarge"
    M4Instance16XLarge InstanceType = "m4.16xlarge"
    M4Instance2XLarge InstanceType = "m4.2xlarge"
    M4Instance4XLarge InstanceType = "m4.4xlarge"
    M4InstanceLarge InstanceType = "m4.large"
    M4InstanceXLarge InstanceType = "m4.xlarge"
    P2Instance16XLarge InstanceType = "p2.16xlarge"
    P2Instance8XLarge InstanceType = "p2.8xlarge"
    P2InstanceXLarge InstanceType = "p2.xlarge"
    R3Instance2XLarge InstanceType = "r3.2xlarge"
    R3Instance4XLarge InstanceType = "r3.4xlarge"
    R3Instance8XLarge InstanceType = "r3.8xlarge"
    R3InstanceLarge InstanceType = "r3.large"
    R3InstanceXLarge InstanceType = "r3.xlarge"
    R4Instance16XLarge InstanceType = "r4.16xlarge"
    R4Instance2XLarge InstanceType = "r4.2xlarge"
    R4Instance4XLarge InstanceType = "r4.4xlarge"
    R4Instance8XLarge InstanceType = "r4.8xlarge"
    R4InstanceLarge InstanceType = "r4.large"
    R4InstanceXLarge InstanceType = "r4.xlarge"
    T2Instance2XLarge InstanceType = "t2.2xlarge"
    T2InstanceLarge InstanceType = "t2.large"
    T2InstanceMedium InstanceType = "t2.medium"
    T2InstanceMicro InstanceType = "t2.micro"
    T2InstanceNano InstanceType = "t2.nano"
    T2InstanceSmall InstanceType = "t2.small"
    T2InstanceXLarge InstanceType = "t2.xlarge"
    X1Instance16XLarge InstanceType = "x1.16xlarge"
    X1Instance32XLarge InstanceType = "x1.32xlarge"
)


