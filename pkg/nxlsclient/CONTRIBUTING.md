# Contributing to nxlsclient

We love your input! We want to make contributing to `nxlsclient` as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code lints
6. Issue that pull request!

### Prerequisites

- Go 1.19 or later
- Node.js 16 or later (for running the nxls server locally if needed)
- An Nx workspace for testing

### Project Structure

```
nxlsclient/
├── client.go           # Main client implementation
├── commands/           # LSP commands implementation directory
│   ├── commands.go     # Base commander implementation
│   └── [command].go    # Individual command implementations
├── examples/           # Example implementations
├── listener.go         # Notification listener implementation
├── notifications.go    # Notification type definitions and utilities
├── nx-types/           # Nx-specific type definitions
├── rwc.go              # ReadWriteCloser interface implementation
├── server.go           # Server management functions
├── stream.go           # Stream handling for JSON-RPC
└── server/             # Embedded nxls server files
    └── nxls/           # Node.js LSP server code
```

## Setup Development Environment

1. Clone the repository

   ```bash
   git clone https://github.com/lazyengs/pkg
   cd pkg/nxlsclient
   ```

2. Install dependencies

   ```bash
   go mod download
   ```

3. Run tests
   ```bash
   go test ./...
   ```

## Making Changes

### Local Testing

To test your changes with a real Nx workspace:

1. Create or use an existing Nx workspace
2. Use the playground example in `internal/playground/main.go` as a reference
3. Update the workspace path to point to your test workspace
4. Run the example:
   ```bash
   go run internal/playground/main.go
   ```

### Adding a New Command

1. Create a new file in the `commands/` directory named after your command (e.g., `my-feature.request.go`)
2. Define the request method constant, parameters, and result types
3. Implement the `SendMyFeatureRequest` method on the `Commander` type
4. Add tests for your command in a corresponding test file
5. Update documentation to include your new command

Example:

```go
package commands

import (
    "context"
)

const (
    MyFeatureRequestMethod = "nx/myFeature"
)

type MyFeatureParams struct {
    Option1 string `json:"option1"`
    Option2 bool   `json:"option2"`
}

type MyFeatureResult struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func (c *Commander) SendMyFeatureRequest(ctx context.Context, params MyFeatureParams) (*MyFeatureResult, error) {
    var result *MyFeatureResult
    err := c.sendRequest(ctx, MyFeatureRequestMethod, params, &result)
    return result, err
}
```

### Adding a New Notification Handler

1. Define the notification method constant in `notifications.go`
2. Create a typed handler struct if needed
3. Use the notification in your code or examples

### Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt` before committing
- Document all exported functions, types, and variables
- Add comprehensive examples for new functionality

## Testing

- Write unit tests for all functionality
- Use End-to-End (E2E) tests for full client/server integration testing
- Use mocks when appropriate to isolate components

## Reporting Bugs

We use GitHub issues to track bugs. Report a bug by [opening a new issue](https://github.com/lazyengs/pkg/issues/new); it's that easy!

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## License

By contributing, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
