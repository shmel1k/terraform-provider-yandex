// Code generated by protoc-gen-goext. DO NOT EDIT.

package compute

import (
	operation "github.com/yandex-cloud/go-genproto/yandex/cloud/operation"
	field_mask "google.golang.org/genproto/protobuf/field_mask"
)

func (m *GetHostGroupRequest) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *ListHostGroupsRequest) SetFolderId(v string) {
	m.FolderId = v
}

func (m *ListHostGroupsRequest) SetPageSize(v int64) {
	m.PageSize = v
}

func (m *ListHostGroupsRequest) SetPageToken(v string) {
	m.PageToken = v
}

func (m *ListHostGroupsRequest) SetFilter(v string) {
	m.Filter = v
}

func (m *ListHostGroupsResponse) SetHostGroups(v []*HostGroup) {
	m.HostGroups = v
}

func (m *ListHostGroupsResponse) SetNextPageToken(v string) {
	m.NextPageToken = v
}

func (m *CreateHostGroupRequest) SetFolderId(v string) {
	m.FolderId = v
}

func (m *CreateHostGroupRequest) SetName(v string) {
	m.Name = v
}

func (m *CreateHostGroupRequest) SetDescription(v string) {
	m.Description = v
}

func (m *CreateHostGroupRequest) SetLabels(v map[string]string) {
	m.Labels = v
}

func (m *CreateHostGroupRequest) SetZoneId(v string) {
	m.ZoneId = v
}

func (m *CreateHostGroupRequest) SetTypeId(v string) {
	m.TypeId = v
}

func (m *CreateHostGroupRequest) SetMaintenancePolicy(v MaintenancePolicy) {
	m.MaintenancePolicy = v
}

func (m *CreateHostGroupRequest) SetScalePolicy(v *ScalePolicy) {
	m.ScalePolicy = v
}

func (m *CreateHostGroupMetadata) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *UpdateHostGroupRequest) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *UpdateHostGroupRequest) SetUpdateMask(v *field_mask.FieldMask) {
	m.UpdateMask = v
}

func (m *UpdateHostGroupRequest) SetName(v string) {
	m.Name = v
}

func (m *UpdateHostGroupRequest) SetDescription(v string) {
	m.Description = v
}

func (m *UpdateHostGroupRequest) SetLabels(v map[string]string) {
	m.Labels = v
}

func (m *UpdateHostGroupRequest) SetMaintenancePolicy(v MaintenancePolicy) {
	m.MaintenancePolicy = v
}

func (m *UpdateHostGroupRequest) SetScalePolicy(v *ScalePolicy) {
	m.ScalePolicy = v
}

func (m *UpdateHostGroupMetadata) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *DeleteHostGroupRequest) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *DeleteHostGroupMetadata) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *ListHostGroupInstancesRequest) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *ListHostGroupInstancesRequest) SetPageSize(v int64) {
	m.PageSize = v
}

func (m *ListHostGroupInstancesRequest) SetPageToken(v string) {
	m.PageToken = v
}

func (m *ListHostGroupInstancesRequest) SetFilter(v string) {
	m.Filter = v
}

func (m *ListHostGroupInstancesResponse) SetInstances(v []*Instance) {
	m.Instances = v
}

func (m *ListHostGroupInstancesResponse) SetNextPageToken(v string) {
	m.NextPageToken = v
}

func (m *ListHostGroupHostsRequest) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *ListHostGroupHostsRequest) SetPageSize(v int64) {
	m.PageSize = v
}

func (m *ListHostGroupHostsRequest) SetPageToken(v string) {
	m.PageToken = v
}

func (m *ListHostGroupHostsResponse) SetHosts(v []*Host) {
	m.Hosts = v
}

func (m *ListHostGroupHostsResponse) SetNextPageToken(v string) {
	m.NextPageToken = v
}

func (m *ListHostGroupOperationsRequest) SetHostGroupId(v string) {
	m.HostGroupId = v
}

func (m *ListHostGroupOperationsRequest) SetPageSize(v int64) {
	m.PageSize = v
}

func (m *ListHostGroupOperationsRequest) SetPageToken(v string) {
	m.PageToken = v
}

func (m *ListHostGroupOperationsResponse) SetOperations(v []*operation.Operation) {
	m.Operations = v
}

func (m *ListHostGroupOperationsResponse) SetNextPageToken(v string) {
	m.NextPageToken = v
}
