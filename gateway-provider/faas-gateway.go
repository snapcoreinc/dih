package gateway_provider

// DeleteFunctionRequest delete a deployed function
type DeleteFunctionRequest struct {
	FunctionName string `json:"functionName"`
}

//FunctionDescription information related to a function
type FunctionDescription struct {
	Name              string
	Status            string
	Replicas          int
	AvailableReplicas int
	InvocationCount   int
	Image             string
	EnvProcess        string
	URL               string
	AsyncURL          string
	Labels            *map[string]string
	Annotations       *map[string]string
	ImagePullSecret   string
}
