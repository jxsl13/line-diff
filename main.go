package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/jxsl13/line-diff/config"
	"github.com/spf13/cobra"
)

func main() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func NewRootCmd() *cobra.Command {
	rootContext := rootContext{}

	// rootCmd represents the run command
	rootCmd := &cobra.Command{
		Use:   "line-diff a.txt b.txt",
		Short: "diff two text files",
		Args:  cobra.ExactArgs(2),
		RunE:  rootContext.RunE,
	}

	// register flags but defer parsing and validation of the final values
	rootCmd.PreRunE = rootContext.PreRunE(rootCmd)

	rootCmd.AddCommand(NewCompletionCmd(rootCmd.Name()))

	return rootCmd
}

type rootContext struct {
	Config     *config.Config
	SourcePath string
	TargetPath string
}

func (c *rootContext) PreRunE(cmd *cobra.Command) func(cmd *cobra.Command, args []string) error {
	c.Config = &config.Config{
		Sorted: false,
	}

	runParser := config.RegisterFlags(c.Config, true, cmd)

	return func(cmd *cobra.Command, args []string) error {
		for idx, a := range args {
			switch idx {
			case 0:
				c.SourcePath = a
			case 1:
				c.TargetPath = a
			}
		}

		return runParser()
	}
}

func (c *rootContext) RunE(cmd *cobra.Command, args []string) (err error) {
	source, target := make(map[string]bool, 2048), make(map[string]bool, 2048)

	added, removed, unchanged := make([]string, 0, 2048), make([]string, 0, 2048), make([]string, 0, 2048)
	b, err := os.ReadFile(c.SourcePath)
	if err != nil {
		return err
	}

	sourceLines := strings.Split(string(b), "\n")
	for idx, line := range sourceLines {
		line = strings.TrimSpace(line)
		sourceLines[idx] = line
		source[line] = true
	}

	b, err = os.ReadFile(c.TargetPath)
	if err != nil {
		return err
	}

	targetLines := strings.Split(string(b), "\n")
	for _, line := range targetLines {
		line = strings.TrimSpace(line)
		target[line] = true

		if source[line] {
			unchanged = append(unchanged, line)
		} else {
			added = append(added, line)
		}
	}

	for line := range source {
		if !target[line] {
			removed = append(removed, line)
		}
	}

	if c.Config.Sorted {
		slices.Sort(added)
		slices.Sort(removed)
		slices.Sort(unchanged)
	}

	if len(added) > 0 {
		fmt.Printf("--- added lines (%s -> %s) ---\n", c.SourcePath, c.TargetPath)
		for _, k := range added {
			fmt.Println(k)
		}
	}

	if len(removed) > 0 {
		fmt.Printf("--- removed lines (%s -> %s) ---\n", c.SourcePath, c.TargetPath)
		for _, k := range removed {
			fmt.Println(k)
		}
	}

	if len(unchanged) > 0 {
		fmt.Printf("--- unchanged lines (%s -> %s) ---\n", c.SourcePath, c.TargetPath)
		for _, k := range unchanged {
			fmt.Println(k)
		}
	}

	return nil
}
