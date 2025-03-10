package nxtypes

type NxVersion struct {
	Full  string `json:"full"`
	Major int    `json:"major"`
	Minor int    `json:"minor"`
}
