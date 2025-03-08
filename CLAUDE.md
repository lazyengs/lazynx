# LazyNX Development Guide

## Build/Run/Test Commands
- Run all tests: `nx run-many -t test`
- Run single test: `nx run <project>:test -- -run <TestName>` (e.g. `nx run nxlsclient:test -- -run TestClientE2E`)
- Run all linters: `nx run-many -t lint`
- Lint single project: `nx run <project>:lint`
- Build application: `nx run lazynx:build`
- Serve application: `nx run lazynx:serve`
- Update Go dependencies: `nx run <project>:tidy`

## Code Style Guidelines
- **Formatting**: Go standard formatting (`gofmt`)
- **Imports**: Standard Go import organization (stdlib, external, internal)
- **Error Handling**: Check errors using `require` for tests, proper error propagation in application code
- **Naming**: Follow Go conventions (CamelCase for exported, camelCase for unexported)
- **Testing**: Use testify for assertions/requirements, cupaloy for snapshot testing
- **Structure**: Organize code by feature/domain in `/internal` for private code, `/pkg` for reusable packages
- **Logging**: Use zap for structured logging (see `pkg/nxlsclient/client.go`)