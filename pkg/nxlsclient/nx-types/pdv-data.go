package nxtypes

type PDVResultType string

const (
	PDVResultTypeNoGraphError PDVResultType = "NO_GRAPH_ERROR"
	PDVResultTypeOldNxVersion PDVResultType = "OLD_NX_VERSION"
	PDVResultTypeError        PDVResultType = "ERROR"
	PDVResultTypeSuccess      PDVResultType = "SUCCESS"
	PDVResultTypeSuccessMulti PDVResultType = "SUCCESS_MULTI"
)

type PDVData struct {
	ResultType             PDVResultType     `json:"resultType"`
	GraphBasePath          string            `json:"graphBasePath,omitempty"`
	PDVDataSerialized      string            `json:"pdvDataSerialized,omitempty"`
	PDVDataSerializedMulti map[string]string `json:"pdvDataSerializedMulti,omitempty"`
	ErrorsSerialized       string            `json:"errorsSerialized,omitempty"`
	ErrorMessage           string            `json:"errorMessage,omitempty"`
}
