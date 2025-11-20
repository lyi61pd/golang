package main

import (
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	dedicatedkmsopenapi "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi"
	dedicatedkmsopenapiutil "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/openapi-util"
	dedicatedkmssdk "github.com/aliyun/alibabacloud-dkms-gcs-go-sdk/sdk"
)

func main() {
	config := &dedicatedkmsopenapi.Config{
		Protocol:         tea.String("https"),
		ClientKeyContent: tea.String("<your client key content>"),
		Password:         tea.String("<your client key password>"),
		Endpoint:         tea.String("<your dkms instance service endpoint>"),
	}
	client, err := dedicatedkmssdk.NewClient(config)
	if err != nil {
		panic(err)
	}
	runtimeOptions := &dedicatedkmsopenapiutil.RuntimeOptions{
		IgnoreSSL: tea.Bool(true),
	}
	advanceEncryptRequest := &dedicatedkmssdk.AdvanceEncryptRequest{
		KeyId:     tea.String("<your key id>"),
		Plaintext: []byte("plaintext"),
	}
	advanceEncryptResponse, err := client.AdvanceEncryptWithOptions(advanceEncryptRequest, runtimeOptions)
	if err != nil {
		panic(err)
	}
	fmt.Println(advanceEncryptResponse)
}