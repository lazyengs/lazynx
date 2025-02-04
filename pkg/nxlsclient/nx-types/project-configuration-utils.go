package nxtypes

type SourceInformation [2]interface{} // [string|null, string]

type ConfigurationSourceMaps map[string]map[string]SourceInformation
