package nxtypes

type GeneratorCollectionInfo struct {
	Type           string     `json:"type"`
	Name           string     `json:"name"`
	ConfigPath     string     `json:"configPath"` // The path to the file that lists all generators in the collection
	SchemaPath     string     `json:"schemaPath"`
	Data           *Generator `json:"data,omitempty"`
	CollectionName string     `json:"collectionName"`
}

// GeneratorType represents the type of generator
type GeneratorType string

const (
	GeneratorTypeApplication GeneratorType = "application"
	GeneratorTypeLibrary     GeneratorType = "library"
	GeneratorTypeOther       GeneratorType = "other"
)

// Generator represents a generator configuration
type Generator struct {
	Collection  string        `json:"collection"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Options     []Option      `json:"options,omitempty"`
	Type        GeneratorType `json:"type"`
	Aliases     []string      `json:"aliases"`
}

// Option represents a CLI option with additional metadata
type Option struct {
	// Base CLI Option fields
	Name         string `json:"name"`
	OriginalName string `json:"originalName,omitempty"`
	Positional   *int   `json:"positional,omitempty"`
	Alias        string `json:"alias,omitempty"`
	Hidden       bool   `json:"hidden,omitempty"`
	Deprecated   any    `json:"deprecated,omitempty"` // can be bool or string

	// From Schema PropertyDescription (embedded)
	PropertyDescription

	// Additional Option fields
	Tooltip      string            `json:"tooltip,omitempty"`
	ItemTooltips map[string]string `json:"itemTooltips,omitempty"`
	Items        any               `json:"items,omitempty"` // can be []string or ItemsWithEnum
	Aliases      []string          `json:"aliases"`
	IsRequired   bool              `json:"isRequired"`
	XDropdown    string            `json:"x-dropdown,omitempty"`
	XPriority    string            `json:"x-priority,omitempty"` // "important" or "internal"
	XHint        string            `json:"x-hint,omitempty"`
}

// Schema represents a JSON Schema structure
type Schema struct {
	Properties           Properties                     `json:"properties"`
	Required             []string                       `json:"required,omitempty"`
	AnyOf                []Schema                       `json:"anyOf,omitempty"`
	OneOf                []Schema                       `json:"oneOf,omitempty"`
	Description          string                         `json:"description,omitempty"`
	Definitions          Properties                     `json:"definitions,omitempty"`
	AdditionalProperties any                            `json:"additionalProperties,omitempty"` // can be bool or PropertyDescription
	Examples             []SchemaExample                `json:"examples,omitempty"`
	PatternProperties    map[string]PropertyDescription `json:"patternProperties,omitempty"`
}

// SchemaExample represents an example in the Schema
type SchemaExample struct {
	Command     string `json:"command"`
	Description string `json:"description,omitempty"`
}

// Properties is a map of property descriptions
type Properties map[string]PropertyDescription

// PropertyDescription would need to be defined based on the full schema specification
type PropertyDescription struct {
	// Add fields based on your needs
	// This would typically include type information, descriptions, etc.
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	// Add other fields as needed
}

// ItemsWithEnum represents the structure when Items is not a string array
type ItemsWithEnum struct {
	// Add fields based on your needs
	Enum []string `json:"enum,omitempty"`
	// Add other fields as needed
}

type TaskExecutionSchema struct {
	Name           string                `json:"name"`
	Command        string                `json:"command"`
	Collection     string                `json:"collection,omitempty"`
	Positional     string                `json:"positional"`
	Builder        string                `json:"builder,omitempty"`
	Description    string                `json:"description"`
	Configurations []TargetConfiguration `json:"configurations,omitempty"`
	Options        []Option              `json:"options"`
	ContextValues  *TaskContextValues    `json:"contextValues,omitempty"`
}

type TaskContextValues struct {
	Path        string `json:"path,omitempty"`
	Directory   string `json:"directory,omitempty"`
	Project     string `json:"project,omitempty"`
	ProjectName string `json:"projectName,omitempty"`
}
