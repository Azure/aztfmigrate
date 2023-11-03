//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armpanngfw

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

// CertificateObjectGlobalRulestackClient contains the methods for the CertificateObjectGlobalRulestack group.
// Don't use this type directly, use NewCertificateObjectGlobalRulestackClient() instead.
type CertificateObjectGlobalRulestackClient struct {
	internal *arm.Client
}

// NewCertificateObjectGlobalRulestackClient creates a new instance of CertificateObjectGlobalRulestackClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCertificateObjectGlobalRulestackClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*CertificateObjectGlobalRulestackClient, error) {
	cl, err := arm.NewClient(moduleName+".CertificateObjectGlobalRulestackClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CertificateObjectGlobalRulestackClient{
		internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a CertificateObjectGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - globalRulestackName - GlobalRulestack resource name
//   - name - certificate name
//   - resource - Resource create parameters.
//   - options - CertificateObjectGlobalRulestackClientBeginCreateOrUpdateOptions contains the optional parameters for the CertificateObjectGlobalRulestackClient.BeginCreateOrUpdate
//     method.
func (client *CertificateObjectGlobalRulestackClient) BeginCreateOrUpdate(ctx context.Context, globalRulestackName string, name string, resource CertificateObjectGlobalRulestackResource, options *CertificateObjectGlobalRulestackClientBeginCreateOrUpdateOptions) (*runtime.Poller[CertificateObjectGlobalRulestackClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, globalRulestackName, name, resource, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CertificateObjectGlobalRulestackClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[CertificateObjectGlobalRulestackClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create a CertificateObjectGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
func (client *CertificateObjectGlobalRulestackClient) createOrUpdate(ctx context.Context, globalRulestackName string, name string, resource CertificateObjectGlobalRulestackResource, options *CertificateObjectGlobalRulestackClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, globalRulestackName, name, resource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *CertificateObjectGlobalRulestackClient) createOrUpdateCreateRequest(ctx context.Context, globalRulestackName string, name string, resource CertificateObjectGlobalRulestackResource, options *CertificateObjectGlobalRulestackClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/certificates/{name}"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, resource)
}

// BeginDelete - Delete a CertificateObjectGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - globalRulestackName - GlobalRulestack resource name
//   - name - certificate name
//   - options - CertificateObjectGlobalRulestackClientBeginDeleteOptions contains the optional parameters for the CertificateObjectGlobalRulestackClient.BeginDelete
//     method.
func (client *CertificateObjectGlobalRulestackClient) BeginDelete(ctx context.Context, globalRulestackName string, name string, options *CertificateObjectGlobalRulestackClientBeginDeleteOptions) (*runtime.Poller[CertificateObjectGlobalRulestackClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, globalRulestackName, name, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[CertificateObjectGlobalRulestackClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[CertificateObjectGlobalRulestackClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Delete a CertificateObjectGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
func (client *CertificateObjectGlobalRulestackClient) deleteOperation(ctx context.Context, globalRulestackName string, name string, options *CertificateObjectGlobalRulestackClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, globalRulestackName, name, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *CertificateObjectGlobalRulestackClient) deleteCreateRequest(ctx context.Context, globalRulestackName string, name string, options *CertificateObjectGlobalRulestackClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/certificates/{name}"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a CertificateObjectGlobalRulestackResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-29
//   - globalRulestackName - GlobalRulestack resource name
//   - name - certificate name
//   - options - CertificateObjectGlobalRulestackClientGetOptions contains the optional parameters for the CertificateObjectGlobalRulestackClient.Get
//     method.
func (client *CertificateObjectGlobalRulestackClient) Get(ctx context.Context, globalRulestackName string, name string, options *CertificateObjectGlobalRulestackClientGetOptions) (CertificateObjectGlobalRulestackClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, globalRulestackName, name, options)
	if err != nil {
		return CertificateObjectGlobalRulestackClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CertificateObjectGlobalRulestackClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return CertificateObjectGlobalRulestackClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *CertificateObjectGlobalRulestackClient) getCreateRequest(ctx context.Context, globalRulestackName string, name string, options *CertificateObjectGlobalRulestackClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/certificates/{name}"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
	if name == "" {
		return nil, errors.New("parameter name cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{name}", url.PathEscape(name))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CertificateObjectGlobalRulestackClient) getHandleResponse(resp *http.Response) (CertificateObjectGlobalRulestackClientGetResponse, error) {
	result := CertificateObjectGlobalRulestackClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CertificateObjectGlobalRulestackResource); err != nil {
		return CertificateObjectGlobalRulestackClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List CertificateObjectGlobalRulestackResource resources by Tenant
//
// Generated from API version 2022-08-29
//   - globalRulestackName - GlobalRulestack resource name
//   - options - CertificateObjectGlobalRulestackClientListOptions contains the optional parameters for the CertificateObjectGlobalRulestackClient.NewListPager
//     method.
func (client *CertificateObjectGlobalRulestackClient) NewListPager(globalRulestackName string, options *CertificateObjectGlobalRulestackClientListOptions) *runtime.Pager[CertificateObjectGlobalRulestackClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[CertificateObjectGlobalRulestackClientListResponse]{
		More: func(page CertificateObjectGlobalRulestackClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CertificateObjectGlobalRulestackClientListResponse) (CertificateObjectGlobalRulestackClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, globalRulestackName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CertificateObjectGlobalRulestackClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return CertificateObjectGlobalRulestackClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CertificateObjectGlobalRulestackClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *CertificateObjectGlobalRulestackClient) listCreateRequest(ctx context.Context, globalRulestackName string, options *CertificateObjectGlobalRulestackClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/PaloAltoNetworks.Cloudngfw/globalRulestacks/{globalRulestackName}/certificates"
	if globalRulestackName == "" {
		return nil, errors.New("parameter globalRulestackName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{globalRulestackName}", url.PathEscape(globalRulestackName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-29")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *CertificateObjectGlobalRulestackClient) listHandleResponse(resp *http.Response) (CertificateObjectGlobalRulestackClientListResponse, error) {
	result := CertificateObjectGlobalRulestackClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CertificateObjectGlobalRulestackResourceListResult); err != nil {
		return CertificateObjectGlobalRulestackClientListResponse{}, err
	}
	return result, nil
}