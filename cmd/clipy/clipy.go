package main

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/szktkfm/clipy"
)

var (
	rootCmd = &cobra.Command{
		Use:   "clipy",
		Short: "Clipboard history daemon",
		Run:   run,
	}

	childCmd = &cobra.Command{
		Use:   "child",
		Short: "Clipboard history daemon",
		Run:   runChild,
	}
	shellCmd = &cobra.Command{
		Use:   "shell",
		Short: "shell",
		Run: func(cmd *cobra.Command, args []string) {
			clipy.Shell(args[0])
		},
	}
	historyCmd = &cobra.Command{
		Use:   "history",
		Short: "Dump your clipboard history",
		Run: func(cmd *cobra.Command, args []string) {
			clipy.ListHistories()
		},
	}
)

func run(cmd *cobra.Command, _ []string) {
	clipy.KillClipies()
	childCmd := exec.Command(os.Args[0], "child")
	childCmd.Start()
}

func runChild(cmd *cobra.Command, _ []string) {
	clipy.Execute()
}

func init() {
	rootCmd.Flags().Bool(
		"debug",
		false,
		"passing this flag will allow writing debug output to debug.log",
	)

	rootCmd.Flags().BoolP(
		"help",
		"h",
		false,
		"help for clipy",
	)

	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(childCmd)
	rootCmd.AddCommand(shellCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
