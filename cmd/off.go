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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("on needs an IP address of a bulb"))
		}
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
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(offCmd)
}
