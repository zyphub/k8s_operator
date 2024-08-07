package main

import (
	"bytes"
	"context"
	encodingjson "encoding/json"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/tools/flow"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/klog/v2"

	"github.com/penk110/k8s_operator/deployment_1/handler"
	"github.com/penk110/k8s_operator/k8s_client"
)

const (
	K8SFlowTpl = "../flow_templates/deploy_flow.cue"
)

func main() {

	apiGroupResources, err := k8s_client.RestMapper()
	if err != nil {
		klog.Errorf("RestNapper err: %v", err)
		return
	}

	restMapping, err := apiGroupResources.RESTMapping(schema.GroupKind{
		Group: "apps",
		Kind:  "Deployment",
	})
	if err != nil {
		klog.Errorf("RESTMapping err: %v", err)
		return
	}
	restMappingJson, _ := encodingjson.Marshal(restMapping)
	var prettyJSON bytes.Buffer
	_ = encodingjson.Indent(&prettyJSON, restMappingJson, "", "  ")

	klog.Infof("restMapping: %v", prettyJSON.String())

	inst := load.Instances([]string{K8SFlowTpl}, nil)[0]

	cc := cuecontext.New()
	cv := cc.BuildInstance(inst)

	flowConfig := &flow.Config{
		Root: cue.ParsePath(handler.K8sTest1Root),
	}
	k8sFlow := flow.New(flowConfig, cv, handler.Handler)

	err = k8sFlow.Run(context.TODO())
	if err != nil {
		klog.Errorf("k8sFlow err: %v", err)
		return
	}

	k8sJson, err := k8sFlow.Value().MarshalJSON()

	err = k8s_client.Apply(k8sJson, k8s_client.GetConfig(), apiGroupResources)
	if err != nil {
		klog.Errorf("Apply err: %v", err)
		return
	}
}
