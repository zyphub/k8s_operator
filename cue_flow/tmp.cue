register: {
	// 模板学习
	username: "tester"
	password: "tester"
}

confire: {
	result: register.username
}

resp: {
	last_username: confire.result
	result:   string
}
