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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		b := bulb.Bulb{
			IP:     &ip,
			Params: &bulb.Params{},
		}
		// Only forward flags the user actually set; the pointer-based
		// Params struct distinguishes "unset" from a real zero value
		// (e.g. red=0 is a valid color channel).
		for _, f := range []struct {
			name  string
			field **int
			src   *int
		}{
			{"temp", &b.Params.Temp, &temp},
			{"dimming", &b.Params.Dimming, &dimming},
			{"red", &b.Params.Red, &red},
			{"green", &b.Params.Green, &green},
			{"blue", &b.Params.Blue, &blue},
			{"cold-white", &b.Params.ColdWhite, &coldWhite},
			{"warm-white", &b.Params.WarmWhite, &warmWhite},
			{"scene-id", &b.Params.SceneID, &sceneID},
			{"speed", &b.Params.Speed, &speed},
		} {
			if cmd.Flag(f.name).Changed {
				*f.field = f.src
			}
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
