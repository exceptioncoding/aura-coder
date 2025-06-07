package cmd

import (
"fmt"
"os"

"github.com/spf13/cobra"

"aura/app"
)

var rootCmd = &cobra.Command{
Use:   "aura [directory]",
Short: "Aura - AI-powered coding agent for the terminal",
Long: `Aura is a terminal-first, AI-powered coding agent designed to augment 
the developer workflow. It operates as an interactive Terminal User Interface (TUI) 
that integrates conversational AI directly into the developer's primary work environment.`,
Args: cobra.MaximumNArgs(1),
RunE: runAura,
}

func runAura(cmd *cobra.Command, args []string) error {
// Change to specified directory if provided
if len(args) > 0 {
if err := os.Chdir(args[0]); err != nil {
return fmt.Errorf("failed to change to directory %s: %w", args[0], err)
}
}

// Initialize and run the application
application, err := app.NewApp()
if err != nil {
return fmt.Errorf("failed to initialize Aura: %w", err)
}

return application.Run()
}

func Execute() {
if err := rootCmd.Execute(); err != nil {
fmt.Fprintf(os.Stderr, "Error: %v\n", err)
os.Exit(1)
}
}

func init() {
rootCmd.Flags().BoolP("version", "v", false, "Show version information")
}
