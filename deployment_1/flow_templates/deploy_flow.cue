package demploy_flow

import (
	"deployment.com/yamls"
)

// 引用模板
workflow: {
	step1: yamls.deployment
	step2: yamls.service
}
