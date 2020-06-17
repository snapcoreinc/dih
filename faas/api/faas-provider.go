package api

import (
	"fmt"
	"time"
)

// FunctionDeployment represents a request to create or update a Function.
type FunctionDeployment struct {

	// Service corresponds to a Service
	Service string `json:"service"`

	// Image corresponds to a Docker image
	Image string `json:"image"`

	// Network is specific to Docker Swarm - default overlay network is: func_functions
	Network string `json:"network"`

	// EnvProcess corresponds to the fprocess variable for your container watchdog.
	EnvProcess string `json:"envProcess"`

	// EnvVars provides overrides for functions.
	EnvVars map[string]string `json:"envVars"`

	// RegistryAuth is the registry authentication (optional)
	// in the same encoded format as Docker native credentials
	// (see ~/.docker/config.json)
	RegistryAuth string `json:"registryAuth,omitempty"`

	// Constraints are specific to back-end orchestration platform
	Constraints []string `json:"constraints"`

	// Secrets list of secrets to be made available to function
	Secrets []string `json:"secrets"`

	// Labels are metadata for functions which may be used by the
	// back-end for making scheduling or routing decisions
	Labels *map[string]string `json:"labels"`

	// Annotations are metadata for functions which may be used by the
	// back-end for management, orchestration, events and build tasks
	Annotations *map[string]string `json:"annotations"`

	// Limits for function
	Limits *FunctionResources `json:"limits"`

	// Requests of resources requested by function
	Requests *FunctionResources `json:"requests"`

	// ReadOnlyRootFilesystem removes write-access from the root filesystem
	// mount-point.
	ReadOnlyRootFilesystem bool `json:"readOnlyRootFilesystem"`

	// Namespace for the function to be deployed into
	Namespace string `json:"namespace,omitempty"`

	// Name of the secret holding the image repository credentials
	ImagePullSecret string `json:"imagePullSecret,omitempty"`
}

// FunctionResources Memory and CPU
type FunctionResources struct {
	Memory string `json:"memory"`
	CPU    string `json:"cpu"`
}

// FunctionStatus exported for system/functions endpoint
type FunctionStatus struct {

	// Name corresponds to a Service
	Name string `json:"name"`

	// Image corresponds to a Docker image
	Image string `json:"image"`

	// InvocationCount count of invocations
	InvocationCount float64 `json:"invocationCount"`

	// Replicas desired within the cluster
	Replicas uint64 `json:"replicas"`

	// EnvProcess is the process to pass to the watchdog, if in use
	EnvProcess string `json:"envProcess"`

	// AvailableReplicas is the count of replicas ready to receive
	// invocations as reported by the backend
	AvailableReplicas uint64 `json:"availableReplicas"`

	// Labels are metadata for functions which may be used by the
	// backend for making scheduling or routing decisions
	Labels *map[string]string `json:"labels"`

	// Annotations are metadata for functions which may be used by the
	// backend for management, orchestration, events and build tasks
	Annotations *map[string]string `json:"annotations"`

	// Namespace where the function can be accessed
	Namespace string `json:"namespace,omitempty"`

	// Name of the secret used to pull the image
	ImagePullSecret string `json:"imagePullSecret,omitempty"`
}

// Secret for underlying orchestrator
type Secret struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Value     string `json:"value,omitempty"`
}

type ScaleServiceRequest struct {
	ServiceName string `json:"serviceName"`
	Replicas    uint64 `json:"replicas"`
	//TODO: add memory & cpu, min & max
}

// InfoResponse provides information about the underlying provider
type InfoResponse struct {
	Provider      string          `json:"provider"`
	Version       ProviderVersion `json:"version"`
	Orchestration string          `json:"orchestration"`
}

// ProviderVersion provides the commit sha and release version number of the underlying provider
type ProviderVersion struct {
	SHA     string `json:"sha"`
	Release string `json:"release"`
}

// --- MESSAGE

// Request is the query to return the function logs.
type Request struct {
	// Name is the function name and is required
	Name string `json:"name"`
	// Namespace is the namespace the function is deployed to, how a namespace is defined
	// is faas-provider specific
	Namespace string `json:"namespace"`
	// Instance is the optional container name, that allows you to request logs from a specific function instance
	Instance string `json:"instance"`
	// Since is the optional datetime value to start the logs from
	Since *time.Time `json:"since"`
	// Tail sets the maximum number of log messages to return, <=0 means unlimited
	Tail int `json:"tail"`
	// Follow is allows the user to request a stream of logs until the timeout
	Follow bool `json:"follow"`
}

// String implements that Stringer interface and prints the log Request in a consistent way that
// allows you to safely compare if two requests have the same value.
func (r Request) String() string {
	return fmt.Sprintf(
		"name: %s namespace: %s instance: %s since: %v tail: %d follow: %v",
		r.Name, r.Namespace, r.Instance, r.Since, r.Tail, r.Follow,
	)
}

// Message is a specific log message from a function container log stream
type Message struct {
	// Name is the function name
	Name string `json:"name"`
	// Namespace is the namespace the function is deployed to, how a namespace is defined
	// is faas-provider specific
	Namespace string `json:"namespace"`
	// instance is the name/id of the specific function instance
	Instance string `json:"instance"`
	// Timestamp is the timestamp of when the log message was recorded
	Timestamp time.Time `json:"timestamp"`
	// Text is the raw log message content
	Text string `json:"text"`
}

// String implements the Stringer interface and allows for nice and simple string formatting of a log Message.
func (m Message) String() string {
	ns := ""
	if len(m.Namespace) > 0 {
		ns = fmt.Sprintf("%s ", m.Namespace)
	}
	return fmt.Sprintf(
		"%s %s (%s%s) %s",
		m.Timestamp.String(), m.Name, ns, m.Instance, m.Text,
	)
}
