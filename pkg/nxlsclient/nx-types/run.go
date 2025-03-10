package nxtypes

type Target struct {
	Project       string `json:"project"`
	Target        string `json:"target"`
	Configuration string `json:"configuration,omitempty"`
}
