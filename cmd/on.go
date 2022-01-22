package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/squarejaw/wiz/udp"
	"github.com/tidwall/gjson"
)

var (
	temp    int
	dimming int
)

var onCmd = &cobra.Command{
	Use:   "on [IP]",
	Short: "Adjust or turn bulb on",
	Long:  `Sends a UDP request to a bulb to turn on with dimming and temperature parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("on needs an IP address of a bulb"))
		}
		ip := args[0]
		udpSession, err := udp.NewSession(ip+":38899", time.Duration(timeout)*time.Second)
		cobra.CheckErr(err)
		defer udpSession.Close()

		msg := []byte(fmt.Sprintf(`{"method":"setState","params":{"temp": %d, "dimming": %d}}`, temp, dimming))
		_, err = udpSession.Write(msg)
		cobra.CheckErr(err)

		buf := make([]byte, 1024)
		_, _, err = udpSession.Read(buf)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return
		}
		cobra.CheckErr(err)
		result := gjson.GetBytes(buf, "result")
		if result.Get("success").Bool() {
			fmt.Println(result)
		} else {
			errorMessage := gjson.GetBytes(buf, "error.message").String()
			fmt.Println("Error: " + errorMessage)
		}
	},
}

func init() {
	rootCmd.AddCommand(onCmd)
	onCmd.Flags().IntVarP(&temp, "kelvin", "k", 2700, "temperature in Kelvin")
	onCmd.Flags().IntVarP(&dimming, "dimming", "d", 50, "dimming")
}
