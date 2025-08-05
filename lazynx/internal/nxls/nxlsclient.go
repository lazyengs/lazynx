package nxls

import (
	"context"
	"encoding/json"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/lazyengs/lazynx/internal/config"
	"github.com/lazyengs/lazynx/internal/logs"
	"github.com/lazyengs/lazynx/pkg/nxlsclient"
	"github.com/lazyengs/lazynx/pkg/nxlsclient/commands"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

func CreateNxlsclient(logger *zap.SugaredLogger, config *config.Config) *nxlsclient.Client {
	// Setup separate logger for nxlsclient
	nxlsclientLogFile := filepath.Join(filepath.Dir(config.LogsPath), "nxlsclient.log")
	nxlsclientLogger, err := logs.SetupFileLogger(nxlsclientLogFile, true)
	if err != nil {
		logger.Errorw("Failed to setup nxlsclient logger, using main logger", "error", err)
		nxlsclientLogger = logger
	}

	// Create the client with custom logger but don't initialize it yet
	currentNxWorkspacePath, _ := filepath.Abs("./")
	client := nxlsclient.NewClientWithLogger(currentNxWorkspacePath, true, nxlsclientLogger)
	logger.Infow("Created nxlsclient", "workspacePath", currentNxWorkspacePath)

	return client
}

func InitializeNxlsclient(ctx context.Context, client *nxlsclient.Client, workspacePath string, p *tea.Program, logger *zap.SugaredLogger) error {
	// Update workspace path if different
	if workspacePath != "" {
		absPath, err := filepath.Abs(workspacePath)
		if err != nil {
			return err
		}
		client.NxWorkspacePath = absPath
	}

	// Channel for initialization result
	ch := make(chan *commands.InitializeRequestResult)

	client.OnNotification(commands.RefreshWorkspaceNotificationMethod, func(method string, params json.RawMessage) error {
		logger.Debugw("Received refresh workspace notification", "method", method)
		p.Send(tea.Msg(commands.RefreshWorkspaceNotificationMethod))
		return nil
	})

	// Start client in a goroutine
	go func() {
		logger.Infow("Starting client...")
		params := &protocol.InitializeParams{
			RootURI: protocol.DocumentURI(client.NxWorkspacePath),
			Capabilities: protocol.ClientCapabilities{
				Workspace: &protocol.WorkspaceClientCapabilities{
					Configuration: true,
				},
				TextDocument: &protocol.TextDocumentClientCapabilities{},
			},
			InitializationOptions: map[string]any{
				"workspacePath": client.NxWorkspacePath,
			},
		}
		err := client.Start(ctx, params, ch)
		if err != nil {
			logger.Errorw("Error starting client", "error", err)
		}
	}()

	res := <-ch

	logger.Debugw("Received initialization result", "result", res)
	p.Send(tea.Msg(res))

	return nil
}
