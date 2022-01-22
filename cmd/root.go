package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var timeout int

var rootCmd = &cobra.Command{
	Use:   "wiz",
	Short: "WiZ light bulb CLI",
	Long:  `List & control Phillips WiZ light bulbs on your network`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 5, "timeout in seconds to wait for response")
}
