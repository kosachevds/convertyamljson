package cmd

import (
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:  "convert",
	Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(convertCmd)
}
