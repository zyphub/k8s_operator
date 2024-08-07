package k8s_client

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// 测试用的配置
const (
	// ClientCA ClientKey pid的机器
	ClientCA  = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURRakNDQWlxZ0F3SUJBZ0lJR0ttaW4wOFQ0c0V3RFFZSktvWklodmNOQVFFTEJRQXdGVEVUTUJFR0ExVUUKQXhNS2EzVmlaWEp1WlhSbGN6QWVGdzB5TkRBM01EUXdPREl3TURWYUZ3MHlOVEE0TURjeE5EVXlNakJhTURZeApGekFWQmdOVkJBb1REbk41YzNSbGJUcHRZWE4wWlhKek1Sc3dHUVlEVlFRREV4SmtiMk5yWlhJdFptOXlMV1JsCmMydDBiM0F3Z2dFaU1BMEdDU3FHU0liM0RRRUJBUVVBQTRJQkR3QXdnZ0VLQW9JQkFRRHV3OVZYZzBHRWl6NDcKaDZlZlZHcUV5aWdVS1RFTTlkUDArM3dmL3pyT2Zvd3BiNTkxbG1aVDZETzVWUDJ2alJ5U2l1R29la0JKMDdwTwpLRTdqeEM1Q0QxUGJtSkw1WGVtM1FuYjFWcldzRGtrL2FaL0pDV2lZNW9EUWhnOHRiZDhiY3Yyd3BlRld0YzF2CjVMV2k5R1QySzZRZmdQNTNPVkZGdVBNS2JXRDJIbVltem9OdUo0OUxQbno2MlFlb1RXQ2E0QS9MK0dwWDVpM0kKVWVLNkFMY1VGNVIvVE5xZEg5VnpybTA4YlMybmhBM085dGJ4Q1NCdzdwSlNLTHNocFFUOFZJRGR2Nzl5OE43QwpjcFVKZXNnVEJaVFArZGhhLzVpV1R3bFhLN1BMZEUwZ04rUUwvbmExQ3ROUWZGTUtOeUlBRFJ1M21NaEs3anpPCjI0cUpkU1JQQWdNQkFBR2pkVEJ6TUE0R0ExVWREd0VCL3dRRUF3SUZvREFUQmdOVkhTVUVEREFLQmdnckJnRUYKQlFjREFqQU1CZ05WSFJNQkFmOEVBakFBTUI4R0ExVWRJd1FZTUJhQUZMcEwralVnN2p1ZlB0dmpvVWNndE8ySwoxSE9YTUIwR0ExVWRFUVFXTUJTQ0VtUnZZMnRsY2kxbWIzSXRaR1Z6YTNSdmNEQU5CZ2txaGtpRzl3MEJBUXNGCkFBT0NBUUVBT0tVTzZ3anhtc203a0M4OCtpRVFJTUZoVFhFeDR2UnZhYjd3YU1ZSGhMUXcxYkpFejNyaTJ4OTAKcHZqWGxuc282WFBjTWtDZDJQTXJCTkhyU2dYQmhMZ1pjUGROYmNTY3V0bkRGSzNPYjdiZWVLNkZuY2FJcU5pVAp1eFpaMFdhS0tKQVBOWVNYaldQbGVETEFoY1laeDlSY2RSNFFEZHFNZUJxb1pCaVN5UGZRVlNBT2Z6Zm5DV0ZXCmFpUUo2Y0d4akxRaithOTV4bVBrNE1CUmRuNlJOM3NyTW0rZDNJZTlTc3dKUzdySXZia0ZWU0xrbm9aRmFSbisKOGd0bVlvQUo1OGRtOXVRNUxQZGpzQVo1K0ZhM0lxdWd4VmwrSEZ2a0RwQ3d0cnZBZmVBY29DNEh5UWpJQjRBVQpTMjI1aEFYazdNMExEMFNEQTl6NWJKa01iVnRXdWc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
	ClientKey = "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBN3NQVlY0TkJoSXMrTzRlbm4xUnFoTW9vRkNreERQWFQ5UHQ4SC84NnpuNk1LVytmCmRaWm1VK2d6dVZUOXI0MGNrb3JocUhwQVNkTzZUaWhPNDhRdVFnOVQyNWlTK1YzcHQwSjI5VmExckE1SlAybWYKeVFsb21PYUEwSVlQTFczZkczTDlzS1hoVnJYTmIrUzFvdlJrOWl1a0g0RCtkemxSUmJqekNtMWc5aDVtSnM2RApiaWVQU3o1OCt0a0hxRTFnbXVBUHkvaHFWK1l0eUZIaXVnQzNGQmVVZjB6YW5SL1ZjNjV0UEcwdHA0UU56dmJXCjhRa2djTzZTVWlpN0lhVUUvRlNBM2IrL2N2RGV3bktWQ1hySUV3V1V6L25ZV3YrWWxrOEpWeXV6eTNSTklEZmsKQy81MnRRclRVSHhUQ2pjaUFBMGJ0NWpJU3U0OHp0dUtpWFVrVHdJREFRQUJBb0lCQURaNzg3UUxuS2pOU1g4MgpIbmNLUVdCWjdUbGtpTy9uTE4zcmdWQ2Y0bUI2bWl0ZWNHblp6ekg0ZTgwZjZ0L2plSkNzSm9CV25WTDdnTGtUCkU2Vi8vL3BOR3hxeVAxK3VJWVlUSWFnc2lEcGg3QzhQUUVvVTNveDlsUW1BZmZnazZWT1BNdnJiYjRkazV5TlEKY201a0RLSHNKWUNXNC9wNjF1UHRKM0RLc3VTV2JHMU5wOG4zdXNHaWdLY2FvU1lnc2tKSEZqeWZIcktvdnBqdAorc3VaZ2pJWVJWc0p6cFhSQjVJVFlKR0xJZVlpTjMxRklPUFBmVEYyMlRRT21iRlBwUFoxRGM1THp0UElXU2RhCndkcEhJUTZsaWlnajNiSlArOFRTdkxxck1JemFFRUJlZ29Zb2NTbGpoSS9pREpuSjN1a1IyaU5EcmV0N3c3cTMKUDFkQUVnRUNnWUVBOW5mYzkvZTFvREdQN29raHhjNThWcTZwTmYveCt3UVZQUmExVTNoc1RwOVJoUHpZZ1MrcApHY0l2d3l0QmtqYmUrc09uYlpMeVBKWGRsZkFZV2IzMFQ0RUdWYUF5Q0xFN29QSWdNRHdxM2cvZ1ArRVlCU0wrCjF5bGVBeTJQK01MSGc1U0ZYM0g5bjhFdCtDTFZJU1k4dmdQYnBIRldHZm5JVFJwSURucVRyUUVDZ1lFQTkvKzAKZDFJN1IrY0NFZ0MveTNkK1NKditTTFdwczNwZ0lNaXRRUjZlbU1VZUlBWHZhczhQYTFic21qTTB3Znc1SW5XTgpHaW5iR01TeFFCN2xiNE5uZVN4U2tGRDlOL1FROFp3Q3ZYU2tJUWlQNExDc3d4TExPZU5zZktXQ2k5YTdMb2pjCkhka0lUREppR29yWnF1cnY0cUFSYno1OTlSMUUxS3IweUU5MXdVOENnWUFOMWZhNm1OWkNVdVh3anhRdFJZVW4KWEpDMUxsUUlNbGQ2NFc1MmJCa3daTE12MHYzWWFyT0VkYWsydkpQbXdGdk9HZk9wTEFtYkt4S1FXelVTdko0ZApaSEhWbHJPWVYxS3dtMGNCVGk5ZDNlaEp6Ym9LZDhkMGpxYnZhTHhmUzVmbHBBM0VxT0tDK0ZZN1NzRktKaHBjCjFGeWRJNXVndzZ2aDRDclJYVUl6QVFLQmdBeitud3dwaU9XcG14Z3FaZUpaYm9xTGNmV0pYMDBDT29zOU9LYlMKM2VpUFc1YTkrTitWM2U3MzdRbmZhUUpKSHcxSkw0MlJaK09TV3Q5TFB5WnFzajlOTFQ1V29BNFFnZHJIRy9XbAphUHc2SUovYlloSU9xQXR1ZVQ3R3hXSmliQWh5TDJaNCt0QlRTNFNzaGQ5STFDMEJ5aWdVRkRHRnlSZURwYlBoClJnQk5Bb0dCQUo1bFdLbnVFVWVtcW1xTmMvT2pKa29iaXI5WW9vTndXYjNhZGZZVU92OVVFRFJmeFhHQ0F6dWYKQ3dDcVlsMFFTcEd1ZThGNHg3WENGY0tPa1BVU3Rvd0xyTk92YTZMaXR5N1dRK01NSjJ0M042RGVBMkdBek1sdgpXeUNPbTdBdFd3ZEthZHk5NEhnN2tuL1JlVWtPa3VGQVFTdTJwcnFDSGFycG1hOHowZERTCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg=="
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

func GetClientSet() *kubernetes.Clientset {
	return ClientSet
}

func RestMapper() (meta.RESTMapper, error) {
	// 动态获取所有的api group resources
	apiGroupResources, err := restmapper.GetAPIGroupResources(ClientSet.Discovery())
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDiscoveryRESTMapper(apiGroupResources)
	return mapper, nil
}

func InitWatch() (informers.SharedInformerFactory, error) {
	factory := informers.NewSharedInformerFactory(ClientSet, 0)

	handler, err := factory.Core().V1().
		Namespaces().
		Informer().
		AddEventHandler(cache.ResourceEventHandlerFuncs{})
	if err != nil {
		klog.Errorf("add namespace informer event handler failed, err: %s", err.Error())
		return nil, err
	}

	_ = handler

	return factory, nil
}
