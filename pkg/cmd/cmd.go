package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/bwplotka/mdox/pkg/clilog"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/observatorium/obsctl/pkg/version"
	"github.com/spf13/cobra"
)

const (
	logFormatLogfmt = "logfmt"
	logFormatJson   = "json"
	logFormatCLILog = "clilog"
)

var logLevel, logFormat string
var logger log.Logger

func setupLogger(*cobra.Command, []string) {
	var lvl level.Option
	switch logLevel {
	case "error":
		lvl = level.AllowError()
	case "warn":
		lvl = level.AllowWarn()
	case "info":
		lvl = level.AllowInfo()
	case "debug":
		lvl = level.AllowDebug()
	default:
		panic("unexpected log level")
	}
	switch logFormat {
	case logFormatJson:
		logger = level.NewFilter(log.NewJSONLogger(log.NewSyncWriter(os.Stderr)), lvl)
	case logFormatLogfmt:
		logger = level.NewFilter(log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr)), lvl)
	case logFormatCLILog:
		fallthrough
	default:
		logger = level.NewFilter(clilog.New(log.NewSyncWriter(os.Stderr)), lvl)
	}
}

func NewObsctlCmd(ctx context.Context, path ...string) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "obsctl",
		Short:            "CLI to interact with Observatorium",
		Long:             `CLI to interact with Observatorium`,
		Version:          version.Version,
		PersistentPreRun: setupLogger,
	}

	cmd.AddCommand(NewMetricsCmd(ctx, path...))
	cmd.AddCommand(NewContextCommand(ctx, path...))
	cmd.AddCommand(NewLoginCmd(ctx))
	cmd.AddCommand(NewLogoutCmd(ctx))

	cmd.PersistentFlags().StringVar(&logLevel, "log.level", "info", "Log filtering level.")
	cmd.PersistentFlags().StringVar(&logFormat, "log.format", logFormatCLILog, "Log format to use.")

	return cmd
}

// prettyPrintJSON prints indented JSON to stdout.
func prettyPrintJSON(b []byte, w io.Writer) error {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "\t")
	if err != nil {
		return fmt.Errorf("indent JSON %w", err)
	}

	fmt.Fprintln(w, out.String())

	return nil
}
