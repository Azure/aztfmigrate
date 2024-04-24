package client

import (
	"context"
	"net/http"

	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type RawClient struct {
	host string
	pl   runtime.Pipeline
}

func (b *ClientBuilder) NewRawClient() (*RawClient, error) {
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := b.ClientOpt.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline("resource", "v0.1.0", b.Cred, runtime.PipelineOptions{}, &b.ClientOpt)
	if err != nil {
		return nil, err
	}
	client := &RawClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

func (client *RawClient) Get(ctx context.Context, resourceID string, apiVersion string) (interface{}, error) {
	req, err := client.getCreateRequest(ctx, resourceID, apiVersion)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return nil, runtime.NewResponseError(resp)
	}

	var responseBody interface{}
	if err := runtime.UnmarshalAsJSON(resp, &responseBody); err != nil {
		return nil, err
	}
	return responseBody, nil
}

func (client *RawClient) getCreateRequest(ctx context.Context, resourceID string, apiVersion string) (*policy.Request, error) {
	urlPath := resourceID
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", apiVersion)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}
