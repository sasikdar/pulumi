package lifecycletest

import (
	"fmt"
	"github.com/blang/semver"
	. "github.com/pulumi/pulumi/pkg/v3/engine"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy/deploytest"
	"github.com/pulumi/pulumi/pkg/v3/resource/deploy/providers"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/result"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsatisfiedWishCreationSkipped(t *testing.T) {

	p := &TestPlan{}
	//prepare the loader
	loaders := []*deploytest.ProviderLoader{
		deploytest.NewProviderLoader("pkgA", semver.MustParse("1.0.0"), func() (plugin.Provider, error) {
			return &deploytest.Provider{}, nil
		}),
	}

	//define type of the two variables

	const resType = "pkgA:index:typ"

	//define a function to create new resource
	type propertyDependencies map[resource.PropertyKey][]resource.URN
	newResource := func(urn resource.URN, id resource.ID, provider string, dependencies []resource.URN,
		propertyDeps propertyDependencies, outputs resource.PropertyMap) *resource.State {

		inputs := resource.PropertyMap{}
		for k := range propertyDeps {
			inputs[k] = resource.NewStringProperty("foo")
		}
		return &resource.State{
			Type:                 urn.Type(),
			URN:                  urn,
			Custom:               true,
			Delete:               false,
			ID:                   id,
			Inputs:               inputs,
			Outputs:              outputs,
			Dependencies:         dependencies,
			Provider:             provider,
			PropertyDependencies: propertyDeps,
		}
	}

	//create two resources

	//resProvider:=providers.MakeProviderType("pkgA")
	reswishUrn := resource.NewURN(tokens.QName("teststack"), "pkgA", "", tokens.Type("pulumi-nodejs:dynamic:Resource"), tokens.QName("wish$Resource"))
	resWish := newResource(reswishUrn, "1", "", []resource.URN{}, nil,
		resource.NewPropertyMapFromMap(map[string]interface{}{"isSatnewResourceisfied": false}))

	resbUrn := p.NewURN(tokens.Type(resType), "resB", "")
	resB := newResource(resbUrn, "1", "", []resource.URN{resWish.URN}, nil, nil)

	//create the environment and run the deployment
	//define variables
	//var provURN, urnA resource.URN
	//var provID resource.ID
	//var err error
	//prepare the dummy program
	//create a validate function as part of test steps

	fmt.Printf("%q", resWish.URN.Type().String())
	fmt.Printf("%q", resWish.URN.Name())
	var provURN, urnA, urnB resource.URN
	var provID resource.ID
	var err error

	//create a program which registers the intended resources
	program := deploytest.NewLanguageRuntime(func(_ plugin.RunInfo, monitor *deploytest.ResourceMonitor) error {

		provURN, provID, _, err = monitor.RegisterResource(providers.MakeProviderType("pkgA"), "provA", true)
		assert.NoError(t, err)

		if provID == "" {
			provID = providers.UnknownID
		}
		provRef, err := providers.NewReference(provURN, provID)
		assert.NoError(t, err)
		provA := provRef.String()

		urnA, _, _, err = monitor.RegisterResource(resWish.URN.Type(), resWish.URN.Name().String(), true, deploytest.ResourceOptions{
			Provider: provA,
			Inputs:   nil,
		})
		assert.NoError(t, err)

		errOut := monitor.RegisterResourceOutputs(urnA, resource.NewPropertyMapFromMap(map[string]interface{}{"isSatisfied": false}))
		assert.NoError(t, errOut)

		urnB, _, _, err = monitor.RegisterResource(resB.URN.Type(), resB.URN.Name().String(), false, deploytest.ResourceOptions{
			Provider:     provA,
			Dependencies: []resource.URN{urnA},
		})

		assert.NoError(t, err)

		return nil
	})

	p.Options.Host = deploytest.NewPluginHost(nil, nil, program, loaders...)

	//create the test steps and the validation function
	p.Steps = []TestStep{{Op: Update}}

	p.Steps = []TestStep{{
		Op: Update,

		Validate: func(project workspace.Project, target deploy.Target, entries JournalEntries,
			evts []Event, res result.Result) result.Result {

			assert.Nil(t, res)

			AssertSameSteps(t, []StepSummary{
				{Op: deploy.OpCreate, URN: provURN},
				{Op: deploy.OpCreate, URN: urnA},
				{Op: deploy.OpSame, URN: urnB},
			}, SuccessfulSteps(entries))

			return res
		},
	}}
	snap := p.Run(t, nil)

	//use the snap shot as old
	//try to again register the two resources
	//try to access the area of step generator which checks for hasold for wishresources
	//end of test case

	//create a program which registers the intended resources
	program1 := deploytest.NewLanguageRuntime(func(_ plugin.RunInfo, monitor *deploytest.ResourceMonitor) error {

		provURN, provID, _, err = monitor.RegisterResource(providers.MakeProviderType("pkgA"), "provA", true)
		assert.NoError(t, err)

		if provID == "" {
			provID = providers.UnknownID
		}
		provRef, err := providers.NewReference(provURN, provID)
		assert.NoError(t, err)
		provA := provRef.String()

		urnA, _, _, err = monitor.RegisterResource(resWish.URN.Type(), resWish.URN.Name().String(), true, deploytest.ResourceOptions{
			Provider: provA,
			Inputs:   nil,
		})
		assert.NoError(t, err)

		errOut := monitor.RegisterResourceOutputs(urnA, resource.NewPropertyMapFromMap(map[string]interface{}{"isSatisfied": false}))
		assert.NoError(t, errOut)

		urnB, _, _, err = monitor.RegisterResource(resB.URN.Type(), resB.URN.Name().String(), false, deploytest.ResourceOptions{
			Provider:     provA,
			Dependencies: []resource.URN{urnA},
		})

		assert.NoError(t, err)

		return nil
	})

	p.Steps = []TestStep{{
		Op: Update,

		Validate: func(project workspace.Project, target deploy.Target, entries JournalEntries,
			evts []Event, res result.Result) result.Result {

			assert.Nil(t, res)
			resBDeleted, resBSame := false, false
			for _, entry := range entries {
				if entry.Step.URN() == urnB {
					switch entry.Step.Op() {
					case deploy.OpDelete, deploy.OpDeleteReplaced:
						resBDeleted = true
					case deploy.OpSame:
						resBSame = true
					}
				}
			}
			assert.True(t, resBSame)
			assert.True(t, resBDeleted)

			return res
		},
	}}
	p.Options.Host = deploytest.NewPluginHost(nil, nil, program1, loaders...)
	_ = p.Run(t, snap)

}

func TestWishResourceDependentDelete(t *testing.T) {

	p := &TestPlan{}

	urns, old, program := ComplexResourceGraphWithWishResource(t, p)
	loaders := []*deploytest.ProviderLoader{
		deploytest.NewProviderLoader("pkgA", semver.MustParse("1.0.0"), func() (plugin.Provider, error) {
			return &deploytest.Provider{}, nil
		}),
	}
	p.Options.Host = deploytest.NewPluginHost(nil, nil, program, loaders...)

	destroyTargets := []resource.URN{p.NewURN(tokens.Type("pulumi-nodejs:dynamic:Resource"), "wish$resA", "")}
	p.Options.DestroyTargets = destroyTargets
	p.Options.TargetDependents = true
	t.Logf("Destroying targets: %v", destroyTargets)

	p.Steps = []TestStep{{
		Op: Destroy,
		//ExpectFailure: !targetDependents,
		Validate: func(project workspace.Project, target deploy.Target, entries JournalEntries,
			evts []Event, res result.Result) result.Result {

			assert.True(t, len(entries) > 0)

			deleted := make(map[resource.URN]bool)
			for _, entry := range entries {
				assert.Equal(t, deploy.OpDelete, entry.Step.Op())
				deleted[entry.Step.URN()] = true
			}

			for _, target := range p.Options.DestroyTargets {
				assert.Contains(t, deleted, target)
			}
			assert.True(t, deleted[urns[6]])
			assert.True(t, deleted[urns[7]])
			assert.True(t, deleted[urns[8]])
			assert.True(t, deleted[urns[9]])

			return res
		},
	}}

	p.Run(t, old)

}

func ComplexResourceGraphWithWishResource(t *testing.T, p *TestPlan) ([]resource.URN, *deploy.Snapshot, plugin.LanguageRuntime) {

	// we are trying to create a ResourceGraph which resembeles the following structure
	// 			     resA
	//           /    |    \
	//       resB    resC   wish$resA
	//                        |
	//                       depA
	//                       /   \
	//                  depB    depC

	resType := tokens.Type("pkgA:m:typA")
	wishType := tokens.Type("pulumi-nodejs:dynamic:Resource")

	type propertyDependencies map[resource.PropertyKey][]resource.URN

	//create the providers
	urnProva := p.NewProviderURN("pkgA", "provA", "")
	urnProvb := p.NewProviderURN("pkgA", "provB", "")
	urnProvc := p.NewProviderURN("pkgA", "provC", "")
	urnProvwish := p.NewProviderURN("pkgA", "provWish", "")

	//create the resources
	urnResa := p.NewURN(resType, "resA", "")
	urnResb := p.NewURN(resType, "resB", "")
	urnResc := p.NewURN(resType, "resC", "")
	urnWishresa := p.NewURN(wishType, "wish$resA", "")
	urnDepa := p.NewURN(resType, "depA", "")
	urnDepb := p.NewURN(resType, "depB", "")
	urnDepc := p.NewURN(resType, "depC", "")

	urns := []resource.URN{
		urnProva, urnProvb, urnProvc,
		urnResa, urnResb, urnResc, urnWishresa, urnDepa, urnDepb,
		urnDepc,
	}

	newResource := func(urn resource.URN, id resource.ID, provider string, dependencies []resource.URN,
		propertyDeps propertyDependencies, outputs resource.PropertyMap) *resource.State {

		inputs := resource.PropertyMap{}
		for k := range propertyDeps {
			inputs[k] = resource.NewStringProperty("foo")
		}

		return &resource.State{
			Type:                 urn.Type(),
			URN:                  urn,
			Custom:               true,
			Delete:               false,
			ID:                   id,
			Inputs:               inputs,
			Outputs:              outputs,
			Dependencies:         dependencies,
			Provider:             provider,
			PropertyDependencies: propertyDeps,
		}
	}

	old := &deploy.Snapshot{
		Resources: []*resource.State{
			newResource(urnProva, "0", "", nil, nil, nil),
			newResource(urnProvb, "1", "", nil, nil, nil),
			newResource(urnProvc, "2", "", nil, nil, nil),
			newResource(urnProvwish, "3", "", nil, nil, nil),
			newResource(urnResa, "4", string(urnProva)+"::0", nil, nil, nil),
			newResource(urnResb, "5", string(urnProvb)+"::1", []resource.URN{urnResa}, nil, nil),
			newResource(urnResc, "6", string(urnProvc)+"::2", []resource.URN{urnResa}, nil, nil),
			newResource(urnWishresa, "7", string(urnProvwish)+"::3", []resource.URN{urnResa}, nil, nil),
			newResource(urnDepa, "8", "", []resource.URN{urnWishresa}, nil, nil),
			newResource(urnDepb, "9", "", []resource.URN{urnDepa}, nil, nil),
			newResource(urnDepc, "10", "", []resource.URN{urnDepa}, nil, nil),
		},
	}

	program := deploytest.NewLanguageRuntime(func(_ plugin.RunInfo, monitor *deploytest.ResourceMonitor) error {
		register := func(urn resource.URN, provider string, inputs resource.PropertyMap) resource.ID {
			_, id, _, err := monitor.RegisterResource(urn.Type(), string(urn.Name()), true, deploytest.ResourceOptions{

				Provider: provider,
				Inputs:   inputs,
			})
			assert.NoError(t, err)
			return id
		}
		//register providers
		idA := register(urnProva, "", nil)
		idB := register(urnProvb, "", nil)
		idC := register(urnProvc, "", nil)

		//register resources
		register(urnResa, string(urnProva)+"::"+string(idA), nil)
		register(urnResb, string(urnProvb)+"::"+string(idB), nil)
		register(urnResc, string(urnProvc)+"::"+string(idC), nil)
		register(urnWishresa, string(urnProvwish)+"::3", nil)
		register(urnDepa, "", nil)
		register(urnDepb, "", nil)
		register(urnDepc, "", nil)

		return nil
	})

	return urns, old, program
}
