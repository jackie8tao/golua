package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golua",
	Short: "A Lua interpreter written in Go",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
