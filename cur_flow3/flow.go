package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/tools/flow"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

const tasks = `
    reg: {
		uname: "tester"
		pass: "12345"
	}
	regrsp: {
		regname: reg.uname
		result: string
	}
	test:{
		data: int
	}

     `

// http://127.0.0.1:8080/

func main() {
	cc := cuecontext.New()
	cv := cc.CompileString(tasks)
	// cue 工作流对象
	regFlow := flow.New(nil, cv, regFlowFunc)
	flowCtx, flowCancel := context.WithCancel(context.Background())

	regFlow.Tasks()[0].State()
	r := gin.New()
	r.LoadHTMLGlob("workflow/*")

	// 1、访问主页，回去当前的未执行的工作流的状态信息
	// 2、启动工作流，重定向会主页，刷新获取当前工作流状态
	// 3、重置工作流，则重新创建flow对象
	// 4、取消，则退出工作流

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"tasks": regFlow.Tasks()})
	})

	r.GET("/reset", func(c *gin.Context) {
		regFlow = flow.New(nil, cv, regFlowFunc)
		c.Redirect(http.StatusFound, "/")
		return
	})

	r.GET("/cancel", func(c *gin.Context) {
		flowCancel()
		c.Redirect(http.StatusFound, "/")
		return
	})

	r.POST("/", func(c *gin.Context) {

		go func() {
			err := regFlow.Run(flowCtx)
			if err != nil {
				klog.Errorf("regFlow err: %v", err)
				c.JSONP(200, gin.H{
					"code": 5000000,
					"msg":  err.Error(),
					"data": nil,
				})
				return
			}
		}()

		c.Redirect(http.StatusFound, "/")

		return
	})

	err := r.Run(":8080")
	if err != nil {
		klog.Fatalf("r.Run err: %v", err)
		return
	}

}
func regFlowFunc(v cue.Value) (flow.Runner, error) {
	uname := v.LookupPath(cue.ParsePath("uname"))
	result := v.LookupPath(cue.ParsePath("result"))
	if !uname.Exists() && !result.Exists() {
		return nil, nil
	}
	return flow.RunnerFunc(func(t *flow.Task) error {
		if t.Path().String() == "reg" {
			//假设模拟 数据库很耗时
			time.Sleep(time.Second * 5)
			unameStr, err := uname.String()
			if err != nil {
				return err
			}
			if unameStr == "admin" {
				return fmt.Errorf("不能注册为admin用户名")
			}
			return nil
		}

		lastUname, err := t.Value().LookupPath(cue.ParsePath("regname")).String()
		fmt.Println(t.Value())
		if err != nil {
			return err
		}
		return t.Fill(map[string]interface{}{
			"result": lastUname + "用户注册成功",
		})

	}), nil
}
