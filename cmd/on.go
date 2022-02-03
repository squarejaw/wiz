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
	sceneID   int
	speed     int
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
			IP: &ip,
			Params: &bulb.Params{
				Dimming: &dimming,
				Speed:   &speed,
			},
		}
		if cmd.Flag("temp").Changed {
			b.Params.Temp = &temp
		}
		if cmd.Flag("red").Changed {
			b.Params.Red = &red
		}
		if cmd.Flag("green").Changed {
			b.Params.Green = &green
		}
		if cmd.Flag("blue").Changed {
			b.Params.Blue = &blue
		}
		if cmd.Flag("cold-white").Changed {
			b.Params.ColdWhite = &coldWhite
		}
		if cmd.Flag("warm-white").Changed {
			b.Params.WarmWhite = &warmWhite
		}
		if cmd.Flag("scene-id").Changed {
			b.Params.SceneID = &sceneID
		}
		result, err := b.SetState(timeout)
		cobra.CheckErr(err)
		if result != "" {
			fmt.Println(result)
		}
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
	onCmd.Flags().SortFlags = false
	onCmd.Flags().IntVarP(&temp, "temp", "k", 0, "temperature in Kelvin")
	onCmd.Flags().IntVarP(&dimming, "dimming", "d", 50, "dimming")
	onCmd.Flags().IntVarP(&red, "red", "r", 0, "red")
	onCmd.Flags().IntVarP(&green, "green", "g", 0, "green")
	onCmd.Flags().IntVarP(&blue, "blue", "b", 0, "blue")
	onCmd.Flags().IntVarP(&coldWhite, "cold-white", "c", 0, "cold white")
	onCmd.Flags().IntVarP(&warmWhite, "warm-white", "w", 0, "warm white")
	onCmd.Flags().IntVarP(&sceneID, "scene-id", "i", 0, "scene ID")
	onCmd.Flags().IntVarP(&speed, "speed", "s", 50, "speed")
}
