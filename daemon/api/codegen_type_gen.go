// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
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

// Defines values for Placement.
const (
	PlacementLastStart  Placement = "last start"
	PlacementLoadAvg    Placement = "load avg"
	PlacementNodesOrder Placement = "nodes order"
	PlacementNone       Placement = "none"
	PlacementScore      Placement = "score"
	PlacementShift      Placement = "shift"
	PlacementSpread     Placement = "spread"
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
	StatusDown      Status = "down"
	StatusNa        Status = "n/a"
	StatusStdbyDown Status = "stdby down"
	StatusStdbyUp   Status = "stdby up"
	StatusUndef     Status = "undef"
	StatusUp        Status = "up"
	StatusWarn      Status = "warn"
)

// Defines values for Topology.
const (
	Failover Topology = "failover"
	Flex     Topology = "flex"
)

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
	ExpireAt time.Time          `json:"expire_at"`
	Id       openapi_types.UUID `json:"id"`
	Minor    int                `json:"minor"`
	Port     int                `json:"port"`
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

// DaemonRunning defines model for DaemonRunning.
type DaemonRunning struct {
	Data []struct {
		Data     bool   `json:"data"`
		Endpoint string `json:"endpoint"`
	} `json:"data"`
	Entrypoint string `json:"entrypoint"`
	Status     int    `json:"status"`
}

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

// InstanceStatus defines model for InstanceStatus.
type InstanceStatus struct {
	App           *string      `json:"app,omitempty"`
	Avail         Status       `json:"avail"`
	Constraints   *bool        `json:"constraints,omitempty"`
	Csum          *string      `json:"csum,omitempty"`
	Drp           *bool        `json:"drp,omitempty"`
	Env           *string      `json:"env,omitempty"`
	FlexMax       *int         `json:"flex_max,omitempty"`
	FlexMin       *int         `json:"flex_min,omitempty"`
	FlexTarget    *int         `json:"flex_target,omitempty"`
	FrozenAt      time.Time    `json:"frozen_at"`
	LastStartedAt time.Time    `json:"last_started_at"`
	Optional      *Status      `json:"optional,omitempty"`
	Orchestrate   *Orchestrate `json:"orchestrate,omitempty"`
	Overall       Status       `json:"overall"`

	// Placement object placement policy
	Placement *Placement `json:"placement,omitempty"`

	// Preserved preserve is true if this status has not been updated due to a
	// heartbeat downtime covered by a maintenance period.
	// when the maintenance period ends, the status should be unchanged,
	// and preserve will be set to false.
	Preserved *bool `json:"preserved,omitempty"`

	// Priority scheduling priority of an object instance on a its node
	Priority *int `json:"priority,omitempty"`

	// Provisioned service, instance or resource provisioned state
	Provisioned Provisioned              `json:"provisioned"`
	Resources   *[]ResourceExposedStatus `json:"resources,omitempty"`
	Running     *[]string                `json:"running,omitempty"`
	Scale       *int                     `json:"scale,omitempty"`
	Slaves      *PathRelation            `json:"slaves,omitempty"`
	StatusGroup *string                  `json:"status_group,omitempty"`

	// Subsets subset properties
	Subsets *map[string]struct {
		Parallel bool `json:"parallel"`
	} `json:"subsets,omitempty"`

	// Topology object topology
	Topology  *Topology `json:"topology,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LogList responseLogList is a list of sse
type LogList = openapi_types.File

// MonitorUpdateQueued defines model for MonitorUpdateQueued.
type MonitorUpdateQueued struct {
	OrchestrationId openapi_types.UUID `json:"orchestration_id"`
}

// NetworkStatus defines model for NetworkStatus.
type NetworkStatus struct {
	Errors  *[]string           `json:"errors,omitempty"`
	Ips     *[]NetworkStatusIp  `json:"ips,omitempty"`
	Name    *string             `json:"name,omitempty"`
	Network *string             `json:"network,omitempty"`
	Type    *string             `json:"type,omitempty"`
	Usage   *NetworkStatusUsage `json:"usage,omitempty"`
}

// NetworkStatusIp defines model for NetworkStatusIp.
type NetworkStatusIp struct {
	Ip   string `json:"ip"`
	Node string `json:"node"`
	Path string `json:"path"`
	Rid  string `json:"rid"`
}

// NetworkStatusList defines model for NetworkStatusList.
type NetworkStatusList = []NetworkStatus

// NetworkStatusUsage defines model for NetworkStatusUsage.
type NetworkStatusUsage struct {
	Free int     `json:"free"`
	Pct  float32 `json:"pct"`
	Size int     `json:"size"`
	Used int     `json:"used"`
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

// NodeLabel defines model for NodeLabel.
type NodeLabel struct {
	// Name name is the label name.
	Name string `json:"name"`

	// Value value is the label value.
	Value string `json:"value"`
}

// NodesInfo defines model for NodesInfo.
type NodesInfo = []NodeInfo

// ObjectConfig defines model for ObjectConfig.
type ObjectConfig struct {
	Data  map[string]interface{} `json:"data"`
	Mtime time.Time              `json:"mtime"`
}

// ObjectFile defines model for ObjectFile.
type ObjectFile struct {
	Data  []byte    `json:"data"`
	Mtime time.Time `json:"mtime"`
}

// ObjectSelection defines model for ObjectSelection.
type ObjectSelection = []string

// Orchestrate defines model for Orchestrate.
type Orchestrate string

// PathRelation defines model for PathRelation.
type PathRelation = []string

// Placement object placement policy
type Placement string

// PoolStatus defines model for PoolStatus.
type PoolStatus struct {
	Capabilities *[]string           `json:"capabilities,omitempty"`
	Errors       *[]string           `json:"errors,omitempty"`
	Head         *string             `json:"head,omitempty"`
	Name         *string             `json:"name,omitempty"`
	Type         *string             `json:"type,omitempty"`
	Usage        *PoolStatusUsage    `json:"usage,omitempty"`
	Volumes      *[]PoolStatusVolume `json:"volumes,omitempty"`
}

// PoolStatusList defines model for PoolStatusList.
type PoolStatusList = []PoolStatus

// PoolStatusUsage defines model for PoolStatusUsage.
type PoolStatusUsage struct {
	Children []string `json:"children"`

	// Orphan an orphan is a volume driven by no svc resource
	Orphan bool   `json:"orphan"`
	Path   string `json:"path"`

	// Size volume size in bytes
	Size float32 `json:"size"`
}

// PoolStatusVolume defines model for PoolStatusVolume.
type PoolStatusVolume struct {
	Ip   string `json:"ip"`
	Node string `json:"node"`
	Path string `json:"path"`
	Rid  string `json:"rid"`
}

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

// PostInstanceStatus defines model for PostInstanceStatus.
type PostInstanceStatus struct {
	Path   string         `json:"path"`
	Status InstanceStatus `json:"status"`
}

// PostNodeDRBDConfigRequest defines model for PostNodeDRBDConfigRequest.
type PostNodeDRBDConfigRequest struct {
	AllocationId openapi_types.UUID `json:"allocation_id"`
	Data         []byte             `json:"data"`
}

// PostNodeMonitor defines model for PostNodeMonitor.
type PostNodeMonitor struct {
	GlobalExpect *string `json:"global_expect,omitempty"`
	LocalExpect  *string `json:"local_expect,omitempty"`
	State        *string `json:"state,omitempty"`
}

// PostObjectAbort defines model for PostObjectAbort.
type PostObjectAbort struct {
	Path string `json:"path"`
}

// PostObjectClear defines model for PostObjectClear.
type PostObjectClear struct {
	Path string `json:"path"`
}

// PostObjectMonitor defines model for PostObjectMonitor.
type PostObjectMonitor struct {
	GlobalExpect *string `json:"global_expect,omitempty"`
	LocalExpect  *string `json:"local_expect,omitempty"`
	Path         string  `json:"path"`
	State        *string `json:"state,omitempty"`
}

// PostObjectProgress defines model for PostObjectProgress.
type PostObjectProgress struct {
	IsPartial *bool              `json:"is_partial,omitempty"`
	Path      string             `json:"path"`
	SessionId openapi_types.UUID `json:"session_id"`
	State     string             `json:"state"`
}

// PostObjectSwitchTo defines model for PostObjectSwitchTo.
type PostObjectSwitchTo struct {
	Destination []string `json:"destination"`
	Path        string   `json:"path"`
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

// ResourceExposedStatus defines model for ResourceExposedStatus.
type ResourceExposedStatus struct {
	// Disable hints the resource ignores all state transition actions
	Disable *bool `json:"disable,omitempty"`

	// Encap indicates that the resource is handled by the encapsulated agents,
	// and ignored at the hypervisor level
	Encap *bool `json:"encap,omitempty"`

	// Info key-value pairs providing interesting information to collect
	// site-wide about this resource
	Info  *map[string]interface{} `json:"info,omitempty"`
	Label string                  `json:"label"`
	Log   *[]struct {
		Level   string `json:"level"`
		Message string `json:"message"`
	} `json:"log,omitempty"`

	// Monitor tells the daemon if it should trigger a monitor action when the
	// resource is not up
	Monitor *bool `json:"monitor,omitempty"`

	// Optional is resource status aggregated into Overall instead of Avail instance status.
	// Errors in optional resource don't stop a state transition action
	Optional    *bool                    `json:"optional,omitempty"`
	Provisioned *ResourceProvisionStatus `json:"provisioned,omitempty"`
	Restart     *int                     `json:"restart,omitempty"`
	Rid         ResourceId               `json:"rid"`

	// Standby resource should always be up, even after a stop state transition action
	Standby *bool  `json:"standby,omitempty"`
	Status  Status `json:"status"`

	// Subset the name of the subset this resource is assigned to
	Subset *string   `json:"subset,omitempty"`
	Tags   *[]string `json:"tags,omitempty"`
	Type   string    `json:"type"`
}

// ResourceId defines model for ResourceId.
type ResourceId = string

// ResourceProvisionStatus defines model for ResourceProvisionStatus.
type ResourceProvisionStatus struct {
	Mtime *time.Time `json:"mtime,omitempty"`

	// State service, instance or resource provisioned state
	State Provisioned `json:"state"`
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

// ObjectPath defines model for ObjectPath.
type ObjectPath = string

// ObjectSelector defines model for ObjectSelector.
type ObjectSelector = string

// Paths defines model for Paths.
type Paths = []string

// RelayClusterId defines model for RelayClusterId.
type RelayClusterId = string

// RelayNodename defines model for RelayNodename.
type RelayNodename = string

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

// GetNetworksParams defines parameters for GetNetworks.
type GetNetworksParams struct {
	// Name the name of a cluster backend network
	Name *string `form:"name,omitempty" json:"name,omitempty"`
}

// GetNodeBacklogsParams defines parameters for GetNodeBacklogs.
type GetNodeBacklogsParams struct {
	// Filter list of log filter
	Filter *LogFilter `form:"filter,omitempty" json:"filter,omitempty"`

	// Paths list of object paths to send logs for
	Paths Paths `form:"paths" json:"paths"`
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

// GetObjectSelectorParams defines parameters for GetObjectSelector.
type GetObjectSelectorParams struct {
	// Selector object selector
	Selector ObjectSelector `form:"selector" json:"selector"`
}

// GetPoolsParams defines parameters for GetPools.
type GetPoolsParams struct {
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

// PostDaemonLogsControlJSONRequestBody defines body for PostDaemonLogsControl for application/json ContentType.
type PostDaemonLogsControlJSONRequestBody = PostDaemonLogsControl

// PostDaemonSubActionJSONRequestBody defines body for PostDaemonSubAction for application/json ContentType.
type PostDaemonSubActionJSONRequestBody = PostDaemonSubAction

// PostInstanceStatusJSONRequestBody defines body for PostInstanceStatus for application/json ContentType.
type PostInstanceStatusJSONRequestBody = PostInstanceStatus

// PostNodeDRBDConfigJSONRequestBody defines body for PostNodeDRBDConfig for application/json ContentType.
type PostNodeDRBDConfigJSONRequestBody = PostNodeDRBDConfigRequest

// PostNodeMonitorJSONRequestBody defines body for PostNodeMonitor for application/json ContentType.
type PostNodeMonitorJSONRequestBody = PostNodeMonitor

// PostObjectAbortJSONRequestBody defines body for PostObjectAbort for application/json ContentType.
type PostObjectAbortJSONRequestBody = PostObjectAbort

// PostObjectClearJSONRequestBody defines body for PostObjectClear for application/json ContentType.
type PostObjectClearJSONRequestBody = PostObjectClear

// PostObjectMonitorJSONRequestBody defines body for PostObjectMonitor for application/json ContentType.
type PostObjectMonitorJSONRequestBody = PostObjectMonitor

// PostObjectProgressJSONRequestBody defines body for PostObjectProgress for application/json ContentType.
type PostObjectProgressJSONRequestBody = PostObjectProgress

// PostObjectSwitchToJSONRequestBody defines body for PostObjectSwitchTo for application/json ContentType.
type PostObjectSwitchToJSONRequestBody = PostObjectSwitchTo

// PostRelayMessageJSONRequestBody defines body for PostRelayMessage for application/json ContentType.
type PostRelayMessageJSONRequestBody = PostRelayMessage
