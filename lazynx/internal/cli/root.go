package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/lazyengs/lazynx/internal/program"
	"github.com/lazyengs/lazynx/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	verbose bool
	logFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lazynx [workspace-path]",
	Short: "A Terminal User Interface for Nx workspace management",
	Long: `LazyNX is a modern TUI application for managing Nx workspaces.

Provide the path to your Nx workspace as an argument, or you will be prompted
to select the nx.json file using an interactive file picker.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLazyNX,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Specify custom log file path")
}

func runLazyNX(cmd *cobra.Command, args []string) error {
	// Setup logging
	var logPath string
	if logFile != "" {
		logPath = logFile
	} else {
		logPath = utils.GetDefaultLogFile()
	}

	logger, err := utils.SetupFileLogger(logPath, verbose)
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
	client := utils.CreateNxlsclient(logger)

	// Create and run the program
	p := program.Create(client, logger, workspacePath)

	// Initialize the nxlsclient
	go func() {
		err := utils.InitializeNxlsclient(cmd.Context(), client, workspacePath, p, logger)
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
	logger.Info("Prompting user to select nx.json file")

	var nxJsonPath string

	picker := huh.NewForm(huh.NewGroup(huh.NewFilePicker().
		Title("Select nx.json file").
		Description("Navigate to and select the nx.json file in your Nx workspace").
		DirAllowed(false).
		FileAllowed(true).
		AllowedTypes([]string{".json"}).
		Value(&nxJsonPath))).WithShowHelp(true)

	err := picker.Run()
	if err != nil {
		logger.Errorw("Error running file picker", "error", err)
		fmt.Fprintf(os.Stderr, "Error selecting nx.json file: %v\n", err)
		os.Exit(1)
	}

	if nxJsonPath == "" {
		logger.Error("No nx.json file selected by user")
		fmt.Fprintf(os.Stderr, "No nx.json file selected\n")
		os.Exit(1)
	}

	fileName := filepath.Base(nxJsonPath)
	logger.Debugw("User selected file", "path", nxJsonPath, "fileName", fileName)

	if fileName != "nx.json" {
		logger.Errorw("Selected file is not nx.json", "fileName", fileName, "expectedFileName", "nx.json")
		fmt.Fprintf(os.Stderr, "Selected file is not nx.json: %s\n", fileName)
		os.Exit(1)
	}

	// Extract the directory path from the nx.json file path
	workspacePath := filepath.Dir(nxJsonPath)
	logger.Infow("Successfully extracted workspace path from nx.json", "workspacePath", workspacePath, "nxJsonPath", nxJsonPath)

	return workspacePath
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
