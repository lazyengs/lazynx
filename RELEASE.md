# Release Strategy for LazyNX

This document outlines our release strategy for packages in this repository.

## Release Philosophy

- We use [Nx Release](https://nx.dev/recipes/nx-release/get-started-with-nx-release) to manage our releases
- All packages are versioned independently
- We follow [Semantic Versioning](https://semver.org/)
- We use [Conventional Commits](https://www.conventionalcommits.org/) to automatically determine version numbers

## Projects Setup

- **nxlsclient** - Go library for communicating with the Nx Language Server
- **lazynx** - Terminal UI application for managing Nx workspaces

## How Releases Work

### Automated Release Process

Releases are automatically triggered when commits are pushed to the `main` branch:

1. A GitHub Action checks for new changes in the projects
2. The action runs tests and linters to ensure quality
3. If tests pass, `nx release` is executed to:
   - Determine the next version based on commit messages
   - Update version references
   - Generate a changelog
   - Create Git tags with the format `{projectName}-v{version}`
   - Create a GitHub release
4. For the nxlsclient library, the module is published to pkg.go.dev

### Manual Releases

To create a release manually:

1. Ensure your changes are committed and pushed to the main branch
2. Run `nx release --projects=nxlsclient` (or `lazynx`) to create a release for a specific project
3. Specify the version bump type when prompted (patch, minor, major)

## Viewing Release History

- GitHub Releases: See the [Releases page](https://github.com/lazyengs/lazynx/releases)
- Changelogs: Each project maintains its own CHANGELOG.md file

## Release Tagging Strategy

Tags follow the format: `{projectName}-v{version}`

Examples:
- `nxlsclient-v0.1.0`
- `lazynx-v0.2.3`

## For Go Module Users

The nxlsclient library can be imported as:

```go
import "github.com/lazyengs/lazynx/pkg/nxlsclient"
```

You can use a specific version with:

```bash
go get github.com/lazyengs/lazynx/pkg/nxlsclient@v0.1.0
```

## Troubleshooting

If you encounter issues with the release process:

1. Check the GitHub Actions logs for errors
2. Ensure all tests are passing
3. Verify that your commit messages follow the Conventional Commits format
4. Check that your go.mod file has the correct module path