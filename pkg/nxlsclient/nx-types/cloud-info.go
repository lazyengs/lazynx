package nxtypes

type CloudOnboardingInfo struct {
	HasNxInCI               bool   `json:"hasNxInCI"`
	HasAffectedCommandsInCI bool   `json:"hasAffectedCommandsInCI"`
	IsConnectedToCloud      bool   `json:"isConnectedToCloud"`
	IsWorkspaceClaimed      bool   `json:"isWorkspaceClaimed"`
	PersonalAccessToken     string `json:"personalAccessToken,omitempty"`
}

type CIPEExecutionStatus string

const (
	CIPEStatusNotStarted CIPEExecutionStatus = "NOT_STARTED"
	CIPEStatusInProgress CIPEExecutionStatus = "IN_PROGRESS"
	CIPEStatusSucceeded  CIPEExecutionStatus = "SUCCEEDED"
	CIPEStatusFailed     CIPEExecutionStatus = "FAILED"
	CIPEStatusCanceled   CIPEExecutionStatus = "CANCELED"
	CIPEStatusTimedOut   CIPEExecutionStatus = "TIMED_OUT"
)

type CIPEInfo struct {
	CiPipelineExecutionId string              `json:"ciPipelineExecutionId"`
	Branch                string              `json:"branch"`
	Status                CIPEExecutionStatus `json:"status"`
	CreatedAt             int64               `json:"createdAt"`
	CompletedAt           *int64              `json:"completedAt"`
	CommitTitle           *string             `json:"commitTitle"`
	CommitUrl             *string             `json:"commitUrl"`
	Author                *string             `json:"author,omitempty"`
	AuthorAvatarUrl       *string             `json:"authorAvatarUrl,omitempty"`
	CipeUrl               string              `json:"cipeUrl"`
	RunGroups             []CIPERunGroup      `json:"runGroups"`
}

type CIPERunGroup struct {
	CiExecutionEnv string              `json:"ciExecutionEnv"`
	RunGroup       string              `json:"runGroup"`
	CreatedAt      int64               `json:"createdAt"`
	CompletedAt    *int64              `json:"completedAt"`
	Status         CIPEExecutionStatus `json:"status"`
	Runs           []CIPERun           `json:"runs"`
}

type CIPERun struct {
	LinkId         string               `json:"linkId"`
	Command        string               `json:"command"`
	Status         *CIPEExecutionStatus `json:"status,omitempty"`
	NumFailedTasks *int                 `json:"numFailedTasks,omitempty"`
	NumTasks       *int                 `json:"numTasks,omitempty"`
	RunUrl         string               `json:"runUrl"`
}

type CIPEErrorType string

const (
	CIPEErrorTypeAuthentication CIPEErrorType = "authentication"
	CIPEErrorTypeNetwork        CIPEErrorType = "network"
	CIPEErrorTypeOther          CIPEErrorType = "other"
)

type CIPEInfoError struct {
	Message string        `json:"message"`
	Type    CIPEErrorType `json:"type"`
}
