package yamls

service: {
	apiVersion: "v1"
	kind:       "Service"
	metadata: name: "flowsvc"
	spec: {
		type: "ClusterIP"
		ports: [
			{
				port:       80
				targetPort: 80
			}]
		selector: {
			// service通过selector和pod建立关联
			app: "flowdeploy"
		}
	}
}
