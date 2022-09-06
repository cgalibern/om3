package env

import "os"

var (
	ActionOriginVar    = "OSVC_ACTION_ORIGIN"
	ActionOriginUser   = "user"
	ActionOriginDaemon = "daemon"

	NameVar      = "OSVC_NAME"
	NamespaceVar = "OSVC_NAMESPACE"
	KindVar      = "OSVC_KIND"
	ContextVar   = "OSVC_CONTEXT"
)

// HasDaemonOrigin returns true if the environment variable OSVC_ACTION_ORIGIN
// is set to "daemon". The opensvc daemon sets this variable on every command
// it executes.
func HasDaemonOrigin() bool {
	return os.Getenv(ActionOriginVar) == ActionOriginDaemon
}

// Origin returns the action origin using a env var that the daemon sets when
// executing a CRM action. The only possible return values are "daemon" or "user".
func Origin() string {
	s := os.Getenv(ActionOriginVar)
	if s == "" {
		s = ActionOriginUser
	}
	return s
}

// DaemonOriginSetenvArgs returns the args to pass to environment variable
// setter functions to hint the called CRM command was launched from a daemon
// policy.
func DaemonOriginSetenvArgs() []string {
	return []string{ActionOriginVar + "=" + ActionOriginDaemon}
}

// Namespace returns the namespace filter forced via the OSVC_NAMESPACE environment
// variable.
func Namespace() string {
	return os.Getenv(NamespaceVar)
}

// Kind returns the object kind filter forced via the OSVC_NAMESPACE environment
// variable.
func Kind() string {
	return os.Getenv(KindVar)
}

// Context returns the identifier of a remote cluster endpoint and credentials
// configuration via the OSVC_CONTEXT variable.
func Context() string {
	return os.Getenv(ContextVar)
}
