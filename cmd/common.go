package cmd

// Response represences response from OreCast service
type Response struct {
	Status string `json:"status"`
	Error  any    `json:"error,omitempty"`
}
