package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/szktkfm/clipy"
)

var (
	rootCmd = &cobra.Command{
		Use:   "clipy",
		Short: "Clipboard history manager and retriever",
		Long:  `clipy is a CLI tool for managing clipboard history and easily retrieving past clipboard entries.`,
		Run:   runDaemon,
	}

	startDaemonCmd = &cobra.Command{
		Use:   "start-daemon",
		Short: "Start clipboard history manager as a daemon",
		Long:  `Start the clipboard history manager as a background daemon process.`,
		Run:   startDaemon,
	}

	generateShellCmd = &cobra.Command{
		Use:   "shell",
		Short: "Generate a shell integration script",
		Long:  `Generate a shell script that integrates clipy into your terminal profile.`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.PrintErrln("Error: Shell type is required (e.g., bash, zsh).")
				return
			}
			clipy.GenerateShellScript(args[0])
		},
	}

	listHistoryCmd = &cobra.Command{
		Use:   "history",
		Short: "List clipboard history",
		Long:  `Dump the clipboard history, showing all saved entries.`,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := resolveConfigDir()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			tmux, err := cmd.Flags().GetBool("tmux")
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			if tmux {
				output, err := runTmuxCmd()
				var tmuxErr *TmuxError
				if err != nil {
					if errors.As(err, &tmuxErr) {
						err = clipy.ListHistories(config)
					}
					if err != nil {
						cmd.PrintErrln(err)
					}
					return
				} else {
					fmt.Print(string(output))
					return
				}
			}

			if err := clipy.ListHistories(config); err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}
)

type TmuxError struct {
	Message string
}

func (e *TmuxError) Error() string {
	return fmt.Sprintf("failed to execute tmux command: %s", e.Message)
}

func runTmuxCmd() (string, error) {
	tmpfile, tmpfileName, err := mkTmpFile("clipy-")
	if err != nil {
		return "", err
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfileName)
	}()
	tmuxCmd := exec.Command("sh", "-c", fmt.Sprintf("tmux popup -E sh -c \"clipy history | tee %v\"", tmpfileName))
	err = tmuxCmd.Run()
	if err != nil {
		return "", &TmuxError{Message: err.Error()}
	}
	output, err := io.ReadAll(tmpfile)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func mkTmpFile(name string) (*os.File, string, error) {
	ns := time.Now().UnixNano()
	output := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%d", name, ns))
	f, err := os.OpenFile(output, os.O_CREATE, 0700)
	if err != nil {
		return nil, "", err
	}

	return f, output, nil
}

func runDaemon(cmd *cobra.Command, _ []string) {
	// Kill any running clipy instances
	err := clipy.StopAllInstances()
	if err != nil {
		cmd.PrintErrln(err)
	}
	daemonProcess := exec.Command(os.Args[0], "start-daemon")
	fmt.Printf("Starting clipboard history daemon...\n")
	_ = daemonProcess.Start()
}

func startDaemon(cmd *cobra.Command, _ []string) {
	config, err := resolveConfigDir()
	if err != nil {
		cmd.PrintErrln(err)
	}

	err = clipy.Execute(config)
	if err != nil {
		cmd.PrintErrln(err)
	}
}

func resolveConfigDir() (string, error) {
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return xdgConfig, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, ".config/clipy")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "clipy.db")

	file, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open or create config file: %w", err)
	}
	defer file.Close()

	return configFile, nil
}

func init() {
	rootCmd.Flags().Bool(
		"debug",
		false,
		"Enable debug logging to debug.log",
	)

	rootCmd.Flags().BoolP(
		"help",
		"h",
		false,
		"Show help for clipy",
	)

	listHistoryCmd.Flags().Bool(
		"tmux",
		false,
		"list histories in tmux popup",
	)

	rootCmd.AddCommand(listHistoryCmd)
	rootCmd.AddCommand(startDaemonCmd)
	rootCmd.AddCommand(generateShellCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
