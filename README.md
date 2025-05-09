![Logo](https://pub-fbe92a1118f84492b43a7f815ef0dd0a.r2.dev/banner.png)

# LazyNX

> ⚠️ **NOTE:** This is NOT an official Nx project. It's a community-driven tool inspired by the Nx ecosystem.

LazyNX is a project to build a terminal-based UI for [Nx](https://nx.dev) workspaces, inspired by [lazygit](https://github.com/jesseduffield/lazygit) & [lazydocker](https://github.com/jesseduffield/lazydocker). It provides a convenient terminal interface to navigate, view, and run commands in your Nx monorepo without leaving the terminal.

## Project Components

The project consists of two main components:

1. [**lazynx**](./lazynx/README.md) - A terminal UI application (TUI) for Nx workspaces that provides a user-friendly interface to navigate and run Nx commands directly from your terminal.

2. [**nxlsclient**](./pkg/nxlsclient/README.md) - A Go client library for the Nx Language Server Protocol (LSP) server that powers the LazyNX interface. This client listens for changes in the Nx workspace and communicates with the Nx LSP server to provide real-time information.

## Contributing

Contributions are welcome! Please check out our [contribution guidelines](CONTRIBUTING.md) for details on how to get started.

## License

[MIT](LICENSE)
