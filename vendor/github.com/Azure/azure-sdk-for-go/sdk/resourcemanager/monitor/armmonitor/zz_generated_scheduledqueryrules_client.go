//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ScheduledQueryRulesClient contains the methods for the ScheduledQueryRules group.
// Don't use this type directly, use NewScheduledQueryRulesClient() instead.
type ScheduledQueryRulesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewScheduledQueryRulesClient creates a new instance of ScheduledQueryRulesClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewScheduledQueryRulesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ScheduledQueryRulesClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &ScheduledQueryRulesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates an log search rule.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-04-16
// resourceGroupName - The name of the resource group. The name is case insensitive.
// ruleName - The name of the rule.
// parameters - The parameters of the rule to create or update.
// options - ScheduledQueryRulesClientCreateOrUpdateOptions contains the optional parameters for the ScheduledQueryRulesClient.CreateOrUpdate
// method.
func (client *ScheduledQueryRulesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, ruleName string, parameters LogSearchRuleResource, options *ScheduledQueryRulesClientCreateOrUpdateOptions) (ScheduledQueryRulesClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, ruleName, parameters, options)
	if err != nil {
		return ScheduledQueryRulesClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ScheduledQueryRulesClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return ScheduledQueryRulesClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ScheduledQueryRulesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, ruleName string, parameters LogSearchRuleResource, options *ScheduledQueryRulesClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Insights/scheduledQueryRules/{ruleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ruleName == "" {
		return nil, errors.New("parameter ruleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ruleName}", url.PathEscape(ruleName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-04-16")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ScheduledQueryRulesClient) createOrUpdateHandleResponse(resp *http.Response) (ScheduledQueryRulesClientCreateOrUpdateResponse, error) {
	result := ScheduledQueryRulesClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LogSearchRuleResource); err != nil {
		return ScheduledQueryRulesClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a Log Search rule
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-04-16
// resourceGroupName - The name of the resource group. The name is case insensitive.
// ruleName - The name of the rule.
// options - ScheduledQueryRulesClientDeleteOptions contains the optional parameters for the ScheduledQueryRulesClient.Delete
// method.
func (client *ScheduledQueryRulesClient) Delete(ctx context.Context, resourceGroupName string, ruleName string, options *ScheduledQueryRulesClientDeleteOptions) (ScheduledQueryRulesClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, ruleName, options)
	if err != nil {
		return ScheduledQueryRulesClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ScheduledQueryRulesClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return ScheduledQueryRulesClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return ScheduledQueryRulesClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ScheduledQueryRulesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, ruleName string, options *ScheduledQueryRulesClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Insights/scheduledQueryRules/{ruleName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ruleName == "" {
		return nil, errors.New("parameter ruleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ruleName}", url.PathEscape(ruleName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-04-16")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets an Log Search rule
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-04-16
// resourceGroupName - The name of the resource group. The name is case insensitive.
// ruleName - The name of the rule.
// options - ScheduledQueryRulesClientGetOptions contains the optional parameters for the ScheduledQueryRulesClient.Get method.
func (client *ScheduledQueryRulesClient) Get(ctx context.Context, resourceGroupName string, ruleName string, options *ScheduledQueryRulesClientGetOptions) (ScheduledQueryRulesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, ruleName, options)
	if err != nil {
		return ScheduledQueryRulesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ScheduledQueryRulesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ScheduledQueryRulesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ScheduledQueryRulesClient) getCreateRequest(ctx context.Context, resourceGroupName string, ruleName string, options *ScheduledQueryRulesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Insights/scheduledQueryRules/{ruleName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ruleName == "" {
		return nil, errors.New("parameter ruleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ruleName}", url.PathEscape(ruleName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-04-16")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ScheduledQueryRulesClient) getHandleResponse(resp *http.Response) (ScheduledQueryRulesClientGetResponse, error) {
	result := ScheduledQueryRulesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LogSearchRuleResource); err != nil {
		return ScheduledQueryRulesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - List the Log Search rules within a resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-04-16
// resourceGroupName - The name of the resource group. The name is case insensitive.
// options - ScheduledQueryRulesClientListByResourceGroupOptions contains the optional parameters for the ScheduledQueryRulesClient.ListByResourceGroup
// method.
func (client *ScheduledQueryRulesClient) NewListByResourceGroupPager(resourceGroupName string, options *ScheduledQueryRulesClientListByResourceGroupOptions) *runtime.Pager[ScheduledQueryRulesClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[ScheduledQueryRulesClientListByResourceGroupResponse]{
		More: func(page ScheduledQueryRulesClientListByResourceGroupResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ScheduledQueryRulesClientListByResourceGroupResponse) (ScheduledQueryRulesClientListByResourceGroupResponse, error) {
			req, err := client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			if err != nil {
				return ScheduledQueryRulesClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ScheduledQueryRulesClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ScheduledQueryRulesClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ScheduledQueryRulesClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ScheduledQueryRulesClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Insights/scheduledQueryRules"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-04-16")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ScheduledQueryRulesClient) listByResourceGroupHandleResponse(resp *http.Response) (ScheduledQueryRulesClientListByResourceGroupResponse, error) {
	result := ScheduledQueryRulesClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LogSearchRuleResourceCollection); err != nil {
		return ScheduledQueryRulesClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListBySubscriptionPager - List the Log Search rules within a subscription group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-04-16
// options - ScheduledQueryRulesClientListBySubscriptionOptions contains the optional parameters for the ScheduledQueryRulesClient.ListBySubscription
// method.
func (client *ScheduledQueryRulesClient) NewListBySubscriptionPager(options *ScheduledQueryRulesClientListBySubscriptionOptions) *runtime.Pager[ScheduledQueryRulesClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[ScheduledQueryRulesClientListBySubscriptionResponse]{
		More: func(page ScheduledQueryRulesClientListBySubscriptionResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *ScheduledQueryRulesClientListBySubscriptionResponse) (ScheduledQueryRulesClientListBySubscriptionResponse, error) {
			req, err := client.listBySubscriptionCreateRequest(ctx, options)
			if err != nil {
				return ScheduledQueryRulesClientListBySubscriptionResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ScheduledQueryRulesClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ScheduledQueryRulesClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *ScheduledQueryRulesClient) listBySubscriptionCreateRequest(ctx context.Context, options *ScheduledQueryRulesClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Insights/scheduledQueryRules"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-04-16")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *ScheduledQueryRulesClient) listBySubscriptionHandleResponse(resp *http.Response) (ScheduledQueryRulesClientListBySubscriptionResponse, error) {
	result := ScheduledQueryRulesClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LogSearchRuleResourceCollection); err != nil {
		return ScheduledQueryRulesClientListBySubscriptionResponse{}, err
	}
	return result, nil
}

// Update - Update log search Rule.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2018-04-16
// resourceGroupName - The name of the resource group. The name is case insensitive.
// ruleName - The name of the rule.
// parameters - The parameters of the rule to update.
// options - ScheduledQueryRulesClientUpdateOptions contains the optional parameters for the ScheduledQueryRulesClient.Update
// method.
func (client *ScheduledQueryRulesClient) Update(ctx context.Context, resourceGroupName string, ruleName string, parameters LogSearchRuleResourcePatch, options *ScheduledQueryRulesClientUpdateOptions) (ScheduledQueryRulesClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, ruleName, parameters, options)
	if err != nil {
		return ScheduledQueryRulesClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ScheduledQueryRulesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ScheduledQueryRulesClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *ScheduledQueryRulesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, ruleName string, parameters LogSearchRuleResourcePatch, options *ScheduledQueryRulesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Insights/scheduledQueryRules/{ruleName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if ruleName == "" {
		return nil, errors.New("parameter ruleName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{ruleName}", url.PathEscape(ruleName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-04-16")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleResponse handles the Update response.
func (client *ScheduledQueryRulesClient) updateHandleResponse(resp *http.Response) (ScheduledQueryRulesClientUpdateResponse, error) {
	result := ScheduledQueryRulesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LogSearchRuleResource); err != nil {
		return ScheduledQueryRulesClientUpdateResponse{}, err
	}
	return result, nil
}