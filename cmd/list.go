package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/squarejaw/wiz/udp"
	"github.com/tidwall/gjson"
)

var addr string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List bulbs",
	Long:  `Sends a UDP broadcast and lists bulbs that respond`,
	Run: func(cmd *cobra.Command, args []string) {
		udpSession, err := udp.NewSession(addr+":38899", time.Duration(timeout)*time.Second)
		cobra.CheckErr(err)
		defer udpSession.Close()

		msg := []byte(`{"method":"registration","params":{"phoneMac":"AAAAAAAAAAAA","register":false,"phoneIp":"1.2.3.4","id":"1"}}`)
		_, err = udpSession.Write(msg)
		cobra.CheckErr(err)

		fmt.Printf("%-16s%s\n", "IP", "MAC")
		for {
			buf := make([]byte, 1024)
			_, addr, err := udpSession.Read(buf)
			if errors.Is(err, os.ErrDeadlineExceeded) {
				break
			}
			cobra.CheckErr(err)
			udpAddr := addr.(*net.UDPAddr)
			mac := gjson.GetBytes(buf, "result.mac")
			if mac.String() != "" {
				fmt.Printf("%-16s%s\n", udpAddr.IP, mac)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&addr, "broadcast-address", "b", "255.255.255.255", "broadcast address")
}
