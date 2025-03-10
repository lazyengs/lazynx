package nxtypes

type PackageJson struct {
	Name                 string            `json:"name"`
	Version              string            `json:"version"`
	License              *string           `json:"license,omitempty"`
	Private              *bool             `json:"private,omitempty"`
	Scripts              map[string]string `json:"scripts,omitempty"`
	Type                 *string           `json:"type,omitempty"`
	Main                 *string           `json:"main,omitempty"`
	Types                *string           `json:"types,omitempty"`
	Typings              *string           `json:"typings,omitempty"`
	Module               *string           `json:"module,omitempty"`
	Exports              interface{}       `json:"exports,omitempty"`
	Dependencies         map[string]string `json:"dependencies,omitempty"`
	DevDependencies      map[string]string `json:"devDependencies,omitempty"`
	OptionalDependencies map[string]string `json:"optionalDependencies,omitempty"`
	PeerDependencies     map[string]string `json:"peerDependencies,omitempty"`
	PeerDependenciesMeta map[string]struct {
		Optional bool `json:"optional"`
	} `json:"peerDependenciesMeta,omitempty"`
	Resolutions map[string]string `json:"resolutions,omitempty"`
	Pnpm        *struct {
		Overrides PackageOverride `json:"overrides,omitempty"`
	} `json:"pnpm,omitempty"`
	Overrides      PackageOverride                    `json:"overrides,omitempty"`
	Bin            interface{}                        `json:"bin,omitempty"`        // can be map[string]string or string
	Workspaces     interface{}                        `json:"workspaces,omitempty"` // can be []string or struct
	PublishConfig  map[string]string                  `json:"publishConfig,omitempty"`
	Files          []string                           `json:"files,omitempty"`
	Nx             *NxProjectPackageJsonConfiguration `json:"nx,omitempty"`
	Generators     *string                            `json:"generators,omitempty"`
	Schematics     *string                            `json:"schematics,omitempty"`
	Builders       *string                            `json:"builders,omitempty"`
	Executors      *string                            `json:"executors,omitempty"`
	NxMigrations   interface{}                        `json:"nx-migrations,omitempty"` // string or NxMigrationsConfiguration
	NgUpdate       interface{}                        `json:"ng-update,omitempty"`     // string or NxMigrationsConfiguration
	PackageManager *string                            `json:"packageManager,omitempty"`
	Description    *string                            `json:"description,omitempty"`
	Keywords       []string                           `json:"keywords,omitempty"`
}

type PackageOverride map[string]interface{} // can be string or PackageOverride

type ArrayPackageGroup []struct {
	Package string `json:"package"`
	Version string `json:"version"`
}

type PackageGroup interface{} // can be MixedPackageGroup or ArrayPackageGroup

type MixedPackageGroup interface{} // can be []interface{} (string or struct) or map[string]string

type NxProjectPackageJsonConfiguration struct {
	ProjectConfiguration
	IncludedScripts []string `json:"includedScripts,omitempty"`
}

type NxMigrationsConfiguration struct {
	Migrations   *string      `json:"migrations,omitempty"`
	PackageGroup PackageGroup `json:"packageGroup,omitempty"`
}
