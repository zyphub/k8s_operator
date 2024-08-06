package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/tools/flow"
	"k8s.io/klog/v2"
)

var curFlowPath string

func main() {

	flag.StringVar(&curFlowPath, "c", "", "path to cue flow")
	flag.Parse()

	klog.Infof("cue file: %v", curFlowPath)

	tmpData, err := ioutil.ReadFile(curFlowPath)
	if err != nil {
		klog.Fatalf("failed to read .cue file: %v", err)
		return
	}

	cueCtx := cuecontext.New()

	cueFlowValue := cueCtx.CompileBytes(tmpData)

	// 打印编译后的cue
	klog.Infof("cue flow value: %v", cueFlowValue)

	// cue flow
	flowConfig := &flow.Config{}
	cueFlow := flow.New(flowConfig, cueFlowValue, CueFlowTaskFunc)

	go func() {
		err = cueFlow.Run(context.Background())
		if err != nil {
			klog.Fatalf("failed to run flow: %v", err)
		}
	}()

	exitCh := make(chan os.Signal, 1)
	go func() {
		for {
			for _, task := range cueFlow.Tasks() {
				if task.State() == 3 && task.Path().String() == "resp" {
					klog.Infof("finish task: %v", task.Value())
					exitCh <- os.Interrupt
					return
				}
				klog.Infof("----- task path: %v state: %v ", task.Path().String(), task.State())
			}
			fmt.Println()
			time.Sleep(1 * time.Second)
		}
	}()

	<-exitCh
}

func CueFlowTaskFunc(cueFlowValue cue.Value) (flow.Runner, error) {
	username := cueFlowValue.LookupPath(cue.ParsePath("username"))
	password := cueFlowValue.LookupPath(cue.ParsePath("password"))
	result := cueFlowValue.LookupPath(cue.ParsePath("result"))

	if !username.Exists() && !result.Exists() {
		//return nil, fmt.Errorf("username or password not found")
		return nil, nil
	}
	funcRet := flow.RunnerFunc(func(t *flow.Task) error {

		if t.Path().String() == "register" {
			usernameValue, err := username.String()
			if err != nil {
				return err
			}
			passwordValue, err := password.String()
			if err != nil {
				return err
			}
			if usernameValue == "admin" {
				return fmt.Errorf("username not allow admin")
			}

			// 模拟业务耗时
			time.Sleep(time.Second * 3)
			klog.Infof("usernameValue: %v passwordValue: %v", usernameValue, passwordValue)
			return nil
		}

		klog.Infof("---- flow path: %v", t.Path().String())

		lastUsername, err := t.Value().LookupPath(cue.ParsePath("last_username")).String()
		if err != nil && t.Path().String() == "resp" { // 最后一个步骤
			return fmt.Errorf("failed to get last last_username: %v", err)
		}
		if t.Path().String() == "resp" {
			return t.Fill(map[string]interface{}{
				"result": lastUsername + "用户注册成功",
			})
		}
		return nil
	})

	return funcRet, nil
}
