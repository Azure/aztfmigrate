//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsql

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ElasticPoolsClient contains the methods for the ElasticPools group.
// Don't use this type directly, use NewElasticPoolsClient() instead.
type ElasticPoolsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewElasticPoolsClient creates a new instance of ElasticPoolsClient with the specified values.
//   - subscriptionID - The subscription ID that identifies an Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewElasticPoolsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ElasticPoolsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ElasticPoolsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - elasticPoolName - The name of the elastic pool.
//   - parameters - The elastic pool parameters.
//   - options - ElasticPoolsClientBeginCreateOrUpdateOptions contains the optional parameters for the ElasticPoolsClient.BeginCreateOrUpdate
//     method.
func (client *ElasticPoolsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, parameters ElasticPool, options *ElasticPoolsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ElasticPoolsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, serverName, elasticPoolName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ElasticPoolsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ElasticPoolsClientCreateOrUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// CreateOrUpdate - Creates or updates an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
func (client *ElasticPoolsClient) createOrUpdate(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, parameters ElasticPool, options *ElasticPoolsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ElasticPoolsClient.BeginCreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serverName, elasticPoolName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ElasticPoolsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, parameters ElasticPool, options *ElasticPoolsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if elasticPoolName == "" {
		return nil, errors.New("parameter elasticPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{elasticPoolName}", url.PathEscape(elasticPoolName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - elasticPoolName - The name of the elastic pool.
//   - options - ElasticPoolsClientBeginDeleteOptions contains the optional parameters for the ElasticPoolsClient.BeginDelete
//     method.
func (client *ElasticPoolsClient) BeginDelete(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientBeginDeleteOptions) (*runtime.Poller[ElasticPoolsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, serverName, elasticPoolName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ElasticPoolsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ElasticPoolsClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
func (client *ElasticPoolsClient) deleteOperation(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "ElasticPoolsClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serverName, elasticPoolName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ElasticPoolsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if elasticPoolName == "" {
		return nil, errors.New("parameter elasticPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{elasticPoolName}", url.PathEscape(elasticPoolName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// BeginFailover - Failovers an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - elasticPoolName - The name of the elastic pool to failover.
//   - options - ElasticPoolsClientBeginFailoverOptions contains the optional parameters for the ElasticPoolsClient.BeginFailover
//     method.
func (client *ElasticPoolsClient) BeginFailover(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientBeginFailoverOptions) (*runtime.Poller[ElasticPoolsClientFailoverResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.failover(ctx, resourceGroupName, serverName, elasticPoolName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ElasticPoolsClientFailoverResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ElasticPoolsClientFailoverResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Failover - Failovers an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
func (client *ElasticPoolsClient) failover(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientBeginFailoverOptions) (*http.Response, error) {
	var err error
	const operationName = "ElasticPoolsClient.BeginFailover"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.failoverCreateRequest(ctx, resourceGroupName, serverName, elasticPoolName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// failoverCreateRequest creates the Failover request.
func (client *ElasticPoolsClient) failoverCreateRequest(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientBeginFailoverOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}/failover"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if elasticPoolName == "" {
		return nil, errors.New("parameter elasticPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{elasticPoolName}", url.PathEscape(elasticPoolName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Gets an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - elasticPoolName - The name of the elastic pool.
//   - options - ElasticPoolsClientGetOptions contains the optional parameters for the ElasticPoolsClient.Get method.
func (client *ElasticPoolsClient) Get(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientGetOptions) (ElasticPoolsClientGetResponse, error) {
	var err error
	const operationName = "ElasticPoolsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, serverName, elasticPoolName, options)
	if err != nil {
		return ElasticPoolsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ElasticPoolsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ElasticPoolsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ElasticPoolsClient) getCreateRequest(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if elasticPoolName == "" {
		return nil, errors.New("parameter elasticPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{elasticPoolName}", url.PathEscape(elasticPoolName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ElasticPoolsClient) getHandleResponse(resp *http.Response) (ElasticPoolsClientGetResponse, error) {
	result := ElasticPoolsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ElasticPool); err != nil {
		return ElasticPoolsClientGetResponse{}, err
	}
	return result, nil
}

// NewListByServerPager - Gets all elastic pools in a server.
//
// Generated from API version 2021-08-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - options - ElasticPoolsClientListByServerOptions contains the optional parameters for the ElasticPoolsClient.NewListByServerPager
//     method.
func (client *ElasticPoolsClient) NewListByServerPager(resourceGroupName string, serverName string, options *ElasticPoolsClientListByServerOptions) *runtime.Pager[ElasticPoolsClientListByServerResponse] {
	return runtime.NewPager(runtime.PagingHandler[ElasticPoolsClientListByServerResponse]{
		More: func(page ElasticPoolsClientListByServerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ElasticPoolsClientListByServerResponse) (ElasticPoolsClientListByServerResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ElasticPoolsClient.NewListByServerPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByServerCreateRequest(ctx, resourceGroupName, serverName, options)
			}, nil)
			if err != nil {
				return ElasticPoolsClientListByServerResponse{}, err
			}
			return client.listByServerHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByServerCreateRequest creates the ListByServer request.
func (client *ElasticPoolsClient) listByServerCreateRequest(ctx context.Context, resourceGroupName string, serverName string, options *ElasticPoolsClientListByServerOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(*options.Skip, 10))
	}
	reqQP.Set("api-version", "2021-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByServerHandleResponse handles the ListByServer response.
func (client *ElasticPoolsClient) listByServerHandleResponse(resp *http.Response) (ElasticPoolsClientListByServerResponse, error) {
	result := ElasticPoolsClientListByServerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ElasticPoolListResult); err != nil {
		return ElasticPoolsClientListByServerResponse{}, err
	}
	return result, nil
}

// NewListMetricDefinitionsPager - Returns elastic pool metric definitions.
//
// Generated from API version 2014-04-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - elasticPoolName - The name of the elastic pool.
//   - options - ElasticPoolsClientListMetricDefinitionsOptions contains the optional parameters for the ElasticPoolsClient.NewListMetricDefinitionsPager
//     method.
func (client *ElasticPoolsClient) NewListMetricDefinitionsPager(resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientListMetricDefinitionsOptions) *runtime.Pager[ElasticPoolsClientListMetricDefinitionsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ElasticPoolsClientListMetricDefinitionsResponse]{
		More: func(page ElasticPoolsClientListMetricDefinitionsResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ElasticPoolsClientListMetricDefinitionsResponse) (ElasticPoolsClientListMetricDefinitionsResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ElasticPoolsClient.NewListMetricDefinitionsPager")
			req, err := client.listMetricDefinitionsCreateRequest(ctx, resourceGroupName, serverName, elasticPoolName, options)
			if err != nil {
				return ElasticPoolsClientListMetricDefinitionsResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ElasticPoolsClientListMetricDefinitionsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ElasticPoolsClientListMetricDefinitionsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listMetricDefinitionsHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listMetricDefinitionsCreateRequest creates the ListMetricDefinitions request.
func (client *ElasticPoolsClient) listMetricDefinitionsCreateRequest(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, options *ElasticPoolsClientListMetricDefinitionsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}/metricDefinitions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if elasticPoolName == "" {
		return nil, errors.New("parameter elasticPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{elasticPoolName}", url.PathEscape(elasticPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2014-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listMetricDefinitionsHandleResponse handles the ListMetricDefinitions response.
func (client *ElasticPoolsClient) listMetricDefinitionsHandleResponse(resp *http.Response) (ElasticPoolsClientListMetricDefinitionsResponse, error) {
	result := ElasticPoolsClientListMetricDefinitionsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MetricDefinitionListResult); err != nil {
		return ElasticPoolsClientListMetricDefinitionsResponse{}, err
	}
	return result, nil
}

// NewListMetricsPager - Returns elastic pool metrics.
//
// Generated from API version 2014-04-01
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - elasticPoolName - The name of the elastic pool.
//   - filter - An OData filter expression that describes a subset of metrics to return.
//   - options - ElasticPoolsClientListMetricsOptions contains the optional parameters for the ElasticPoolsClient.NewListMetricsPager
//     method.
func (client *ElasticPoolsClient) NewListMetricsPager(resourceGroupName string, serverName string, elasticPoolName string, filter string, options *ElasticPoolsClientListMetricsOptions) *runtime.Pager[ElasticPoolsClientListMetricsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ElasticPoolsClientListMetricsResponse]{
		More: func(page ElasticPoolsClientListMetricsResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ElasticPoolsClientListMetricsResponse) (ElasticPoolsClientListMetricsResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ElasticPoolsClient.NewListMetricsPager")
			req, err := client.listMetricsCreateRequest(ctx, resourceGroupName, serverName, elasticPoolName, filter, options)
			if err != nil {
				return ElasticPoolsClientListMetricsResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ElasticPoolsClientListMetricsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ElasticPoolsClientListMetricsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listMetricsHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listMetricsCreateRequest creates the ListMetrics request.
func (client *ElasticPoolsClient) listMetricsCreateRequest(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, filter string, options *ElasticPoolsClientListMetricsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}/metrics"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if elasticPoolName == "" {
		return nil, errors.New("parameter elasticPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{elasticPoolName}", url.PathEscape(elasticPoolName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2014-04-01")
	reqQP.Set("$filter", filter)
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listMetricsHandleResponse handles the ListMetrics response.
func (client *ElasticPoolsClient) listMetricsHandleResponse(resp *http.Response) (ElasticPoolsClientListMetricsResponse, error) {
	result := ElasticPoolsClientListMetricsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.MetricListResult); err != nil {
		return ElasticPoolsClientListMetricsResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Updates an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - elasticPoolName - The name of the elastic pool.
//   - parameters - The elastic pool update parameters.
//   - options - ElasticPoolsClientBeginUpdateOptions contains the optional parameters for the ElasticPoolsClient.BeginUpdate
//     method.
func (client *ElasticPoolsClient) BeginUpdate(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, parameters ElasticPoolUpdate, options *ElasticPoolsClientBeginUpdateOptions) (*runtime.Poller[ElasticPoolsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, serverName, elasticPoolName, parameters, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ElasticPoolsClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ElasticPoolsClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Updates an elastic pool.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-08-01-preview
func (client *ElasticPoolsClient) update(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, parameters ElasticPoolUpdate, options *ElasticPoolsClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "ElasticPoolsClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serverName, elasticPoolName, parameters, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *ElasticPoolsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serverName string, elasticPoolName string, parameters ElasticPoolUpdate, options *ElasticPoolsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/elasticPools/{elasticPoolName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if elasticPoolName == "" {
		return nil, errors.New("parameter elasticPoolName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{elasticPoolName}", url.PathEscape(elasticPoolName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}
