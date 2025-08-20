package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/lazyengs/lazynx/internal/config"
	"github.com/lazyengs/lazynx/internal/logs"
	"github.com/lazyengs/lazynx/internal/nxls"
	"github.com/lazyengs/lazynx/internal/tui"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lazynx [workspace-path]",
	Short: "A Terminal User Interface for Nx workspace management",
	Long: `LazyNX is a modern TUI application for managing Nx workspaces.

Provide the path to your Nx workspace as an argument, or you will be prompted
to enter the workspace path with validation to ensure it contains nx.json.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLazyNX,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
}

func runLazyNX(cmd *cobra.Command, args []string) error {
	config := config.LoadConfiguration()

	logger, err := logs.SetupFileLogger(config.Logs, verbose)
	if err != nil {
		return fmt.Errorf("error setting up logger: %w", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Starting LazyNX")

	var workspacePath string
	if len(args) > 0 {
		workspacePath = args[0]
		logger.Infow("Using provided workspace path", "path", workspacePath)
	} else {
		// Prompt for workspace path
		workspacePath = promptForWorkspacePath(logger)
		logger.Infow("Using prompted workspace path", "path", workspacePath)
	}

	// Create nxlsclient but don't initialize it yet
	client := nxls.CreateNxlsclient(logger, config)

	// Create and run the program
	p := tui.Create(client, logger, workspacePath)

	// Initialize the nxlsclient
	go func() {
		err := nxls.InitializeNxlsclient(cmd.Context(), client, workspacePath, p, logger)
		if err != nil {
			logger.Errorw("Failed to initialize nxlsclient", "error", err)
		}
	}()

	logger.Info("Starting Bubble Tea program")
	if _, err := p.Run(); err != nil {
		logger.Errorw("Error starting program", "error", err)
		return fmt.Errorf("error starting program: %w", err)
	}

	logger.Info("LazyNX shutting down")
	return nil
}

func promptForWorkspacePath(logger *zap.SugaredLogger) string {
	logger.Info("Prompting user for workspace path")

	var workspacePath string

	input := huh.NewInput().
		Title("Nx Workspace Path").
		Description("Enter the path to your Nx workspace (directory containing nx.json)").
		Placeholder("./").
		Value(&workspacePath).
		Validate(func(s string) error {
			return validateWorkspacePath(s, logger)
		})

	err := input.Run()
	if err != nil {
		logger.Errorw("Error running input prompt", "error", err)
		fmt.Fprintf(os.Stderr, "Error getting workspace path: %v\n", err)
		os.Exit(1)
	}

	if workspacePath == "" {
		workspacePath = "./"
	}

	logger.Infow("User provided workspace path", "workspacePath", workspacePath)
	return workspacePath
}

func validateWorkspacePath(path string, logger *zap.SugaredLogger) error {
	if path == "" {
		path = "./"
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		logger.Debugw("Failed to convert to absolute path", "path", path, "error", err)
		return fmt.Errorf("invalid path: %s", path)
	}

	// Check if directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		logger.Debugw("Path does not exist", "path", absPath, "error", err)
		return fmt.Errorf("path does not exist: %s", path)
	}

	if !info.IsDir() {
		logger.Debugw("Path is not a directory", "path", absPath)
		return fmt.Errorf("path is not a directory: %s", path)
	}

	// Check if nx.json exists in the directory
	nxJsonPath := filepath.Join(absPath, "nx.json")
	if _, err := os.Stat(nxJsonPath); err != nil {
		logger.Debugw("nx.json not found in directory", "path", absPath, "nxJsonPath", nxJsonPath)
		return fmt.Errorf("nx.json not found in directory: %s", path)
	}

	logger.Debugw("Valid Nx workspace found", "workspacePath", absPath, "nxJsonPath", nxJsonPath)
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
