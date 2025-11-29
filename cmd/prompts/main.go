package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"prompts-mcp/internal/logging"
	"prompts-mcp/internal/mcp"
	"prompts-mcp/internal/prompts"
	"prompts-mcp/internal/validate"
)

func main() {
	promptsDirFlag := flag.String("prompts-dir", "", "absolute path to prompts directory containing YAML prompt definitions")
	flag.Parse()

	logger, err := logging.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() //nolint:errcheck


	promptsDir, err := validate.Dir(*promptsDirFlag)
	if err != nil {
		logger.Fatal("invalid prompts-dir", zap.Error(err))
	}

	promptsRepo := prompts.NewYAMLRepository(promptsDir)


	server := mcp.NewServer(logger, promptsRepo)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := server.Serve(ctx, os.Stdin, os.Stdout); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
