package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/tools/flow"
	"github.com/gin-gonic/gin"
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
	go func() {
		if err := regFlow.Run(context.TODO()); err != nil {
			log.Fatalln(err)
		}
	}()

	regFlow.Tasks()[0].State()
	r := gin.New()
	r.LoadHTMLGlob("workflow/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"tasks": regFlow.Tasks()})
	})
	r.Run(":8080")

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
