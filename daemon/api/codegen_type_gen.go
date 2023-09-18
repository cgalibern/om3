// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/opensvc/om3/core/instance"
	"github.com/opensvc/om3/core/node"
	"github.com/opensvc/om3/core/resource"
)

const (
	BasicAuthScopes  = "basicAuth.Scopes"
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for Orchestrate.
const (
	OrchestrateHa    Orchestrate = "ha"
	OrchestrateNo    Orchestrate = "no"
	OrchestrateStart Orchestrate = "start"
)

// Defines values for PlacementPolicy.
const (
	PlacementPolicyLastStart  PlacementPolicy = "last start"
	PlacementPolicyLoadAvg    PlacementPolicy = "load avg"
	PlacementPolicyNodesOrder PlacementPolicy = "nodes order"
	PlacementPolicyNone       PlacementPolicy = "none"
	PlacementPolicyScore      PlacementPolicy = "score"
	PlacementPolicyShift      PlacementPolicy = "shift"
	PlacementPolicySpread     PlacementPolicy = "spread"
)

// Defines values for PlacementState.
const (
	PlacementStateNa         PlacementState = "n/a"
	PlacementStateNonOptimal PlacementState = "non-optimal"
	PlacementStateOptimal    PlacementState = "optimal"
	PlacementStateUndef      PlacementState = "undef"
)

// Defines values for PostDaemonLogsControlLevel.
const (
	PostDaemonLogsControlLevelDebug PostDaemonLogsControlLevel = "debug"
	PostDaemonLogsControlLevelError PostDaemonLogsControlLevel = "error"
	PostDaemonLogsControlLevelFatal PostDaemonLogsControlLevel = "fatal"
	PostDaemonLogsControlLevelInfo  PostDaemonLogsControlLevel = "info"
	PostDaemonLogsControlLevelNone  PostDaemonLogsControlLevel = "none"
	PostDaemonLogsControlLevelPanic PostDaemonLogsControlLevel = "panic"
	PostDaemonLogsControlLevelWarn  PostDaemonLogsControlLevel = "warn"
)

// Defines values for PostDaemonSubActionAction.
const (
	PostDaemonSubActionActionStart PostDaemonSubActionAction = "start"
	PostDaemonSubActionActionStop  PostDaemonSubActionAction = "stop"
)

// Defines values for Provisioned.
const (
	ProvisionedFalse Provisioned = "false"
	ProvisionedMixed Provisioned = "mixed"
	ProvisionedNa    Provisioned = "n/a"
	ProvisionedTrue  Provisioned = "true"
)

// Defines values for Role.
const (
	Admin          Role = "admin"
	Blacklistadmin Role = "blacklistadmin"
	Guest          Role = "guest"
	Heartbeat      Role = "heartbeat"
	Join           Role = "join"
	Leave          Role = "leave"
	Root           Role = "root"
	Squatter       Role = "squatter"
)

// Defines values for Status.
const (
	Down      Status = "down"
	Na        Status = "n/a"
	StdbyDown Status = "stdby down"
	StdbyUp   Status = "stdby up"
	Undef     Status = "undef"
	Up        Status = "up"
	Warn      Status = "warn"
)

// Defines values for Topology.
const (
	Failover Topology = "failover"
	Flex     Topology = "flex"
)

// ArbitratorStatus defines model for ArbitratorStatus.
type ArbitratorStatus struct {
	Status Status `json:"status"`
	Url    string `json:"url"`
}

// AuthToken defines model for AuthToken.
type AuthToken struct {
	ExpiredAt time.Time `json:"expired_at"`
	Token     string    `json:"token"`
}

// Cluster defines model for Cluster.
type Cluster struct {
	Config ClusterConfig `json:"config"`
	Node   ClusterNode   `json:"node"`
	Object ClusterObject `json:"object"`
	Status ClusterStatus `json:"status"`
}

// ClusterConfig defines model for ClusterConfig.
type ClusterConfig = map[string]interface{}

// ClusterNode defines model for ClusterNode.
type ClusterNode = map[string]interface{}

// ClusterObject defines model for ClusterObject.
type ClusterObject = map[string]interface{}

// ClusterStatus defines model for ClusterStatus.
type ClusterStatus = map[string]interface{}

// DNSRecord defines model for DNSRecord.
type DNSRecord struct {
	Class string `json:"class"`
	Data  string `json:"data"`
	Name  string `json:"name"`
	Ttl   int    `json:"ttl"`
	Type  string `json:"type"`
}

// DNSZone defines model for DNSZone.
type DNSZone = []DNSRecord

// DRBDAllocation defines model for DRBDAllocation.
type DRBDAllocation struct {
	ExpiredAt time.Time          `json:"expired_at"`
	Id        openapi_types.UUID `json:"id"`
	Minor     int                `json:"minor"`
	Port      int                `json:"port"`
}

// DRBDConfig defines model for DRBDConfig.
type DRBDConfig struct {
	Data []byte `json:"data"`
}

// Daemon defines model for Daemon.
type Daemon struct {
	Collector DaemonCollector `json:"collector"`
	Dns       DaemonDNS       `json:"dns"`
	Hb        DaemonHb        `json:"hb"`
	Listener  DaemonListener  `json:"listener"`
	Monitor   DaemonMonitor   `json:"monitor"`
	Routines  int             `json:"routines"`
	Scheduler DaemonScheduler `json:"scheduler"`
}

// DaemonCollector defines model for DaemonCollector.
type DaemonCollector = DaemonSubsystemStatus

// DaemonDNS defines model for DaemonDNS.
type DaemonDNS = DaemonSubsystemStatus

// DaemonHb defines model for DaemonHb.
type DaemonHb struct {
	Modes   []DaemonHbMode   `json:"modes"`
	Streams []DaemonHbStream `json:"streams"`
}

// DaemonHbMode defines model for DaemonHbMode.
type DaemonHbMode struct {
	// Mode the type of hb message used by node except when Type is patch where it is the patch queue length
	Mode string `json:"mode"`

	// Node a cluster node
	Node string `json:"node"`

	// Type the heartbeat message type used by node
	Type string `json:"type"`
}

// DaemonHbStream defines model for DaemonHbStream.
type DaemonHbStream struct {
	Alerts     []DaemonSubsystemAlert `json:"alerts"`
	Configured time.Time              `json:"configured"`
	CreatedAt  time.Time              `json:"created_at"`
	Id         string                 `json:"id"`
	IsBeating  bool                   `json:"is_beating"`
	LastAt     time.Time              `json:"last_at"`
	State      string                 `json:"state"`

	// Type hb stream type
	Type string `json:"type"`
}

// DaemonHbStreamPeer defines model for DaemonHbStreamPeer.
type DaemonHbStreamPeer struct {
	IsBeating bool      `json:"is_beating"`
	LastAt    time.Time `json:"last_at"`
}

// DaemonHbStreamType defines model for DaemonHbStreamType.
type DaemonHbStreamType struct {
	// Type hb stream type
	Type string `json:"type"`
}

// DaemonListener defines model for DaemonListener.
type DaemonListener = DaemonSubsystemStatus

// DaemonMonitor defines model for DaemonMonitor.
type DaemonMonitor = DaemonSubsystemStatus

// DaemonScheduler defines model for DaemonScheduler.
type DaemonScheduler = DaemonSubsystemStatus

// DaemonStatus defines model for DaemonStatus.
type DaemonStatus struct {
	Cluster Cluster `json:"cluster"`
	Daemon  Daemon  `json:"daemon"`
}

// DaemonSubsystemAlert defines model for DaemonSubsystemAlert.
type DaemonSubsystemAlert struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// DaemonSubsystemStatus defines model for DaemonSubsystemStatus.
type DaemonSubsystemStatus struct {
	Alerts     []DaemonSubsystemAlert `json:"alerts"`
	Configured time.Time              `json:"configured"`
	CreatedAt  time.Time              `json:"created_at"`
	Id         string                 `json:"id"`
	State      string                 `json:"state"`
}

// EventList responseEventList is a list of sse
type EventList = openapi_types.File

// Instance defines model for Instance.
type Instance struct {
	Config  *InstanceConfig  `json:"config,omitempty"`
	Monitor *InstanceMonitor `json:"monitor,omitempty"`
	Status  *InstanceStatus  `json:"status,omitempty"`
}

// InstanceArray defines model for InstanceArray.
type InstanceArray = []InstanceItem

// InstanceConfig defines model for InstanceConfig.
type InstanceConfig = instance.Config

// InstanceConfigArray defines model for InstanceConfigArray.
type InstanceConfigArray = []InstanceConfigItem

// InstanceConfigItem defines model for InstanceConfigItem.
type InstanceConfigItem struct {
	Data InstanceConfig `json:"data"`
	Meta InstanceMeta   `json:"meta"`
}

// InstanceItem defines model for InstanceItem.
type InstanceItem struct {
	Data Instance     `json:"data"`
	Meta InstanceMeta `json:"meta"`
}

// InstanceMeta defines model for InstanceMeta.
type InstanceMeta struct {
	Node   string `json:"node"`
	Object string `json:"object"`
}

// InstanceMonitor defines model for InstanceMonitor.
type InstanceMonitor = instance.Monitor

// InstanceMonitorArray defines model for InstanceMonitorArray.
type InstanceMonitorArray = []InstanceMonitorItem

// InstanceMonitorItem defines model for InstanceMonitorItem.
type InstanceMonitorItem struct {
	Data InstanceMonitor `json:"data"`
	Meta InstanceMeta    `json:"meta"`
}

// InstanceStatus defines model for InstanceStatus.
type InstanceStatus = instance.Status

// InstanceStatusArray defines model for InstanceStatusArray.
type InstanceStatusArray = []InstanceStatusItem

// InstanceStatusItem defines model for InstanceStatusItem.
type InstanceStatusItem struct {
	Data InstanceStatus `json:"data"`
	Meta InstanceMeta   `json:"meta"`
}

// LogList responseLogList is a list of sse
type LogList = openapi_types.File

// Network defines model for Network.
type Network struct {
	Errors  []string     `json:"errors"`
	Name    string       `json:"name"`
	Network string       `json:"network"`
	Type    string       `json:"type"`
	Usage   NetworkUsage `json:"usage"`
}

// NetworkArray defines model for NetworkArray.
type NetworkArray = []Network

// NetworkIp defines model for NetworkIp.
type NetworkIp struct {
	Ip      string           `json:"ip"`
	Network NetworkIpNetwork `json:"network"`
	Node    string           `json:"node"`
	Path    string           `json:"path"`
	Rid     string           `json:"rid"`
}

// NetworkIpArray defines model for NetworkIpArray.
type NetworkIpArray = []NetworkIp

// NetworkIpNetwork defines model for NetworkIpNetwork.
type NetworkIpNetwork struct {
	Name    string `json:"name"`
	Network string `json:"network"`
	Type    string `json:"type"`
}

// NetworkUsage defines model for NetworkUsage.
type NetworkUsage struct {
	Free int     `json:"free"`
	Pct  float32 `json:"pct"`
	Size int     `json:"size"`
	Used int     `json:"used"`
}

// Node defines model for Node.
type Node struct {
	Config  *NodeConfig  `json:"config,omitempty"`
	Monitor *NodeMonitor `json:"monitor,omitempty"`
	Status  *NodeStatus  `json:"status,omitempty"`
}

// NodeArray defines model for NodeArray.
type NodeArray = []NodeItem

// NodeConfig defines model for NodeConfig.
type NodeConfig = node.Config

// NodeConfigArray defines model for NodeConfigArray.
type NodeConfigArray = []NodeConfigItem

// NodeConfigItem defines model for NodeConfigItem.
type NodeConfigItem struct {
	Data NodeConfig `json:"data"`
	Meta NodeMeta   `json:"meta"`
}

// NodeInfo defines model for NodeInfo.
type NodeInfo struct {
	// Labels labels is the list of node labels.
	Labels []NodeLabel `json:"labels"`

	// Nodename nodename is the name of the node where the labels and paths are coming from.
	Nodename string `json:"nodename"`

	// Paths paths is the list of node to storage array san paths.
	Paths []SANPath `json:"paths"`
}

// NodeItem defines model for NodeItem.
type NodeItem struct {
	Data Node     `json:"data"`
	Meta NodeMeta `json:"meta"`
}

// NodeLabel defines model for NodeLabel.
type NodeLabel struct {
	// Name name is the label name.
	Name string `json:"name"`

	// Value value is the label value.
	Value string `json:"value"`
}

// NodeMeta defines model for NodeMeta.
type NodeMeta struct {
	Node string `json:"node"`
}

// NodeMonitor defines model for NodeMonitor.
type NodeMonitor = node.Monitor

// NodeMonitorArray defines model for NodeMonitorArray.
type NodeMonitorArray = []NodeMonitorItem

// NodeMonitorItem defines model for NodeMonitorItem.
type NodeMonitorItem struct {
	Data NodeMonitor `json:"data"`
	Meta NodeMeta    `json:"meta"`
}

// NodeStatus defines model for NodeStatus.
type NodeStatus = node.Status

// NodeStatusArray defines model for NodeStatusArray.
type NodeStatusArray = []NodeStatusItem

// NodeStatusItem defines model for NodeStatusItem.
type NodeStatusItem struct {
	Data NodeStatus `json:"data"`
	Meta NodeMeta   `json:"meta"`
}

// NodesInfo defines model for NodesInfo.
type NodesInfo = []NodeInfo

// ObjectArray defines model for ObjectArray.
type ObjectArray = []ObjectItem

// ObjectConfig defines model for ObjectConfig.
type ObjectConfig struct {
	Data  map[string]interface{} `json:"data"`
	Mtime time.Time              `json:"mtime"`
}

// ObjectData defines model for ObjectData.
type ObjectData struct {
	Avail       Status              `json:"avail"`
	FlexMax     int                 `json:"flex_max"`
	FlexMin     int                 `json:"flex_min"`
	FlexTarget  int                 `json:"flex_target"`
	Frozen      string              `json:"frozen"`
	Instances   map[string]Instance `json:"instances"`
	Orchestrate Orchestrate         `json:"orchestrate"`
	Overall     Status              `json:"overall"`

	// PlacementPolicy object placement policy
	PlacementPolicy PlacementPolicy `json:"placement_policy"`

	// PlacementState object placement state
	PlacementState PlacementState `json:"placement_state"`
	Pool           *string        `json:"pool,omitempty"`
	Priority       int            `json:"priority"`

	// Provisioned service, instance or resource provisioned state
	Provisioned Provisioned `json:"provisioned"`
	Scope       []string    `json:"scope"`
	Size        *int64      `json:"size,omitempty"`

	// Topology object topology
	Topology         Topology `json:"topology"`
	UpInstancesCount int      `json:"up_instances_count"`
	UpdatedAt        string   `json:"updated_at"`
}

// ObjectFile defines model for ObjectFile.
type ObjectFile struct {
	Data  []byte    `json:"data"`
	Mtime time.Time `json:"mtime"`
}

// ObjectItem defines model for ObjectItem.
type ObjectItem struct {
	Data ObjectData `json:"data"`
	Meta ObjectMeta `json:"meta"`
}

// ObjectMeta defines model for ObjectMeta.
type ObjectMeta struct {
	Object string `json:"object"`
}

// ObjectPaths defines model for ObjectPaths.
type ObjectPaths = []string

// Orchestrate defines model for Orchestrate.
type Orchestrate string

// OrchestrationQueued defines model for OrchestrationQueued.
type OrchestrationQueued struct {
	OrchestrationId openapi_types.UUID `json:"orchestration_id"`
}

// PlacementPolicy object placement policy
type PlacementPolicy string

// PlacementState object placement state
type PlacementState string

// Pool defines model for Pool.
type Pool struct {
	Capabilities []string  `json:"capabilities"`
	Errors       *[]string `json:"errors,omitempty"`
	Free         int64     `json:"free"`
	Head         string    `json:"head"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	Type         string    `json:"type"`
	Used         int64     `json:"used"`
	VolumeCount  int       `json:"volume_count"`
}

// PoolArray defines model for PoolArray.
type PoolArray = []Pool

// PoolVolume defines model for PoolVolume.
type PoolVolume struct {
	Children []string `json:"children"`
	IsOrphan bool     `json:"is_orphan"`
	Path     string   `json:"path"`
	Pool     string   `json:"pool"`
	Size     int64    `json:"size"`
}

// PoolVolumeArray defines model for PoolVolumeArray.
type PoolVolumeArray = []PoolVolume

// PostDaemonLogsControl defines model for PostDaemonLogsControl.
type PostDaemonLogsControl struct {
	Level PostDaemonLogsControlLevel `json:"level"`
}

// PostDaemonLogsControlLevel defines model for PostDaemonLogsControl.Level.
type PostDaemonLogsControlLevel string

// PostDaemonSubAction defines model for PostDaemonSubAction.
type PostDaemonSubAction struct {
	Action PostDaemonSubActionAction `json:"action"`

	// Subs daemon component list
	Subs []string `json:"subs"`
}

// PostDaemonSubActionAction defines model for PostDaemonSubAction.Action.
type PostDaemonSubActionAction string

// PostNodeDRBDConfigRequest defines model for PostNodeDRBDConfigRequest.
type PostNodeDRBDConfigRequest struct {
	AllocationId openapi_types.UUID `json:"allocation_id"`
	Data         []byte             `json:"data"`
}

// PostObjectAction defines model for PostObjectAction.
type PostObjectAction struct {
	Path string `json:"path"`
}

// PostObjectActionSwitch defines model for PostObjectActionSwitch.
type PostObjectActionSwitch struct {
	Destination []string `json:"destination"`
	Path        string   `json:"path"`
}

// PostObjectClear defines model for PostObjectClear.
type PostObjectClear struct {
	Path string `json:"path"`
}

// PostObjectProgress defines model for PostObjectProgress.
type PostObjectProgress struct {
	IsPartial *bool              `json:"is_partial,omitempty"`
	Path      string             `json:"path"`
	SessionId openapi_types.UUID `json:"session_id"`
	State     string             `json:"state"`
}

// PostRelayMessage defines model for PostRelayMessage.
type PostRelayMessage struct {
	ClusterId   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Msg         string `json:"msg"`
	Nodename    string `json:"nodename"`
}

// Problem defines model for Problem.
type Problem struct {
	// Detail A human-readable explanation specific to this occurrence of the
	// problem.
	Detail string `json:"detail"`

	// Status The HTTP status code ([RFC7231], Section 6) generated by the
	// origin server for this occurrence of the problem.
	Status int `json:"status"`

	// Title A short, human-readable summary of the problem type.  It SHOULD
	// NOT change from occurrence to occurrence of the problem, except
	// for purposes of localization (e.g., using proactive content
	// negotiation; see [RFC7231], Section 3.4).
	Title string `json:"title"`
}

// Provisioned service, instance or resource provisioned state
type Provisioned string

// RelayMessage defines model for RelayMessage.
type RelayMessage struct {
	Addr        string    `json:"addr"`
	ClusterId   string    `json:"cluster_id"`
	ClusterName string    `json:"cluster_name"`
	Msg         string    `json:"msg"`
	Nodename    string    `json:"nodename"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RelayMessageList defines model for RelayMessageList.
type RelayMessageList = []RelayMessage

// RelayMessages defines model for RelayMessages.
type RelayMessages struct {
	Messages RelayMessageList `json:"messages"`
}

// Resource defines model for Resource.
type Resource struct {
	Config  *ResourceConfig  `json:"config,omitempty"`
	Monitor *ResourceMonitor `json:"monitor,omitempty"`
	Status  *ResourceStatus  `json:"status,omitempty"`
}

// ResourceArray defines model for ResourceArray.
type ResourceArray = []ResourceItem

// ResourceConfig defines model for ResourceConfig.
type ResourceConfig = instance.ResourceConfig

// ResourceConfigArray defines model for ResourceConfigArray.
type ResourceConfigArray = []ResourceConfigItem

// ResourceConfigItem defines model for ResourceConfigItem.
type ResourceConfigItem struct {
	Data ResourceConfig `json:"data"`
	Meta ResourceMeta   `json:"meta"`
}

// ResourceItem defines model for ResourceItem.
type ResourceItem struct {
	Data Resource     `json:"data"`
	Meta ResourceMeta `json:"meta"`
}

// ResourceLog defines model for ResourceLog.
type ResourceLog = []ResourceLogEntry

// ResourceLogEntry defines model for ResourceLogEntry.
type ResourceLogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// ResourceMeta defines model for ResourceMeta.
type ResourceMeta struct {
	Node   string `json:"node"`
	Object string `json:"object"`
	Rid    string `json:"rid"`
}

// ResourceMonitor defines model for ResourceMonitor.
type ResourceMonitor = instance.ResourceMonitor

// ResourceMonitorArray defines model for ResourceMonitorArray.
type ResourceMonitorArray = []ResourceMonitorItem

// ResourceMonitorItem defines model for ResourceMonitorItem.
type ResourceMonitorItem struct {
	Data ResourceMonitor `json:"data"`
	Meta ResourceMeta    `json:"meta"`
}

// ResourceMonitorRestart defines model for ResourceMonitorRestart.
type ResourceMonitorRestart struct {
	LastAt    time.Time `json:"last_at"`
	Remaining int       `json:"remaining"`
}

// ResourceProvisionStatus defines model for ResourceProvisionStatus.
type ResourceProvisionStatus struct {
	Mtime time.Time `json:"mtime"`

	// State service, instance or resource provisioned state
	State Provisioned `json:"state"`
}

// ResourceStatus defines model for ResourceStatus.
type ResourceStatus = resource.Status

// ResourceStatusArray defines model for ResourceStatusArray.
type ResourceStatusArray = []ResourceStatusItem

// ResourceStatusItem defines model for ResourceStatusItem.
type ResourceStatusItem struct {
	Data ResourceStatus `json:"data"`
	Meta ResourceMeta   `json:"meta"`
}

// Role defines model for Role.
type Role string

// SANPath defines model for SANPath.
type SANPath struct {
	// Initiator initiator is the host side san path endpoint.
	Initiator SANPathInitiator `json:"initiator"`

	// Target target is the storage array side san path endpoint.
	Target SANPathTarget `json:"target"`
}

// SANPathInitiator initiator is the host side san path endpoint.
type SANPathInitiator struct {
	// Name name is a worldwide unique path endpoint identifier.
	Name *string `json:"name,omitempty"`

	// Type type is the endpoint type.
	Type *string `json:"type,omitempty"`
}

// SANPathTarget target is the storage array side san path endpoint.
type SANPathTarget struct {
	// Name name is a worldwide unique path endpoint identifier.
	Name *string `json:"name,omitempty"`

	// Type type is a the endpoint type.
	Type *string `json:"type,omitempty"`
}

// Status defines model for Status.
type Status string

// SubsetConfig defines model for SubsetConfig.
type SubsetConfig = instance.SubsetConfig

// SubsetsConfig defines model for SubsetsConfig.
type SubsetsConfig = []SubsetConfig

// Topology object topology
type Topology string

// DRBDConfigName defines model for DRBDConfigName.
type DRBDConfigName = string

// Duration defines model for Duration.
type Duration = string

// EventFilter defines model for EventFilter.
type EventFilter = []string

// Limit defines model for Limit.
type Limit = int64

// LogFilter defines model for LogFilter.
type LogFilter = []string

// NamespaceOptional defines model for NamespaceOptional.
type NamespaceOptional = string

// NodeOptional defines model for NodeOptional.
type NodeOptional = string

// ObjectPath defines model for ObjectPath.
type ObjectPath = string

// Path defines model for Path.
type Path = string

// PathOptional defines model for PathOptional.
type PathOptional = string

// Paths defines model for Paths.
type Paths = []string

// RelayClusterId defines model for RelayClusterId.
type RelayClusterId = string

// RelayNodename defines model for RelayNodename.
type RelayNodename = string

// RidOptional defines model for RidOptional.
type RidOptional = string

// Roles defines model for Roles.
type Roles = []Role

// SelectorOptional defines model for SelectorOptional.
type SelectorOptional = string

// N200 defines model for 200.
type N200 = Problem

// N400 defines model for 400.
type N400 = Problem

// N401 defines model for 401.
type N401 = Problem

// N403 defines model for 403.
type N403 = Problem

// N408 defines model for 408.
type N408 = Problem

// N409 defines model for 409.
type N409 = Problem

// N500 defines model for 500.
type N500 = Problem

// N503 defines model for 503.
type N503 = Problem

// PostAuthTokenParams defines parameters for PostAuthToken.
type PostAuthTokenParams struct {
	// Role list of api role
	Role *Roles `form:"role,omitempty" json:"role,omitempty"`

	// Duration max token duration, maximum value 24h
	Duration *string `form:"duration,omitempty" json:"duration,omitempty"`
}

// GetDaemonEventsParams defines parameters for GetDaemonEvents.
type GetDaemonEventsParams struct {
	// Duration max duration
	Duration *Duration `form:"duration,omitempty" json:"duration,omitempty"`

	// Limit limit items count
	Limit *Limit `form:"limit,omitempty" json:"limit,omitempty"`

	// Filter list of event filter
	Filter *EventFilter `form:"filter,omitempty" json:"filter,omitempty"`
}

// PostDaemonJoinParams defines parameters for PostDaemonJoin.
type PostDaemonJoinParams struct {
	// Node The node to add to cluster nodes
	Node string `form:"node" json:"node"`
}

// PostDaemonLeaveParams defines parameters for PostDaemonLeave.
type PostDaemonLeaveParams struct {
	// Node The leaving cluster node
	Node string `form:"node" json:"node"`
}

// GetDaemonStatusParams defines parameters for GetDaemonStatus.
type GetDaemonStatusParams struct {
	// Namespace namespace
	Namespace *NamespaceOptional `form:"namespace,omitempty" json:"namespace,omitempty"`

	// Selector selector
	Selector *SelectorOptional `form:"selector,omitempty" json:"selector,omitempty"`
}

// GetInstanceParams defines parameters for GetInstance.
type GetInstanceParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`

	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`
}

// GetInstanceConfigParams defines parameters for GetInstanceConfig.
type GetInstanceConfigParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`

	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`
}

// GetInstanceMonitorParams defines parameters for GetInstanceMonitor.
type GetInstanceMonitorParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`

	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`
}

// GetInstanceStatusParams defines parameters for GetInstanceStatus.
type GetInstanceStatusParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`

	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`
}

// GetNetworkParams defines parameters for GetNetwork.
type GetNetworkParams struct {
	// Name the name of a cluster backend network
	Name *string `form:"name,omitempty" json:"name,omitempty"`
}

// GetNetworkIpParams defines parameters for GetNetworkIp.
type GetNetworkIpParams struct {
	// Name the name of a cluster backend network
	Name *string `form:"name,omitempty" json:"name,omitempty"`
}

// GetNodeParams defines parameters for GetNode.
type GetNodeParams struct {
	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`
}

// GetNodeBacklogsParams defines parameters for GetNodeBacklogs.
type GetNodeBacklogsParams struct {
	// Filter list of log filter
	Filter *LogFilter `form:"filter,omitempty" json:"filter,omitempty"`

	// Paths list of object paths to send logs for
	Paths Paths `form:"paths" json:"paths"`
}

// GetNodeConfigParams defines parameters for GetNodeConfig.
type GetNodeConfigParams struct {
	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`
}

// GetNodeDRBDConfigParams defines parameters for GetNodeDRBDConfig.
type GetNodeDRBDConfigParams struct {
	// Name the full path of the file is deduced from the name
	Name DRBDConfigName `form:"name" json:"name"`
}

// PostNodeDRBDConfigParams defines parameters for PostNodeDRBDConfig.
type PostNodeDRBDConfigParams struct {
	// Name the full path of the file is deduced from the name
	Name DRBDConfigName `form:"name" json:"name"`
}

// GetNodeLogsParams defines parameters for GetNodeLogs.
type GetNodeLogsParams struct {
	// Filter list of log filter
	Filter *LogFilter `form:"filter,omitempty" json:"filter,omitempty"`

	// Paths list of object paths to send logs for
	Paths Paths `form:"paths" json:"paths"`
}

// GetNodeMonitorParams defines parameters for GetNodeMonitor.
type GetNodeMonitorParams struct {
	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`
}

// GetNodeStatusParams defines parameters for GetNodeStatus.
type GetNodeStatusParams struct {
	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`
}

// GetObjectParams defines parameters for GetObject.
type GetObjectParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`
}

// GetObjectBacklogsParams defines parameters for GetObjectBacklogs.
type GetObjectBacklogsParams struct {
	// Filter list of log filter
	Filter *LogFilter `form:"filter,omitempty" json:"filter,omitempty"`

	// Paths list of object paths to send logs for
	Paths Paths `form:"paths" json:"paths"`
}

// GetObjectConfigParams defines parameters for GetObjectConfig.
type GetObjectConfigParams struct {
	// Path object path
	Path ObjectPath `form:"path" json:"path"`

	// Evaluate evaluate
	Evaluate *bool `form:"evaluate,omitempty" json:"evaluate,omitempty"`

	// Impersonate impersonate the evaluation as node
	Impersonate *string `form:"impersonate,omitempty" json:"impersonate,omitempty"`
}

// GetObjectFileParams defines parameters for GetObjectFile.
type GetObjectFileParams struct {
	// Path object path
	Path ObjectPath `form:"path" json:"path"`
}

// GetObjectLogsParams defines parameters for GetObjectLogs.
type GetObjectLogsParams struct {
	// Filter list of log filter
	Filter *LogFilter `form:"filter,omitempty" json:"filter,omitempty"`

	// Paths list of object paths to send logs for
	Paths Paths `form:"paths" json:"paths"`
}

// GetObjectPathsParams defines parameters for GetObjectPaths.
type GetObjectPathsParams struct {
	// Path object selector expression.
	Path Path `form:"path" json:"path"`
}

// GetPoolParams defines parameters for GetPool.
type GetPoolParams struct {
	// Name the name of a backend storage pool
	Name *string `form:"name,omitempty" json:"name,omitempty"`
}

// GetPoolVolumeParams defines parameters for GetPoolVolume.
type GetPoolVolumeParams struct {
	// Name the name of a backend storage pool
	Name *string `form:"name,omitempty" json:"name,omitempty"`
}

// GetRelayMessageParams defines parameters for GetRelayMessage.
type GetRelayMessageParams struct {
	// Nodename the nodename component of the slot id on the relay
	Nodename *RelayNodename `form:"nodename,omitempty" json:"nodename,omitempty"`

	// ClusterId the cluster id component of the slot id on the relay
	ClusterId *RelayClusterId `form:"cluster_id,omitempty" json:"cluster_id,omitempty"`
}

// GetResourceParams defines parameters for GetResource.
type GetResourceParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`

	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`

	// Resource resource selector expression.
	Resource *RidOptional `form:"resource,omitempty" json:"resource,omitempty"`
}

// GetResourceConfigParams defines parameters for GetResourceConfig.
type GetResourceConfigParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`

	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`

	// Resource resource selector expression.
	Resource *RidOptional `form:"resource,omitempty" json:"resource,omitempty"`
}

// GetResourceMonitorParams defines parameters for GetResourceMonitor.
type GetResourceMonitorParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`

	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`

	// Resource resource selector expression.
	Resource *RidOptional `form:"resource,omitempty" json:"resource,omitempty"`
}

// GetResourceStatusParams defines parameters for GetResourceStatus.
type GetResourceStatusParams struct {
	// Path object selector expression.
	Path *PathOptional `form:"path,omitempty" json:"path,omitempty"`

	// Node node selector expression.
	Node *NodeOptional `form:"node,omitempty" json:"node,omitempty"`

	// Resource resource selector expression.
	Resource *RidOptional `form:"resource,omitempty" json:"resource,omitempty"`
}

// PostDaemonLogsControlJSONRequestBody defines body for PostDaemonLogsControl for application/json ContentType.
type PostDaemonLogsControlJSONRequestBody = PostDaemonLogsControl

// PostDaemonSubActionJSONRequestBody defines body for PostDaemonSubAction for application/json ContentType.
type PostDaemonSubActionJSONRequestBody = PostDaemonSubAction

// PostInstanceStatusJSONRequestBody defines body for PostInstanceStatus for application/json ContentType.
type PostInstanceStatusJSONRequestBody = InstanceStatusItem

// PostNodeDRBDConfigJSONRequestBody defines body for PostNodeDRBDConfig for application/json ContentType.
type PostNodeDRBDConfigJSONRequestBody = PostNodeDRBDConfigRequest

// PostObjectActionAbortJSONRequestBody defines body for PostObjectActionAbort for application/json ContentType.
type PostObjectActionAbortJSONRequestBody = PostObjectAction

// PostObjectActionDeleteJSONRequestBody defines body for PostObjectActionDelete for application/json ContentType.
type PostObjectActionDeleteJSONRequestBody = PostObjectAction

// PostObjectActionFreezeJSONRequestBody defines body for PostObjectActionFreeze for application/json ContentType.
type PostObjectActionFreezeJSONRequestBody = PostObjectAction

// PostObjectActionGivebackJSONRequestBody defines body for PostObjectActionGiveback for application/json ContentType.
type PostObjectActionGivebackJSONRequestBody = PostObjectAction

// PostObjectActionProvisionJSONRequestBody defines body for PostObjectActionProvision for application/json ContentType.
type PostObjectActionProvisionJSONRequestBody = PostObjectAction

// PostObjectActionPurgeJSONRequestBody defines body for PostObjectActionPurge for application/json ContentType.
type PostObjectActionPurgeJSONRequestBody = PostObjectAction

// PostObjectActionStartJSONRequestBody defines body for PostObjectActionStart for application/json ContentType.
type PostObjectActionStartJSONRequestBody = PostObjectAction

// PostObjectActionStopJSONRequestBody defines body for PostObjectActionStop for application/json ContentType.
type PostObjectActionStopJSONRequestBody = PostObjectAction

// PostObjectActionSwitchJSONRequestBody defines body for PostObjectActionSwitch for application/json ContentType.
type PostObjectActionSwitchJSONRequestBody = PostObjectActionSwitch

// PostObjectActionUnfreezeJSONRequestBody defines body for PostObjectActionUnfreeze for application/json ContentType.
type PostObjectActionUnfreezeJSONRequestBody = PostObjectAction

// PostObjectActionUnprovisionJSONRequestBody defines body for PostObjectActionUnprovision for application/json ContentType.
type PostObjectActionUnprovisionJSONRequestBody = PostObjectAction

// PostObjectClearJSONRequestBody defines body for PostObjectClear for application/json ContentType.
type PostObjectClearJSONRequestBody = PostObjectClear

// PostObjectProgressJSONRequestBody defines body for PostObjectProgress for application/json ContentType.
type PostObjectProgressJSONRequestBody = PostObjectProgress

// PostRelayMessageJSONRequestBody defines body for PostRelayMessage for application/json ContentType.
type PostRelayMessageJSONRequestBody = PostRelayMessage
