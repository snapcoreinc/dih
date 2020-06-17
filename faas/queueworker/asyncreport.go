package queueworker

type AsyncReport struct {
	FunctionNameSpace string  `json:"namespace"`
	FunctionName      string  `json:"name"`
	StatusCode        int     `json:"statusCode"`
	TimeTaken         float64 `json:"timeTaken"`
	// TODO: Return URL to fetch status
}
