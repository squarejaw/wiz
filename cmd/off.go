package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/squarejaw/wiz/bulb"
)

var offCmd = &cobra.Command{
	Use:   "off [IP]",
	Short: "Turn bulb off",
	Long:  `Sends a UDP request to set bulb state off`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		state := false
		b := bulb.Bulb{
			IP: &ip,
			Params: &bulb.Params{
				State: &state,
			},
		}
		result, err := b.SetState(timeout)
		cobra.CheckErr(err)
		if result != "" {
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
