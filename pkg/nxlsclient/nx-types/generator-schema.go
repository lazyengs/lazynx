package nxtypes

type GeneratorContext struct {
	Project             string            `json:"project,omitempty"`
	NormalizedDirectory string            `json:"normalizedDirectory,omitempty"`
	Directory           string            `json:"directory,omitempty"`
	PrefillValues       map[string]string `json:"prefillValues,omitempty"`
	FixedFormValues     map[string]string `json:"fixedFormValues,omitempty"`
	NxVersion           *NxVersion        `json:"nxVersion,omitempty"`
}

type GeneratorSchema struct {
	CollectionName string            `json:"collectionName"`
	GeneratorName  string            `json:"generatorName"`
	Description    string            `json:"description"`
	Options        []Option          `json:"options"`
	Context        *GeneratorContext `json:"context,omitempty"`
}
