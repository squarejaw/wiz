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

var offCmd = &cobra.Command{
	Use:   "off [IP]",
	Short: "Turn bulb off",
	Long:  `Sends a UDP request to a bulb to turn off`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("on needs an IP address of a bulb"))
		}
		ip := args[0]
		udpSession, err := udp.NewSession(ip+":38899", time.Duration(timeout)*time.Second)
		cobra.CheckErr(err)
		defer udpSession.Close()

		msg := []byte(`{"method":"setPilot","params":{"state": "false"}}`)
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
	rootCmd.AddCommand(offCmd)
}
