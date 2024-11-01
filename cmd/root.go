package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/llmos-ai/llmos/utils/cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	klog "k8s.io/klog/v2"

	"github.com/llmos-ai/llmos-gpu-stack/cmd/devicemanager"
	"github.com/llmos-ai/llmos-gpu-stack/cmd/version"
)

type gpuStack struct {
	Debug      bool   `usage:"Enable debug logging" env:"LLMOS_DEBUG"`
	DebugLevel int    `usage:"Debug log level (valid 0-9) (default 4)" env:"LLMOS_DEBUG_LEVEL"`
	LogFormat  string `usage:"Log format (text or json)" default:"text" short:"l" env:"LLMOS_LOG_FORMAT"`
}

func (g *gpuStack) Run(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func NewRootCmd() *cobra.Command {
	g := &gpuStack{}
	root := cli.Command(g, cobra.Command{
		Use:   "llmos-gpu-stack",
		Short: "LLMOS GPU Stack Management Tool",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	})

	root.AddCommand(
		devicemanager.NewManager(),
		version.NewVersion(),
	)
	root.InitDefaultHelpCmd()
	return root
}

func (l *gpuStack) PersistentPre(_ *cobra.Command, _ []string) error {
	switch l.LogFormat {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	logging := flag.NewFlagSet("", flag.PanicOnError)
	klog.InitFlags(logging)
	if l.Debug || l.DebugLevel > 0 {
		level := l.DebugLevel
		if level == 0 {
			level = 4
		}
		if level > 4 {
			logrus.SetLevel(logrus.TraceLevel)
			logs.Debug = log.New(os.Stderr, "ggcr: ", log.LstdFlags)
		} else {
			logrus.SetLevel(logrus.DebugLevel)
		}
		if err := logging.Parse([]string{
			fmt.Sprintf("-v=%d", level),
		}); err != nil {
			return err
		}
	}

	return nil
}
