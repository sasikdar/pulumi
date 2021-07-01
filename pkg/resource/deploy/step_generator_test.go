package deploy

import (
	"github.com/pulumi/pulumi/pkg/v3/resource/graph"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"os"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/stretchr/testify/assert"
)

func TestIgnoreChanges(t *testing.T) {
	cases := []struct {
		name          string
		oldInputs     map[string]interface{}
		newInputs     map[string]interface{}
		expected      map[string]interface{}
		ignoreChanges []string
		expectFailure bool
	}{
		{
			name: "Present in old and new sets",
			oldInputs: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "foo",
				},
			},
			newInputs: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "bar",
				},
				"c": 42,
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "foo",
				},
				"c": 42,
			},
			ignoreChanges: []string{"a.b"},
		},
		{
			name: "Missing in new sets",
			oldInputs: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "foo",
				},
			},
			newInputs: map[string]interface{}{
				"a": map[string]interface{}{},
				"c": 42,
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "foo",
				},
				"c": 42,
			},
			ignoreChanges: []string{"a.b"},
		},
		{
			name:      "Missing in old deletes",
			oldInputs: map[string]interface{}{},
			newInputs: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "foo",
				},
				"c": 42,
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{},
				"c": 42,
			},
			ignoreChanges: []string{"a.b"},
		},
		{
			name:      "Missing keys in old and new are OK",
			oldInputs: map[string]interface{}{},
			newInputs: map[string]interface{}{},
			ignoreChanges: []string{
				"a",
				"a.b",
				"a.c[0]",
			},
		},
		{
			name: "Missing parent keys in only new fail",
			oldInputs: map[string]interface{}{
				"a": map[string]interface{}{
					"b": "foo",
				},
			},
			newInputs:     map[string]interface{}{},
			ignoreChanges: []string{"a.b"},
			expectFailure: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			olds, news := resource.NewPropertyMapFromMap(c.oldInputs), resource.NewPropertyMapFromMap(c.newInputs)

			expected := olds
			if c.expected != nil {
				expected = resource.NewPropertyMapFromMap(c.expected)
			}

			processed, res := processIgnoreChanges(news, olds, c.ignoreChanges)
			if c.expectFailure {
				assert.NotNil(t, res)
			} else {
				assert.Nil(t, res)
				assert.Equal(t, expected, processed)
			}
		})
	}
}

func createResource(name string) *resource.State {
	ty := tokens.Type("test")
	return &resource.State{
		Type:    ty,
		URN:     resource.NewURN(tokens.QName("teststack"), tokens.PackageName("pkg"), "", ty, tokens.QName(name)),
		Inputs:  make(resource.PropertyMap),
		Outputs: make(resource.PropertyMap),
	}
}

func createCustomResource(name string) *resource.State {
	ty := tokens.Type("test")
	return &resource.State{
		Type:    ty,
		URN:     resource.NewURN(tokens.QName("teststack"), tokens.PackageName("pkg"), "", ty, tokens.QName(name)),
		Inputs:  make(resource.PropertyMap),
		Outputs: make(resource.PropertyMap),
		Custom:  true,
	}
}

func TestIsWish(t *testing.T) {
	stack := tokens.QName("testbook")
	pckgName := tokens.PackageName("foo/bar/baz")
	parentType := tokens.Type("")
	typ := tokens.Type("pulumi-nodejs:dynamic:Resource")
	name := tokens.QName("wish$Example")
	urn := resource.NewURN(stack, pckgName, parentType, typ, name)
	assert.Equal(t, true, isWish(urn), "The URN is not a wish resource URN")

}

func TestUnsatifiedWishSkippedCreate(t *testing.T) {

	//create two resources, childResource is dependent on dependcyResource,
	//which is a wish-resource
	childResource := createResource("testresource")
	dependcyResource := &resource.State{
		Type:    tokens.Type("pulumi-nodejs:dynamic:Resource"),
		URN:     resource.NewURN(tokens.QName("teststack"), tokens.PackageName("pkg"), "", tokens.Type("pulumi-nodejs:dynamic:Resource"), tokens.QName("wish$Resourcedepency")),
		Inputs:  make(resource.PropertyMap),
		Outputs: make(resource.PropertyMap),
	}
	ma1 := map[string]interface{}{"isSatnewResourceisfied": false}
	dependcyResource.Outputs = resource.NewPropertyMapFromMap(ma1)

	//create a deployment goal for with the child resource.
	testgoal := resource.NewGoal(childResource.URN.Type(), "testresource", false,
		resource.PropertyMap{}, childResource.Parent, false, []resource.URN{dependcyResource.URN},
		"", []string{}, nil, nil,
		nil, nil, nil, "", nil)

	//create snapshot of this resource
	initsnap := newSnapshot([]*resource.State{}, []resource.Operation{})

	//prepare a dummy deployment
	//prepare the context of the deployment.
	//create a deployment which has dependencyResource as news
	cwd, _ := os.Getwd()
	sink := cmdutil.Diag()
	ctx, err := plugin.NewContext(sink, sink, nil, nil, cwd, nil, true, nil)
	runInfo := &EvalRunInfo{
		Proj:   &workspace.Project{Name: "pkg"},
		Target: &Target{Name: "teststack"},
	}
	src := NewEvalSource(ctx, runInfo, nil, false)

	if err != nil {
		print("error creating context")
	}
	deployment := &Deployment{
		ctx:                  ctx,
		target:               runInfo.Target,
		prev:                 initsnap,
		olds:                 nil,
		imports:              nil,
		isImport:             false,
		schemaLoader:         nil,
		source:               src,
		localPolicyPackPaths: nil,
		preview:              true,
		depGraph:             nil,
		providers:            nil,
		news:                 &resourceMap{},
		goals:                &goalMap{},
	}
	deployment.news.set(dependcyResource.URN, dependcyResource)

	//create dummy stepgen
	sg := &stepGenerator{
		deployment:        deployment,
		urns:              map[resource.URN]bool{},
		sames:             map[resource.URN]bool{},
		skippedCreates:    map[resource.URN]bool{},
		updates:           map[resource.URN]bool{},
		deletes:           map[resource.URN]bool{},
		replaces:          map[resource.URN]bool{},
		replaceTargetsOpt: map[resource.URN]bool{},
		updateTargetsOpt:  map[resource.URN]bool{},
	}

	event := registerResourceEvent{goal: testgoal}
	steps, _ := sg.generateSteps(&event)
	if len(steps) == 1 {
		assert.Equal(t, sg.skippedCreates[childResource.URN], true)
	}

	//verify that resource will be deleted if it existed in old
	deployment.olds = map[resource.URN]*resource.State{childResource.URN: childResource}
	eventwithold := registerResourceEvent{goal: testgoal}
	stepsold, _ := sg.generateSteps(&eventwithold)
	if len(stepsold) == 1 {
		assert.Equal(t, sg.skippedCreates[childResource.URN], true)
		assert.Equal(t, sg.deletes[childResource.URN], true)

	}
}

func TestWishDependentResourcePruning(t *testing.T) {
	//first create some resources create dependency among them
	//create dependency graph and link it to a deployment
	//create a step_generator and link the deployment instance
	//call the calculateDependentReplacements using step_generator instance

	//Create Resource and dependency, dependency structure:
	// 			     root
	//           /    |    \
	//       res2    res3   wish$res1
	//                        |
	//                       dep1
	//                       /   \
	//                  dep11    dep12
	//end of diag

	res1 := createCustomResource("root")

	res2 := createCustomResource("res2")
	res2.Dependencies = []resource.URN{res1.URN}

	res3 := createCustomResource("res3")
	res3.Dependencies = []resource.URN{res1.URN}

	//create a wish resource
	res4 := &resource.State{
		Type:    tokens.Type("pulumi-nodejs:dynamic:Resource"),
		URN:     resource.NewURN(tokens.QName("teststack"), tokens.PackageName("pkg"), "", tokens.Type("pulumi-nodejs:dynamic:Resource"), tokens.QName("wish$res1")),
		Inputs:  make(resource.PropertyMap),
		Outputs: make(resource.PropertyMap),
	}

	res4.Dependencies = []resource.URN{res1.URN}

	res5 := createCustomResource("dep1")
	res5.Dependencies = []resource.URN{res4.URN}

	res6 := createCustomResource("dep11")
	res6.Dependencies = []resource.URN{res5.URN}

	res7 := createCustomResource("dep12")
	res7.Dependencies = []resource.URN{res5.URN}

	//create a dependency graph
	depGraph := graph.NewDependencyGraph([]*resource.State{res1, res2, res3, res4, res5, res6, res7})

	//prepare a deployment
	initsnap := newSnapshot([]*resource.State{}, []resource.Operation{})

	//prepare a dummy deployment
	//prepare the context of the deployment.
	//create a deployment which has dependencyResource as news
	cwd, _ := os.Getwd()
	sink := cmdutil.Diag()
	ctx, err := plugin.NewContext(sink, sink, nil, nil, cwd, nil, true, nil)
	runInfo := &EvalRunInfo{
		Proj:   &workspace.Project{Name: "pkg"},
		Target: &Target{Name: "teststack"},
	}
	src := NewEvalSource(ctx, runInfo, nil, false)

	if err != nil {
		print("error creating context")
	}
	deployment := &Deployment{
		ctx:                  ctx,
		target:               runInfo.Target,
		prev:                 initsnap,
		olds:                 nil,
		imports:              nil,
		isImport:             false,
		schemaLoader:         nil,
		source:               src,
		localPolicyPackPaths: nil,
		preview:              true,
		depGraph:             depGraph,
		providers:            nil,
		news:                 &resourceMap{},
		goals:                &goalMap{},
	}

	//create a step_generator instance and link the deployment with depGraph
	//calculateDependentReplacements will need the deployment instance of SG to access depGraph
	sg := &stepGenerator{
		deployment:        deployment,
		urns:              map[resource.URN]bool{},
		sames:             map[resource.URN]bool{},
		skippedCreates:    map[resource.URN]bool{},
		updates:           map[resource.URN]bool{},
		deletes:           map[resource.URN]bool{},
		replaces:          map[resource.URN]bool{},
		replaceTargetsOpt: map[resource.URN]bool{},
		updateTargetsOpt:  map[resource.URN]bool{},
	}
	dependentset, _ := sg.calculateDependentReplacements(res4)
	assert.Equal(t, len(dependentset), 3)
	assert.Equal(t, dependentset[0].res.URN, res5.URN)
	assert.Equal(t, dependentset[1].res.URN, res6.URN)
	assert.Equal(t, dependentset[2].res.URN, res7.URN)

}
