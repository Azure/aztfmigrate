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
	"strings"
)

// ReplicationLinksClient contains the methods for the ReplicationLinks group.
// Don't use this type directly, use NewReplicationLinksClient() instead.
type ReplicationLinksClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewReplicationLinksClient creates a new instance of ReplicationLinksClient with the specified values.
//   - subscriptionID - The subscription ID that identifies an Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewReplicationLinksClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ReplicationLinksClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ReplicationLinksClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Delete - Deletes the replication link.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - databaseName - The name of the database.
//   - options - ReplicationLinksClientDeleteOptions contains the optional parameters for the ReplicationLinksClient.Delete method.
func (client *ReplicationLinksClient) Delete(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientDeleteOptions) (ReplicationLinksClientDeleteResponse, error) {
	var err error
	const operationName = "ReplicationLinksClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serverName, databaseName, linkID, options)
	if err != nil {
		return ReplicationLinksClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ReplicationLinksClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ReplicationLinksClientDeleteResponse{}, err
	}
	return ReplicationLinksClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ReplicationLinksClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/replicationLinks/{linkId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if linkID == "" {
		return nil, errors.New("parameter linkID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkId}", url.PathEscape(linkID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// BeginFailover - Fails over from the current primary server to this server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - databaseName - The name of the database.
//   - linkID - The name of the replication link.
//   - options - ReplicationLinksClientBeginFailoverOptions contains the optional parameters for the ReplicationLinksClient.BeginFailover
//     method.
func (client *ReplicationLinksClient) BeginFailover(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientBeginFailoverOptions) (*runtime.Poller[ReplicationLinksClientFailoverResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.failover(ctx, resourceGroupName, serverName, databaseName, linkID, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReplicationLinksClientFailoverResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ReplicationLinksClientFailoverResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Failover - Fails over from the current primary server to this server.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
func (client *ReplicationLinksClient) failover(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientBeginFailoverOptions) (*http.Response, error) {
	var err error
	const operationName = "ReplicationLinksClient.BeginFailover"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.failoverCreateRequest(ctx, resourceGroupName, serverName, databaseName, linkID, options)
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
func (client *ReplicationLinksClient) failoverCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientBeginFailoverOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/replicationLinks/{linkId}/failover"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if linkID == "" {
		return nil, errors.New("parameter linkID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkId}", url.PathEscape(linkID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginFailoverAllowDataLoss - Fails over from the current primary server to this server allowing data loss.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - databaseName - The name of the database.
//   - linkID - The name of the replication link.
//   - options - ReplicationLinksClientBeginFailoverAllowDataLossOptions contains the optional parameters for the ReplicationLinksClient.BeginFailoverAllowDataLoss
//     method.
func (client *ReplicationLinksClient) BeginFailoverAllowDataLoss(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientBeginFailoverAllowDataLossOptions) (*runtime.Poller[ReplicationLinksClientFailoverAllowDataLossResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.failoverAllowDataLoss(ctx, resourceGroupName, serverName, databaseName, linkID, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReplicationLinksClientFailoverAllowDataLossResponse]{
			Tracer: client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ReplicationLinksClientFailoverAllowDataLossResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// FailoverAllowDataLoss - Fails over from the current primary server to this server allowing data loss.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
func (client *ReplicationLinksClient) failoverAllowDataLoss(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientBeginFailoverAllowDataLossOptions) (*http.Response, error) {
	var err error
	const operationName = "ReplicationLinksClient.BeginFailoverAllowDataLoss"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.failoverAllowDataLossCreateRequest(ctx, resourceGroupName, serverName, databaseName, linkID, options)
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

// failoverAllowDataLossCreateRequest creates the FailoverAllowDataLoss request.
func (client *ReplicationLinksClient) failoverAllowDataLossCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientBeginFailoverAllowDataLossOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/replicationLinks/{linkId}/forceFailoverAllowDataLoss"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if linkID == "" {
		return nil, errors.New("parameter linkID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkId}", url.PathEscape(linkID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets a replication link.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - databaseName - The name of the database.
//   - linkID - The name of the replication link.
//   - options - ReplicationLinksClientGetOptions contains the optional parameters for the ReplicationLinksClient.Get method.
func (client *ReplicationLinksClient) Get(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientGetOptions) (ReplicationLinksClientGetResponse, error) {
	var err error
	const operationName = "ReplicationLinksClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, serverName, databaseName, linkID, options)
	if err != nil {
		return ReplicationLinksClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ReplicationLinksClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ReplicationLinksClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ReplicationLinksClient) getCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, linkID string, options *ReplicationLinksClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/replicationLinks/{linkId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if linkID == "" {
		return nil, errors.New("parameter linkID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkId}", url.PathEscape(linkID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ReplicationLinksClient) getHandleResponse(resp *http.Response) (ReplicationLinksClientGetResponse, error) {
	result := ReplicationLinksClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReplicationLink); err != nil {
		return ReplicationLinksClientGetResponse{}, err
	}
	return result, nil
}

// NewListByDatabasePager - Gets a list of replication links on database.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - databaseName - The name of the database.
//   - options - ReplicationLinksClientListByDatabaseOptions contains the optional parameters for the ReplicationLinksClient.NewListByDatabasePager
//     method.
func (client *ReplicationLinksClient) NewListByDatabasePager(resourceGroupName string, serverName string, databaseName string, options *ReplicationLinksClientListByDatabaseOptions) *runtime.Pager[ReplicationLinksClientListByDatabaseResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReplicationLinksClientListByDatabaseResponse]{
		More: func(page ReplicationLinksClientListByDatabaseResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReplicationLinksClientListByDatabaseResponse) (ReplicationLinksClientListByDatabaseResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ReplicationLinksClient.NewListByDatabasePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByDatabaseCreateRequest(ctx, resourceGroupName, serverName, databaseName, options)
			}, nil)
			if err != nil {
				return ReplicationLinksClientListByDatabaseResponse{}, err
			}
			return client.listByDatabaseHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByDatabaseCreateRequest creates the ListByDatabase request.
func (client *ReplicationLinksClient) listByDatabaseCreateRequest(ctx context.Context, resourceGroupName string, serverName string, databaseName string, options *ReplicationLinksClientListByDatabaseOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/databases/{databaseName}/replicationLinks"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverName == "" {
		return nil, errors.New("parameter serverName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverName}", url.PathEscape(serverName))
	if databaseName == "" {
		return nil, errors.New("parameter databaseName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseName}", url.PathEscape(databaseName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByDatabaseHandleResponse handles the ListByDatabase response.
func (client *ReplicationLinksClient) listByDatabaseHandleResponse(resp *http.Response) (ReplicationLinksClientListByDatabaseResponse, error) {
	result := ReplicationLinksClientListByDatabaseResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReplicationLinkListResult); err != nil {
		return ReplicationLinksClientListByDatabaseResponse{}, err
	}
	return result, nil
}

// NewListByServerPager - Gets a list of replication links.
//
// Generated from API version 2021-11-01-preview
//   - resourceGroupName - The name of the resource group that contains the resource. You can obtain this value from the Azure
//     Resource Manager API or the portal.
//   - serverName - The name of the server.
//   - options - ReplicationLinksClientListByServerOptions contains the optional parameters for the ReplicationLinksClient.NewListByServerPager
//     method.
func (client *ReplicationLinksClient) NewListByServerPager(resourceGroupName string, serverName string, options *ReplicationLinksClientListByServerOptions) *runtime.Pager[ReplicationLinksClientListByServerResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReplicationLinksClientListByServerResponse]{
		More: func(page ReplicationLinksClientListByServerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReplicationLinksClientListByServerResponse) (ReplicationLinksClientListByServerResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ReplicationLinksClient.NewListByServerPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByServerCreateRequest(ctx, resourceGroupName, serverName, options)
			}, nil)
			if err != nil {
				return ReplicationLinksClientListByServerResponse{}, err
			}
			return client.listByServerHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByServerCreateRequest creates the ListByServer request.
func (client *ReplicationLinksClient) listByServerCreateRequest(ctx context.Context, resourceGroupName string, serverName string, options *ReplicationLinksClientListByServerOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Sql/servers/{serverName}/replicationLinks"
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
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByServerHandleResponse handles the ListByServer response.
func (client *ReplicationLinksClient) listByServerHandleResponse(resp *http.Response) (ReplicationLinksClientListByServerResponse, error) {
	result := ReplicationLinksClientListByServerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReplicationLinkListResult); err != nil {
		return ReplicationLinksClientListByServerResponse{}, err
	}
	return result, nil
}
