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

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "Get current bulb state",
	Long:  `Sends a UDP request to get the current state of the bulb and prints it to standard out`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		udpSession, err := udp.NewSession(ip, time.Duration(timeout)*time.Second)
		cobra.CheckErr(err)
		defer udpSession.Close()

		msg := []byte(`{"method":"getPilot","params":{}}`)
		_, err = udpSession.Write(msg)
		cobra.CheckErr(err)

		buf := make([]byte, udp.MaxSafePayloadSize)
		_, _, err = udpSession.Read(buf)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			cobra.CheckErr(fmt.Errorf("no response from bulb %s within %ds", ip, timeout))
		}
		cobra.CheckErr(err)
		result := gjson.GetBytes(buf, "result")
		fmt.Println(result)
	},
}

func init() {
	rootCmd.AddCommand(stateCmd)
}
