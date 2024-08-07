package main

import (
	"context"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/tools/flow"
	"k8s.io/klog/v2"

	"github.com/penk110/k8s_operator/deployment_1/handler"
)

const (
	K8SFlowTpl = "../flow_templates/deploy_flow.cue"
)

func main() {
	inst := load.Instances([]string{K8SFlowTpl}, nil)[0]

	cc := cuecontext.New()
	cv := cc.BuildInstance(inst)

	flowConfig := &flow.Config{
		Root: cue.ParsePath(handler.K8sTest1Root),
	}
	k8sFlow := flow.New(flowConfig, cv, handler.Handler)

	err := k8sFlow.Run(context.TODO())
	if err != nil {
		klog.Errorf("k8sFlow err: %v", err)
		return
	}
}
