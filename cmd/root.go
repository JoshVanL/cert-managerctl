package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/joshvanl/cert-managerctl/cmd/options"
)

const defaultLevel = log.InfoLevel

var flags = &options.Flags{}

var rootCmd = &cobra.Command{
	Use:   "cert-managerctl",
	Short: "A tool to interact with cert-manager via the Kubernetes API server.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		level := defaultLevel

		parsed, err := log.ParseLevel(flags.LogLevel)
		if err != nil {
			log.Warnf("Invalid log level '%s', defaulting to '%s'", flags.LogLevel, level)
		} else {
			level = parsed
		}
		log.SetLevel(level)

		return nil
	},
}

func Execute(args []string) {
	mustDie(rootCmd.Execute())
	os.Exit(0)
}

func mustDie(err error) {
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&flags.Kubeconfig,
		"kubeconfig",
		"",
		"Path location to the kubeconfig file. If empty, try in-cluster authentication.",
	)

	rootCmd.PersistentFlags().StringVarP(
		&flags.LogLevel,
		"loglevel",
		"v",
		defaultLevel.String(),
		"logrus log level "+levelsString(),
	)
}

// LevelsString returns a string representing all log levels
// this is useful for help text / flag info
func levelsString() string {
	var b strings.Builder
	b.WriteString("[")
	for i, level := range log.AllLevels {
		b.WriteString(level.String())
		if i+1 != len(log.AllLevels) {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")
	return b.String()
}
