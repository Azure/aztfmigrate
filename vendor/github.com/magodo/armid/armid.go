package armid

import (
	"fmt"
	"strings"
)

type ResourceId interface {
	// RootScope returns the root scope of this resource.
	// For scoped resource, it is the including root scope resource id.
	// For root scopes, it returns itself.
	RootScope() RootScope

	// ParentScope returns the parent scope of this resource. Normally, scopes are seperated by "/providers/".
	// This is nil if the resource itself is a root scope.
	// E.g.
	// - /subscriptions/0000/resourceGroups/rg1/providers/Microsoft.Foo/foos/foo1 	-(parent scope)-> /subscriptions/0000/resourceGroups/rg1
	// - /subscriptions/0000/resourceGroups/rg1 									-(parent scope)-> nil
	ParentScope() ResourceId

	// Parent returns the parent resource. The parent resource belongs to the same provider as the current resource.
	// Nil is returned if the current resource is a root scoped resource, or this is a root scope.
	Parent() ResourceId

	// Provider returns the provider namespace of this resource id.
	// For scoped resource, it is the provider namespace of its routing scope, i.e. the scope of the resource itself.
	// For root scopes, it is a builtin provider namespace, e.g. "Microsoft.Resources".
	Provider() string

	// Types returns the resource type array of this resource.
	// For scoped resource, it is the sub-types of its routing scope, i.e. the scope of the resource itself.
	// e.g. ["virtualNetworks", "subnets"] for "Microsoft.Network/virtualNetworks/subnets"
	// For root scopes, it is a builtin type.
	Types() []string

	// Names returns the resource name array of this resource.
	// For scoped resource, it is the names of each sub-type of the Types(), which indicates it always has the same length as the return value of Types().
	// For root scopes, it is nil.
	Names() []string

	// TypeString returns the resource type string literal.
	// For scoped resource, this is the same as their RouteScopeString, with the leading "/" trimmed.
	// For root scope resource, this is the same as "Microsoft.Resources" + id.ScopeString()
	TypeString() string

	// String returns the resource id literal.
	String() string

	// Equal checkes the equality of two resource id. They are regarded as equal only when all the components are equal (case insensitively).
	Equal(ResourceId) bool

	// ScopeEqual checkes the equality of two resource id without taking the Names() into consideration.
	ScopeEqual(ResourceId) bool

	// ScopeString returns the string literal of the resource id without taking the Names() into consideration, for the whole resource id.
	ScopeString() string

	// RouteScopeString is similar as ScopeString, but only for the router scope (i.e. the last scope).
	// For root scope resource, this is the same as ScopeString.
	RouteScopeString() string

	// Normalize normalizes the invariant parts (e.g. Provider, Types) of the id  based on the input scope string.
	// The input scope string must be the same as calling the `ScopeString` of this id, except the casing.
	Normalize(string) error

	// Clone deep clones a ResourceId
	Clone() ResourceId
}

func ParseResourceId(id string) (ResourceId, error) {
	if id == "/" {
		return &TenantId{}, nil
	}
	if !strings.HasPrefix(id, "/") {
		return nil, fmt.Errorf(`id should start with "/"`)
	}
	segs := strings.Split(id[1:], "/")

	for idx, seg := range segs {
		if seg == "" {
			return nil, fmt.Errorf(`empty segment found behind %dth "/"`, idx+1)
		}
	}

	var rootScope RootScope = &TenantId{}
	if len(segs) >= 4 && strings.EqualFold(segs[0], "subscriptions") && strings.EqualFold(segs[2], "resourcegroups") {
		rootScope = &ResourceGroup{
			SubscriptionId: segs[1],
			Name:           segs[3],
		}
		segs = segs[4:]
	} else if len(segs) >= 2 && strings.EqualFold(segs[0], "subscriptions") {
		rootScope = &SubscriptionId{
			Id: segs[1],
		}
		segs = segs[2:]
	} else if len(segs) >= 4 && strings.EqualFold(segs[0], "providers") && strings.EqualFold(segs[1], "Microsoft.Management") && strings.EqualFold(segs[2], "managementgroups") {
		rootScope = &ManagementGroup{
			Name: segs[3],
		}
		segs = segs[4:]
	}

	if len(segs) == 0 {
		return rootScope, nil
	}

	var rid ResourceId = rootScope
	var err error

	// Root scope level resources, indicating the ARM 1st-class resource types
	if !strings.EqualFold(segs[0], "providers") {
		rid, segs, err = extendRootScopeByOneRP(rootScope, "", segs)
		if err != nil {
			return nil, fmt.Errorf("extending for root level RP: %v", err)
		}
	}

	for len(segs) != 0 {
		if !strings.EqualFold(segs[0], "providers") {
			return nil, fmt.Errorf(`scopes should be split by "/providers/"`)
		}
		segs = segs[1:]
		if len(segs) == 0 {
			return nil, fmt.Errorf("missing provider namespace segment")
		}
		rp := segs[0]
		segs = segs[1:]
		rid, segs, err = extendIdByOneRP(rid, rp, segs)
		if err != nil {
			return nil, fmt.Errorf("extending for RP %s: %v", rp, err)
		}
	}
	return rid, nil
}

func extendRootScopeByOneRP(pid RootScope, rp string, segs []string) (ResourceId, []string, error) {
	types := []string{}
	names := []string{}
	for len(segs) != 0 {
		if strings.EqualFold(segs[0], "providers") {
			break
		}
		types = append(types, segs[0])
		segs = segs[1:]
		if len(segs) == 0 {
			return nil, nil, fmt.Errorf("missing resource type name after type %s", types[len(types)-1])
		}
		names = append(names, segs[0])
		segs = segs[1:]
	}
	switch pid := pid.(type) {
	case *ManagementGroup:
		pid.AttrTypes = types
		pid.AttrNames = names
		return pid, segs, nil
	case *SubscriptionId:
		pid.AttrTypes = types
		pid.AttrNames = names
		return pid, segs, nil
	case *ResourceGroup:
		pid.AttrTypes = types
		pid.AttrNames = names
		return pid, segs, nil
	default:
		return nil, nil, fmt.Errorf("unsupported type %T", pid)
	}
}

func extendIdByOneRP(pid ResourceId, rp string, segs []string) (ResourceId, []string, error) {
	types := []string{}
	names := []string{}
	for len(segs) != 0 {
		if strings.EqualFold(segs[0], "providers") {
			break
		}
		types = append(types, segs[0])
		segs = segs[1:]
		if len(segs) == 0 {
			return nil, nil, fmt.Errorf("missing resource type name after type %s", types[len(types)-1])
		}
		names = append(names, segs[0])
		segs = segs[1:]
	}
	return &ScopedResourceId{
		AttrParentScope: pid,
		AttrProvider:    rp,
		AttrTypes:       types,
		AttrNames:       names,
	}, segs, nil
}

// RootScope is a special resource id, that represents a root scope as defined by ARM.
// This is a sealed interface, that has a limited set of implementors.
type RootScope interface {
	ResourceId
	isRootScope()
}

// TenantId represents the tenant scope, which is a pesudo resource id.
type TenantId struct{}

var _ RootScope = &TenantId{}

func (id *TenantId) RootScope() RootScope {
	return id
}

func (*TenantId) ParentScope() ResourceId {
	return nil
}

func (*TenantId) Parent() ResourceId {
	return nil
}

func (*TenantId) Provider() string {
	return "Microsoft.Resources"
}

func (*TenantId) Types() []string {
	return []string{"tenants"}
}

func (*TenantId) Names() []string {
	return nil
}

func (*TenantId) TypeString() string {
	return ""
}

func (*TenantId) String() string {
	return "/"
}

func (*TenantId) Equal(oid ResourceId) bool {
	_, ok := oid.(*TenantId)
	return ok
}

func (*TenantId) ScopeEqual(oid ResourceId) bool {
	_, ok := oid.(*TenantId)
	return ok
}

func (id *TenantId) ScopeString() string {
	return "/"
}

func (id *TenantId) RouteScopeString() string {
	return id.ScopeString()
}

func (id *TenantId) Normalize(string) error {
	return nil
}

func (id *TenantId) Clone() ResourceId {
	return &TenantId{}
}

func (*TenantId) isRootScope() {}

// SubscriptionId represents the subscription scope
type SubscriptionId struct {
	// Id is the UUID of this subscription
	Id string

	AttrTypes []string
	AttrNames []string

	subscriptionsLiteralOverride string
}

var _ RootScope = &SubscriptionId{}

func (id *SubscriptionId) subscriptionsLiteral() string {
	if id.subscriptionsLiteralOverride != "" {
		return id.subscriptionsLiteralOverride
	}
	return "subscriptions"
}

func (id *SubscriptionId) RootScope() RootScope {
	return id
}

func (*SubscriptionId) ParentScope() ResourceId {
	return nil
}

func (*SubscriptionId) Provider() string {
	return "Microsoft.Resources"
}

func (id *SubscriptionId) Parent() ResourceId {
	length := len(id.AttrTypes)
	if length == 0 {
		return nil
	}
	return &SubscriptionId{
		Id:                           id.Id,
		subscriptionsLiteralOverride: id.subscriptionsLiteralOverride,
		AttrTypes:                    id.AttrTypes[0 : length-1],
		AttrNames:                    id.AttrNames[0 : length-1],
	}
}

func (id *SubscriptionId) Types() []string {
	l := []string{id.subscriptionsLiteral()}
	l = append(l, id.AttrTypes...)
	return l
}

func (id *SubscriptionId) Names() []string {
	l := []string{id.Id}
	l = append(l, id.AttrNames...)
	return l
}

func (id *SubscriptionId) TypeString() string {
	return typeString(id)
}

func (id *SubscriptionId) String() string {
	segs := []string{id.subscriptionsLiteral(), id.Id}
	for i := range id.AttrTypes {
		segs = append(segs, id.AttrTypes[i])
		segs = append(segs, id.AttrNames[i])
	}
	return "/" + strings.Join(segs, "/")
}

func (id *SubscriptionId) Equal(oid ResourceId) bool {
	if !id.ScopeEqual(oid) {
		return false
	}
	oSubId := oid.(*SubscriptionId)
	if !strings.EqualFold(id.Id, oSubId.Id) {
		return false
	}
	for i, v := range id.AttrNames {
		if !strings.EqualFold(v, oSubId.AttrNames[i]) {
			return false
		}
	}
	return true
}

func (id *SubscriptionId) ScopeEqual(oid ResourceId) bool {
	oSubId, ok := oid.(*SubscriptionId)
	if !ok {
		return false
	}
	if len(id.AttrTypes) != len(oSubId.AttrTypes) {
		return false
	}
	for i, v := range id.AttrTypes {
		if !strings.EqualFold(v, oSubId.AttrTypes[i]) {
			return false
		}
	}
	return true
}

func (id *SubscriptionId) ScopeString() string {
	segs := []string{id.subscriptionsLiteral()}
	for i := range id.AttrTypes {
		segs = append(segs, id.AttrTypes[i])
	}
	return "/" + strings.Join(segs, "/")
}

func (id *SubscriptionId) RouteScopeString() string {
	return id.ScopeString()
}

func (id *SubscriptionId) Normalize(scopeStr string) error {
	if !strings.EqualFold(id.ScopeString(), scopeStr) {
		return fmt.Errorf("mismatch route scope string (%q) for id %q", scopeStr, id.String())
	}
	segs := strings.Split(strings.TrimPrefix(scopeStr, "/"), "/")
	id.subscriptionsLiteralOverride = segs[0]
	id.AttrTypes = segs[1:]
	return nil
}

func (id *SubscriptionId) Clone() ResourceId {
	out := &SubscriptionId{
		Id:                           id.Id,
		subscriptionsLiteralOverride: id.subscriptionsLiteralOverride,
	}
	if id.AttrTypes != nil {
		out.AttrTypes = append([]string{}, id.AttrTypes...)
	}
	if id.AttrNames != nil {
		out.AttrNames = append([]string{}, id.AttrNames...)
	}
	return out
}

func (*SubscriptionId) isRootScope() {}

// ResourceGroup represents the resource group scope
type ResourceGroup struct {
	// SubscriptionId is the UUID of the containing subscription
	SubscriptionId string
	// Name is the name of this resource group
	Name string

	AttrTypes []string
	AttrNames []string

	subscriptionsLiteralOverride  string
	resourceGroupsLiteralOverride string
}

var _ RootScope = &ResourceGroup{}

func (id *ResourceGroup) subscriptionsLiteral() string {
	if id.subscriptionsLiteralOverride != "" {
		return id.subscriptionsLiteralOverride
	}
	return "subscriptions"
}

func (id *ResourceGroup) resourceGroupsLiteral() string {
	if id.resourceGroupsLiteralOverride != "" {
		return id.resourceGroupsLiteralOverride
	}
	return "resourceGroups"
}

func (id *ResourceGroup) RootScope() RootScope {
	return id
}

func (*ResourceGroup) ParentScope() ResourceId {
	return nil
}

func (id *ResourceGroup) Parent() ResourceId {
	length := len(id.AttrTypes)
	if length == 0 {
		return nil
	}
	return &ResourceGroup{
		SubscriptionId:                id.SubscriptionId,
		Name:                          id.Name,
		subscriptionsLiteralOverride:  id.subscriptionsLiteralOverride,
		resourceGroupsLiteralOverride: id.resourceGroupsLiteralOverride,
		AttrTypes:                     id.AttrTypes[0 : length-1],
		AttrNames:                     id.AttrNames[0 : length-1],
	}
}

func (*ResourceGroup) Provider() string {
	return "Microsoft.Resources"
}

func (id *ResourceGroup) Types() []string {
	l := []string{id.subscriptionsLiteral(), id.resourceGroupsLiteral()}
	l = append(l, id.AttrTypes...)
	return l
}

func (id *ResourceGroup) Names() []string {
	l := []string{id.SubscriptionId, id.Name}
	l = append(l, id.AttrNames...)
	return l
}

func (id *ResourceGroup) TypeString() string {
	return typeString(id)
}

func (id *ResourceGroup) String() string {
	segs := []string{id.subscriptionsLiteral(), id.SubscriptionId, id.resourceGroupsLiteral(), id.Name}
	for i := range id.AttrTypes {
		segs = append(segs, id.AttrTypes[i])
		segs = append(segs, id.AttrNames[i])
	}
	return "/" + strings.Join(segs, "/")
}

func (id *ResourceGroup) Equal(oid ResourceId) bool {
	if !id.ScopeEqual(oid) {
		return false
	}
	oRgId := oid.(*ResourceGroup)
	if !strings.EqualFold(id.SubscriptionId, oRgId.SubscriptionId) {
		return false
	}
	if !strings.EqualFold(id.Name, oRgId.Name) {
		return false
	}
	for i, v := range id.AttrNames {
		if !strings.EqualFold(v, oRgId.AttrNames[i]) {
			return false
		}
	}
	return true
}

func (id *ResourceGroup) ScopeEqual(oid ResourceId) bool {
	oRgId, ok := oid.(*ResourceGroup)
	if !ok {
		return false
	}
	if len(id.AttrTypes) != len(oRgId.AttrTypes) {
		return false
	}
	for i, v := range id.AttrTypes {
		if !strings.EqualFold(v, oRgId.AttrTypes[i]) {
			return false
		}
	}
	return true
}

func (id *ResourceGroup) ScopeString() string {
	segs := []string{id.subscriptionsLiteral(), id.resourceGroupsLiteral()}
	for i := range id.AttrTypes {
		segs = append(segs, id.AttrTypes[i])
	}
	return "/" + strings.Join(segs, "/")
}

func (id *ResourceGroup) RouteScopeString() string {
	return id.ScopeString()
}

func (id *ResourceGroup) Normalize(scopeStr string) error {
	if !strings.EqualFold(id.ScopeString(), scopeStr) {
		return fmt.Errorf("mismatch route scope string (%q) for id %q", scopeStr, id.String())
	}
	segs := strings.Split(strings.TrimPrefix(scopeStr, "/"), "/")
	id.subscriptionsLiteralOverride = segs[0]
	id.resourceGroupsLiteralOverride = segs[1]
	id.AttrTypes = segs[2:]
	return nil
}

func (id *ResourceGroup) Clone() ResourceId {
	out := &ResourceGroup{
		SubscriptionId:                id.SubscriptionId,
		Name:                          id.Name,
		subscriptionsLiteralOverride:  id.subscriptionsLiteralOverride,
		resourceGroupsLiteralOverride: id.resourceGroupsLiteralOverride,
	}
	if id.AttrTypes != nil {
		out.AttrTypes = append([]string{}, id.AttrTypes...)
	}
	if id.AttrNames != nil {
		out.AttrNames = append([]string{}, id.AttrNames...)
	}
	return out
}

func (*ResourceGroup) isRootScope() {}

// ManagementGroup represents the management group scope
type ManagementGroup struct {
	// Name is the name of this management group
	Name string

	AttrTypes []string
	AttrNames []string

	microsoftManagementLiteralOverride string
	managementGroupsLiteralOverride    string
}

var _ RootScope = &ManagementGroup{}

func (id *ManagementGroup) microsoftManagementLiteral() string {
	if id.microsoftManagementLiteralOverride != "" {
		return id.microsoftManagementLiteralOverride
	}
	return "Microsoft.Management"
}

func (id *ManagementGroup) managementGroupsLiteral() string {
	if id.managementGroupsLiteralOverride != "" {
		return id.managementGroupsLiteralOverride
	}
	return "managementGroups"
}

func (id *ManagementGroup) RootScope() RootScope {
	return id
}

func (*ManagementGroup) ParentScope() ResourceId {
	return nil
}

func (id *ManagementGroup) Parent() ResourceId {
	length := len(id.AttrTypes)
	if length == 0 {
		return nil
	}
	return &ManagementGroup{
		Name:                               id.Name,
		AttrTypes:                          id.AttrTypes[0 : length-1],
		AttrNames:                          id.AttrNames[0 : length-1],
		microsoftManagementLiteralOverride: id.microsoftManagementLiteralOverride,
		managementGroupsLiteralOverride:    id.managementGroupsLiteralOverride,
	}
}

func (id *ManagementGroup) Provider() string {
	return id.microsoftManagementLiteral()
}

func (id *ManagementGroup) Types() []string {
	l := []string{id.managementGroupsLiteral()}
	l = append(l, id.AttrTypes...)
	return l
}

func (id *ManagementGroup) Names() []string {
	l := []string{id.Name}
	l = append(l, id.AttrNames...)
	return l
}

func (id *ManagementGroup) TypeString() string {
	return typeString(id)
}

func (id *ManagementGroup) String() string {
	return formatScope(id.Provider(), id.Types(), id.Names())
}

func (id *ManagementGroup) Equal(oid ResourceId) bool {
	if !id.ScopeEqual(oid) {
		return false
	}
	oMgmtId := oid.(*ManagementGroup)
	if !strings.EqualFold(id.Name, oMgmtId.Name) {
		return false
	}
	for i, v := range id.AttrNames {
		if !strings.EqualFold(v, oMgmtId.AttrNames[i]) {
			return false
		}
	}
	return true
}

func (id *ManagementGroup) ScopeEqual(oid ResourceId) bool {
	oMgmtId, ok := oid.(*ManagementGroup)
	if !ok {
		return false
	}
	if len(id.AttrTypes) != len(oMgmtId.AttrTypes) {
		return false
	}
	for i, v := range id.AttrTypes {
		if !strings.EqualFold(v, oMgmtId.AttrTypes[i]) {
			return false
		}
	}
	return true
}

func (id *ManagementGroup) ScopeString() string {
	var segs []string
	segs = append(segs, id.Provider())
	segs = append(segs, id.Types()...)
	return "/" + strings.Join(segs, "/")
}

func (id *ManagementGroup) RouteScopeString() string {
	return id.ScopeString()
}

func (id *ManagementGroup) Normalize(scopeStr string) error {
	if !strings.EqualFold(id.ScopeString(), scopeStr) {
		return fmt.Errorf("mismatch route scope string (%q) for id %q", scopeStr, id.String())
	}
	segs := strings.Split(strings.TrimPrefix(scopeStr, "/"), "/")
	id.microsoftManagementLiteralOverride = segs[0]
	id.managementGroupsLiteralOverride = segs[1]
	id.AttrTypes = segs[2:]
	return nil
}

func (id *ManagementGroup) Clone() ResourceId {
	out := &ManagementGroup{
		Name:                               id.Name,
		microsoftManagementLiteralOverride: id.microsoftManagementLiteralOverride,
		managementGroupsLiteralOverride:    id.managementGroupsLiteralOverride,
	}
	if id.AttrTypes != nil {
		out.AttrTypes = append([]string{}, id.AttrTypes...)
	}
	if id.AttrNames != nil {
		out.AttrNames = append([]string{}, id.AttrNames...)
	}
	return out
}

func (ManagementGroup) isRootScope() {}

// ScopedResourceId represents a resource id that is scoped within a root scope or another scoped resource.
var _ ResourceId = &ScopedResourceId{}

type ScopedResourceId struct {
	AttrParentScope ResourceId
	AttrProvider    string
	AttrTypes       []string
	AttrNames       []string
}

func (id *ScopedResourceId) RootScope() RootScope {
	var rid ResourceId = id
	for rid.ParentScope() != nil {
		rid = rid.ParentScope()
	}
	return rid.(RootScope)
}

func (id *ScopedResourceId) ParentScope() ResourceId {
	return id.AttrParentScope
}

func (id *ScopedResourceId) Parent() ResourceId {
	length := len(id.AttrTypes)
	if length == 0 {
		return nil
	}
	return &ScopedResourceId{
		AttrParentScope: id.AttrParentScope,
		AttrProvider:    id.AttrProvider,
		AttrTypes:       id.AttrTypes[0 : length-1],
		AttrNames:       id.AttrNames[0 : length-1],
	}
}

func (id *ScopedResourceId) Provider() string {
	return id.AttrProvider
}

func (id *ScopedResourceId) Types() []string {
	return id.AttrTypes
}

func (id *ScopedResourceId) Names() []string {
	return id.AttrNames
}

func (id *ScopedResourceId) TypeString() string {
	return typeString(id)
}

func (id *ScopedResourceId) String() string {
	builder := strings.Builder{}
	if _, ok := id.ParentScope().(*TenantId); !ok {
		builder.WriteString(id.ParentScope().String())
	}
	builder.WriteString(formatScope(id.Provider(), id.Types(), id.Names()))
	return builder.String()
}

func (id *ScopedResourceId) Equal(oid ResourceId) bool {
	oRid, ok := oid.(*ScopedResourceId)
	if !ok {
		return false
	}
	if !id.AttrParentScope.Equal(oRid.AttrParentScope) {
		return false
	}
	if !id.ScopeEqual(oid) {
		return false
	}
	if len(id.AttrNames) != len(oRid.AttrNames) {
		return false
	}
	for i := 0; i < len(id.AttrNames); i++ {
		if !strings.EqualFold(id.AttrNames[i], oRid.AttrNames[i]) {
			return false
		}
	}
	return true
}

func (id *ScopedResourceId) ScopeEqual(oid ResourceId) bool {
	oRid, ok := oid.(*ScopedResourceId)
	if !ok {
		return false
	}
	if !id.AttrParentScope.ScopeEqual(oRid.AttrParentScope) {
		return false
	}
	if !strings.EqualFold(id.AttrProvider, oRid.AttrProvider) {
		return false
	}
	if len(id.AttrTypes) != len(oRid.AttrTypes) {
		return false
	}
	for i := 0; i < len(id.AttrTypes); i++ {
		if !strings.EqualFold(id.AttrTypes[i], oRid.AttrTypes[i]) {
			return false
		}
	}
	return true
}

func (id *ScopedResourceId) ScopeString() string {
	var out string
	traverseScopes(id, func(id ResourceId) {
		if _, ok := id.(*TenantId); ok {
			return
		}
		out += id.RouteScopeString()
	})
	return out
}

func (id *ScopedResourceId) RouteScopeString() string {
	var segs []string
	segs = append(segs, id.Provider())
	segs = append(segs, id.Types()...)
	return "/" + strings.Join(segs, "/")
}

// Normalize normalizes the invariant parts (e.g. Provider, Types) of the id  based on the input scope string.
// The input scope string must be the same as calling the `ScopeString` of this id, except the casing.
func (id *ScopedResourceId) Normalize(scopeStr string) error {
	if !strings.EqualFold(id.ScopeString(), scopeStr) {
		return fmt.Errorf("mismatch scope string (%q) for id %q", scopeStr, id.String())
	}

	traverseScopes(id, func(id ResourceId) {
		if id.ParentScope() == nil {
			if _, ok := id.(*TenantId); !ok {
				id.Normalize(scopeStr[:len(id.ScopeString())])
				scopeStr = scopeStr[len(id.ScopeString()):]
			}
			return
		}
		segs := []string{id.Provider()}
		segs = append(segs, id.Types()...)
		thisScopeStrLen := len("/" + strings.Join(segs, "/"))

		var thisScopeStr string
		thisScopeStr, scopeStr = scopeStr[0:thisScopeStrLen], scopeStr[thisScopeStrLen:]

		segs = strings.Split(strings.TrimPrefix(thisScopeStr, "/"), "/")
		rid := id.(*ScopedResourceId)
		rid.AttrProvider = segs[0]
		rid.AttrTypes = segs[1:]
	})
	return nil
}

func (id *ScopedResourceId) Clone() ResourceId {
	out := &ScopedResourceId{
		AttrParentScope: id.ParentScope().Clone(),
		AttrProvider:    id.AttrProvider,
	}
	if id.AttrTypes != nil {
		out.AttrTypes = append([]string{}, id.AttrTypes...)
	}
	if id.AttrNames != nil {
		out.AttrNames = append([]string{}, id.AttrNames...)
	}
	return out
}

// NormalizeRouteScope is similar to Normalize, while only for the current route scope, and didn't affect the parent scopes.
// The input scope string must be the same as the ScopeString of this of this id with its parent scope's ScopeString prefix trimmed, except the casing.
func (id *ScopedResourceId) NormalizeRouteScope(scopeStr string) error {
	if !strings.EqualFold(id.RouteScopeString(), scopeStr) {
		return fmt.Errorf("mismatch route scope string (%q) for id %q", scopeStr, id.String())
	}
	segs := strings.Split(strings.TrimPrefix(scopeStr, "/"), "/")
	id.AttrProvider = segs[0]
	id.AttrTypes = segs[1:]
	return nil
}

func formatScope(provider string, types []string, names []string) string {
	if len(types) != len(names) {
		panic(fmt.Sprintf("invalid input: len(%v) != len(%v)", types, names))
	}
	var segs []string
	segs = append(segs, "/providers", provider)
	for i := 0; i < len(types); i++ {
		segs = append(segs, types[i])
		segs = append(segs, names[i])
	}
	return strings.Join(segs, "/")
}

func typeString(id ResourceId) string {
	var segs []string
	segs = append(segs, id.Provider())
	segs = append(segs, id.Types()...)
	return strings.Join(segs, "/")
}

// traverScopes traverse the scopes of the given resource id from the root scope to the router scope.
func traverseScopes(id ResourceId, f func(ResourceId)) {
	for ; id != nil; id = id.ParentScope() {
		id := id
		defer func() {
			f(id)
		}()
	}
}
