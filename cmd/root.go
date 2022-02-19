package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var defaultCmd = convertCmd

var rootCmd = &cobra.Command{
	Use:          "convertyamljson",
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	insertIfNoCommand(rootCmd, defaultCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func insertIfNoCommand(root, defaultCmd *cobra.Command) {
	defaultCmd.Short += " (call by default)"
	cmd, _, err := root.Find(os.Args[1:])
	if err != nil && cmd == root && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		newArgs := append([]string{defaultCmd.Name()}, os.Args[1:]...)
		root.SetArgs(newArgs)
	}
}
