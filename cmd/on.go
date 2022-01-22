package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/squarejaw/wiz/bulb"
)

var (
	temp      int
	dimming   int
	red       int
	green     int
	blue      int
	coldWhite int
	warmWhite int
)

var onCmd = &cobra.Command{
	Use:   "on [IP]",
	Short: "Turn bulb on",
	Long:  `Sends a UDP request to set bulb state on with parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("on needs an IP address of a bulb"))
		}
		ip := args[0]
		b := bulb.Bulb{
			IP: ip,
			Params: bulb.Params{
				Dimming: &dimming,
			},
		}
		if temp >= 1000 {
			b.Params.Temp = &temp
		}
		if red > 0 {
			b.Params.Red = &red
		}
		if green > 0 {
			b.Params.Green = &green
		}
		if blue > 0 {
			b.Params.Blue = &blue
		}
		if coldWhite > 0 {
			b.Params.ColdWhite = &coldWhite
		}
		if warmWhite > 0 {
			b.Params.WarmWhite = &warmWhite
		}
		result, err := b.SetState(timeout)
		cobra.CheckErr(err)
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
	onCmd.Flags().IntVarP(&temp, "kelvin", "k", 0, "temperature in Kelvin")
	onCmd.Flags().IntVarP(&dimming, "dimming", "d", 50, "dimming")
	onCmd.Flags().IntVarP(&red, "red", "r", 0, "red")
	onCmd.Flags().IntVarP(&green, "green", "g", 0, "green")
	onCmd.Flags().IntVarP(&blue, "blue", "b", 0, "blue")
	onCmd.Flags().IntVarP(&coldWhite, "cold-white", "c", 0, "cold white")
	onCmd.Flags().IntVarP(&warmWhite, "warm-white", "w", 0, "warm white")
}
