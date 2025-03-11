# Contributing to LazyNX

We love your input! We want to make contributing to LazyNX as easy and transparent as possible, whether it's:

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
- Node.js 16 or later
- Nx CLI

### Project Structure

This is a monorepo managed with Nx. It consists of these main components:

```
lazynx/
├── lazynx/             # Main TUI application
│   ├── main.go         # Application entry point
│   └── internal/       # Internal application code
├── pkg/                # Package libraries
│   └── nxlsclient/     # Nx LSP client library
└── internal/           # Internal utilities and examples
```

## Setup Development Environment

1. Clone the repository

   ```bash
   git clone https://github.com/lazyengs/lazynx.git
   cd lazynx
   ```

2. Install dependencies

   ```bash
   pnpm install
   ```

3. Run tests

   ```bash
   nx run-many -t test
   ```

## Making Changes

### Running the Application Locally

```bash
nx run lazynx:serve
```

### Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Run `go fmt` before committing
- Document all exported functions, types, and variables

## Testing

Each package has its own tests that can be run with:

```bash
nx run <package>:test
```

## Reporting Bugs

We use GitHub issues to track bugs. Report a bug by [opening a new issue](https://github.com/lazyengs/lazynx/issues/new); it's that easy!

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
