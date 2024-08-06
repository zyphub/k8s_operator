package k8s_client

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// 测试用的配置
const (
	// ClientCA ClientKey pid的机器
	ClientCA  = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURRakNDQWlxZ0F3SUJBZ0lJTjRkT0RkUzA1eEl3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TXpFd01qQXhNVEkyTlRGYUZ3MHlOREV3TVRreE1UUXdORGRhTURZeApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sc3dHUVlEVlFRREV4SmtiMk5yWlhJdFptOXlMV1JsCmMydDBiM0F3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRQzRGa3VDRFZ3b2Fkb2oKYzlQcFQrM2Njek1rclFMeEdReU5WTVpreTcrUlQyQnF2UFdWN2dqeE0vTHBiUWlIeUUzWHhEeW9sMEQyMjJyeApuUHdtc2ppeWo5eXJtVXhJbmJFam9yNUpnSkVoWExnZWR3N0xYZzlhRm0rNEJkRTBhbEl6bDZXVTdla3QvOGhqCjZicWtMT2ZIdFRCQTVsbU52OWpjbUtjWkRjVHpTUXVlRjNuYzRLS3lUb0R2T1BZSytPdWt2dDB4WWsybDRQWVMKdU91NlRPZitCRGlVWS9PK2hpcVQ4MVE3dWtIRlc1a2NwLzBWRGlVMUM3L21QMnhLU2p1NkRHcHdHSGgraHdFZQpQRjMxNU9mN2N6dmw2cnlEY3JNU1FUSXlDVzJ0TGZ1TGFLeC9iZDV5bHpYYTdGZUluVzcwT1BKbTRXdWZkVnVpCnhPeHEwNUk1QWdNQkFBR2pkVEJ6TUE0R0ExVWREd0VCL3dRRUF3SUZvREFUQmdOVkhTVUVEREFLQmdnckJnRUYKQlFjREFqQU1CZ05WSFJNQkFmOEVBakFBTUI4R0ExVWRJd1FZTUJhQUZMZWxvRTBGaFo0YWdSUnFMZWRPSllFNgpSRnpRTUIwR0ExVWRFUVFXTUJTQ0VtUnZZMnRsY2kxbWIzSXRaR1Z6YTNSdmNEQU5CZ2txaGtpRzl3MEJBUXNGCkFBT0NBUUVBYlphUlNSRFJ3V0p6U1Q0SWFRYnl1UDBNZml0bThlTncvT2x0Yk1ac2w5VllxOWxPa2FkazJnelMKWDNQcnZGUUpBVW9MVk5SSE9kVFhHYXBmVEJWS2JIK1I0MWg0bnZkaDVVcS9iZ1JaWWFxaTAwUmxpajVwcGRPZQpsejFYb3I5TkJLSG94U3UyeDRjd2tzMTJiWDRxbUh5VllSZnJpY1UyR2k4ZmxkUlR2UWpWZnYvdWFDbmxnSzhHCk1rUnQ4YXlZZlJZakF4SkxpVnJmUUhoRCt2bjRFWlZ0VnE5SUxQNGh4QiswZkwyQVpuZTFRK3M0ZTB4VFFiY1oKblBwMjR2WDd6dS91TWwxall3Smk1VmxEZUd3Z3hnMmZ2T1U2dFEyWE5nWTZINXdPZkNGYncwVjZLdmkrdm5EOQpRMXJvUjdVRjNvTzFtaWZBRzczMzNkOXZWTlA1c2c9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
	ClientKey = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBdUJaTGdnMWNLR25hSTNQVDZVL3QzSE16SkswQzhSa01qVlRHWk11L2tVOWdhcnoxCmxlNEk4VFB5NlcwSWg4aE4xOFE4cUpkQTl0dHE4Wno4SnJJNHNvL2NxNWxNU0oyeEk2SytTWUNSSVZ5NEhuY08KeTE0UFdoWnZ1QVhSTkdwU001ZWxsTzNwTGYvSVkrbTZwQ3pueDdVd1FPWlpqYi9ZM0ppbkdRM0U4MGtMbmhkNQozT0Npc2s2QTd6ajJDdmpycEw3ZE1XSk5wZUQyRXJqcnVrem4vZ1E0bEdQenZvWXFrL05VTzdwQnhWdVpIS2Y5CkZRNGxOUXUvNWo5c1Nrbzd1Z3hxY0JoNGZvY0JIanhkOWVUbiszTTc1ZXE4ZzNLekVrRXlNZ2x0clMzN2kyaXMKZjIzZWNwYzEydXhYaUoxdTlEanladUZybjNWYm9zVHNhdE9TT1FJREFRQUJBb0lCQVFDdVlIQW1RWUdLeHJwYgoydHhocGRVcmZmUjBTVzcvOHpwd3BsMUlIYmpaYk5kb1JKWmQ3NTJJM2l5NzhReWprcG9xU1Rrc2VocVB2RWtSCmxpTkVoSTR3bHhYeGRzVk1CQlJJTFdFVFB6WTY1Qm1Fd2tMQllkZ28vaGZWdWF6eWVjUmtHc0krMFI2UTlEcGUKYW9qaCs2ZVRCWTh2NndQcHdsRXFwVytqeStkRWk5RVY0c1V2cTJOWWtOYmNTaGRNQTJjN04yUE1wMk5FVGhuZwpmUFJ6R0JKajZ3WHJGbThQRTlpYWkvVGtZUnM5emdwTGNLWnVTdnNZTEZYSXFyYy9ldVpLSnRJL0VBUXV1a1VxCk9FTmk2OVEvSkxPUktGK0VXSnVKbk1WV1JFZGVZdTNvci9MV2RoQnBHR2VmZ3J6T3EzSlJEalJ3TlVIejRWeFMKbWlMUUNtZ0JBb0dCQU9VYXFMSVR1VmdDdXI0UkpNUTErd0hZdVNEQ0NxZWNCeFhuRWlhUFRhdTdXYUZ5SWIvYwpQVHJGMlR2ZmlVNEhNTDUwbFRVbVVncnNCbTZTN3RocXFmdU40WDRnMys0clNGQjEzc2U1YmMrY1VDVndqNEtXClFLZlROdENTZVM2TjNsSGZmRTlxZkFlT0Ric3hRYWF0VFhqRHMvOTc0Z1JXZlYxYzkrWjFTY1FGQW9HQkFNMnkKdVZ1QTUxaGtyRisyV3ZuS1V6WVRKMk85OFpRdEF5K3VkZHZ1UG8yUVJNOUJ6OUdOU3A3YTRoUzYyTTRqbXZ4VgpPVjhIbW1vYmR1QnFBblZaV1YxeVVNOWkzNVkwNXBkdlhyZm4wU0NqeU9IN0RCMDlNZW1XYmRFTXlTRnJyRDlECjBSRU1FUlhYTEQ5UWpVZDdXQ09aQjAwdkxpeUdsRzMxbS9wdnp6K2xBb0dCQUxta0VHMjdiY1BTOGw1d3Bjb1gKczN5YmorYnJWSmJiNXlIb1N0elQ0YXYxODNyT2NHcDJtMmEwU29JcGI2aTZTdFVJd3A1K25wd2JCRnMwMURTbwp1WFFNVTF0UWFDTWxEME9qUHhHM1B6T3JCWVpRM3ZpQnA0SlZzMlR1U3lOZDhYZUdEOFNLRkZaSzFQV0p4Qmk2CjlMdVdXSlA4WGZnRjNTOTUxYVg0QS8zQkFvR0FONWE0THZsY0MvQlJBU0MzMzArRlExVFR6VW0wc3BXamljdzgKLzYyWDdBdnovSXJOamRVQU9JUHdteWVQbGMzYmdadktnRnIrcVBRNUlSYWxDVytYRGdEcHc5SDFtSk05U2VtSQpFRzB1Z0FLak5DYnpOQ2VvaUhibHdKd1M4dHcxVlhlUFZXc01adm1hZEpYaFNGTVdFN0MwWDNDRHF2Ykh3QnVqCkJvQVc0eDBDZ1lCeEN2d2FKUDBrdTFaemE3Q2trSCtqNmI4WUlORkJ4WERobVlUWFhhOUh6REJIelNaYXlDVzYKbG1VcVUrVVRTRXFOeUF6WVAzVmxzRWRQMmVFQW1LMjNINXZZVmpmSllHQXFodXZ5dUxab21vamNBZGc5UWxXMQpKcnNiRCsxVFFkenJ4VVRJbVBwRTh4KzA0VXRaSm5xc1BtNmRrNGpzZDMyQXJ3d3I4ZGtWU1E9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo="
)

var ClientSet *kubernetes.Clientset
var LocalClientSet *kubernetes.Clientset
var MetricClientSet *versioned.Clientset
var config *rest.Config

func init() {
	var err error
	// 当前机器
	host := "https://127.0.0.1:6443"

	clientCA := ClientCA
	clientKey := ClientKey
	// 环境变量有值直接从环境取
	if os.Getenv("CLOUD_K8S_HOST") != "" && os.Getenv("CLOUD_K8S_ClientCA") != "" && os.Getenv("CLOUD_K8S_ClientKey") != "" {
		// TODO: 临时代码，后期改成先校验网络连通性再判断是否使用环境变量
		CloudK8sHost := os.Getenv("CLOUD_K8S_HOST")
		conn, err := net.Dial("tcp", CloudK8sHost+":6443")
		if err == nil && conn != nil {
			host = fmt.Sprintf("https://%s:6443", os.Getenv("CLOUD_K8S_HOST"))
			log.Println("get env CLOUD_K8S_HOST: " + host)
			clientCA = os.Getenv("CLOUD_K8S_ClientCA")
			clientKey = os.Getenv("CLOUD_K8S_ClientKey")
			log.Println("get env CLOUD_K8S_ClientCA: " + clientCA[:20] + "......")
			log.Println("get env CLOUD_K8S_ClientKey: " + clientKey[:20] + "......")
		}
	}
	keyData, _ := base64.StdEncoding.DecodeString(clientKey)
	certData, _ := base64.StdEncoding.DecodeString(clientCA)

	tlsConfig := rest.TLSClientConfig{
		Insecure: true, // 取消TLS验证
		CertData: certData,
		KeyData:  keyData,
	}
	config = &rest.Config{
		Host:                host,
		APIPath:             "",
		ContentConfig:       rest.ContentConfig{},
		Username:            "",
		Password:            "",
		BearerToken:         "",
		BearerTokenFile:     "",
		Impersonate:         rest.ImpersonationConfig{},
		AuthProvider:        nil,
		AuthConfigPersister: nil,
		ExecProvider:        nil,
		TLSClientConfig:     tlsConfig,
		Timeout:             time.Second * 300,
	}

	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("new clientset failed, err: %s", err.Error())
	}

	// TEST connect
	ctx := context.Background()
	listOption := v1.ListOptions{}

	nodeInterface := ClientSet.CoreV1().Nodes()
	_, err = nodeInterface.List(ctx, listOption)
	if err != nil {
		log.Fatalf("list node failed, err: %s", err.Error())
	}
	MetricClientSet, err = versioned.NewForConfig(config)
	if err != nil {
		log.Fatalf("new metric clientset failed, err: %s", err.Error())
	}
	LocalClientSet = ClientSet
}

func GetConfig() *rest.Config {
	return config
}
