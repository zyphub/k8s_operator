package handler

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/tools/flow"
	"k8s.io/klog/v2"
)

const K8sTest1Root = "workflow" // 代表 根节点

func Handler(v cue.Value) (flow.Runner, error) {
	l, b := v.Label()

	if !b || l == K8sTest1Root {
		klog.Warningf("l: %v, b: %v", l, b)
		return nil, nil
	}

	// 实际处理的模板的函数体，需要在这里面 调用k8s api 执行资源部署 的逻辑
	return flow.RunnerFunc(func(t *flow.Task) error {
		fmt.Println("工作流节点", t.Path())

		k8sJson, err := t.Value().MarshalJSON()
		if err != nil {
			klog.Errorf("t.Value().MarshalJSON() err: %v", err)
			return err
		}

		klog.Infof("k8sJson: %v", string(k8sJson))

		//err = utils.K8sApply(k8sJson, utils.K8sRestConfig, *utils.K8sRestMapper)
		//if err != nil {
		//	return err
		//}

		return nil
	}), nil
}
