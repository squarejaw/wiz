package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/squarejaw/wiz/bulb"
	"github.com/squarejaw/wiz/udp"
	"github.com/tidwall/gjson"
)

const BUFFER_SIZE = 1024

var (
	addr       string
	outputJSON bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List bulbs",
	Long:  `Sends a UDP broadcast and lists bulbs that respond`,
	Run: func(cmd *cobra.Command, args []string) {
		udpSession, err := udp.NewSession(addr+":38899", time.Duration(timeout)*time.Second)
		cobra.CheckErr(err)
		defer udpSession.Close()

		err = sendRegistration(*udpSession)
		cobra.CheckErr(err)

		if outputJSON {
			cobra.CheckErr(printJSON(*udpSession))
		} else {
			cobra.CheckErr(printPlain(*udpSession))
		}
	},
}

func sendRegistration(udpSession udp.UDPSession) error {
	msg := []byte(`{"method":"registration","params":{"phoneMac":"AAAAAAAAAAAA","register":false,"phoneIp":"1.2.3.4","id":"1"}}`)
	_, err := udpSession.Write(msg)
	return err
}

func printJSON(udpSession udp.UDPSession) error {
	var bulbs []bulb.Bulb
	for {
		buf := make([]byte, BUFFER_SIZE)
		_, addr, err := udpSession.Read(buf)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			b, err := json.Marshal(bulbs)
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		} else if err != nil {
			return err
		}
		udpAddr := addr.(*net.UDPAddr)
		mac := gjson.GetBytes(buf, "result.mac").String()
		if mac != "" {
			ip := udpAddr.IP.String()
			b := bulb.Bulb{IP: &ip, Mac: &mac}
			bulbs = append(bulbs, b)
		}
	}
}

func printPlain(udpSession udp.UDPSession) error {
	fmt.Printf("%-16s%s\n", "IP", "MAC")
	for {
		buf := make([]byte, BUFFER_SIZE)
		_, addr, err := udpSession.Read(buf)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return nil
		} else if err != nil {
			return err
		}
		udpAddr := addr.(*net.UDPAddr)
		mac := gjson.GetBytes(buf, "result.mac").String()
		if mac != "" {
			fmt.Printf("%-16s%s\n", udpAddr.IP, mac)
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&addr, "broadcast-address", "b", "255.255.255.255", "broadcast address")
	listCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "output list as JSON")
}
