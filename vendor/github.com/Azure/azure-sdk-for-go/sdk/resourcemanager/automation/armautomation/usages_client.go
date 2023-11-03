//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armautomation

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

// UsagesClient contains the methods for the Usages group.
// Don't use this type directly, use NewUsagesClient() instead.
type UsagesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewUsagesClient creates a new instance of UsagesClient with the specified values.
//   - subscriptionID - Gets subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID
//     forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewUsagesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*UsagesClient, error) {
	cl, err := arm.NewClient(moduleName+".UsagesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &UsagesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// NewListByAutomationAccountPager - Retrieve the usage for the account id.
//
// Generated from API version 2021-06-22
//   - resourceGroupName - Name of an Azure Resource group.
//   - automationAccountName - The name of the automation account.
//   - options - UsagesClientListByAutomationAccountOptions contains the optional parameters for the UsagesClient.NewListByAutomationAccountPager
//     method.
func (client *UsagesClient) NewListByAutomationAccountPager(resourceGroupName string, automationAccountName string, options *UsagesClientListByAutomationAccountOptions) *runtime.Pager[UsagesClientListByAutomationAccountResponse] {
	return runtime.NewPager(runtime.PagingHandler[UsagesClientListByAutomationAccountResponse]{
		More: func(page UsagesClientListByAutomationAccountResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *UsagesClientListByAutomationAccountResponse) (UsagesClientListByAutomationAccountResponse, error) {
			req, err := client.listByAutomationAccountCreateRequest(ctx, resourceGroupName, automationAccountName, options)
			if err != nil {
				return UsagesClientListByAutomationAccountResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return UsagesClientListByAutomationAccountResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return UsagesClientListByAutomationAccountResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByAutomationAccountHandleResponse(resp)
		},
	})
}

// listByAutomationAccountCreateRequest creates the ListByAutomationAccount request.
func (client *UsagesClient) listByAutomationAccountCreateRequest(ctx context.Context, resourceGroupName string, automationAccountName string, options *UsagesClientListByAutomationAccountOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/usages"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationAccountName == "" {
		return nil, errors.New("parameter automationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationAccountName}", url.PathEscape(automationAccountName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-06-22")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByAutomationAccountHandleResponse handles the ListByAutomationAccount response.
func (client *UsagesClient) listByAutomationAccountHandleResponse(resp *http.Response) (UsagesClientListByAutomationAccountResponse, error) {
	result := UsagesClientListByAutomationAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UsageListResult); err != nil {
		return UsagesClientListByAutomationAccountResponse{}, err
	}
	return result, nil
}