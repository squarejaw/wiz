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

var (
	addr       string
	outputJSON bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List bulbs",
	Long:  `Sends a UDP broadcast and lists bulbs that respond`,
	Run: func(cmd *cobra.Command, args []string) {
		udpSession, err := udp.NewSession(addr, time.Duration(timeout)*time.Second)
		cobra.CheckErr(err)
		defer udpSession.Close()

		err = sendRegistration(udpSession)
		cobra.CheckErr(err)

		bulbs, err := collectBulbs(udpSession)
		cobra.CheckErr(err)
		if outputJSON {
			printJSON(bulbs)
		} else {
			printPlain(bulbs)
		}
	},
}

func sendRegistration(udpSession *udp.UDPSession) error {
	msg := []byte(`{"method":"registration","params":{"phoneMac":"AAAAAAAAAAAA","register":false,"phoneIp":"1.2.3.4","id":"1"}}`)
	_, err := udpSession.Write(msg)
	return err
}

// packetReader is the subset of *udp.UDPSession used to read bulb responses.
// Defining it as an interface makes the collection logic unit-testable.
type packetReader interface {
	Read(buf []byte) (int, net.Addr, error)
}

// collectBulbs reads registration responses until the session's read deadline
// expires, deduplicating by MAC so a bulb that responds more than once is
// only listed once.
func collectBulbs(r packetReader) ([]bulb.Bulb, error) {
	bulbs := make([]bulb.Bulb, 0)
	seen := make(map[string]bool)
	buf := make([]byte, udp.MaxSafePayloadSize)
	for {
		_, addr, err := r.Read(buf)
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return bulbs, nil
		} else if err != nil {
			return bulbs, err
		}
		udpAddr, ok := addr.(*net.UDPAddr)
		if !ok {
			continue
		}
		mac := gjson.GetBytes(buf, "result.mac").String()
		if mac == "" || seen[mac] {
			continue
		}
		seen[mac] = true
		ip := udpAddr.IP.String()
		bulbs = append(bulbs, bulb.Bulb{IP: &ip, Mac: &mac})
	}
}

func printJSON(bulbs []bulb.Bulb) {
	b, err := json.Marshal(bulbs)
	cobra.CheckErr(err)
	fmt.Println(string(b))
}

func printPlain(bulbs []bulb.Bulb) {
	fmt.Printf("%-16s%s\n", "IP", "MAC")
	for _, b := range bulbs {
		fmt.Printf("%-16s%s\n", *b.IP, *b.Mac)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&addr, "broadcast-address", "b", "255.255.255.255", "broadcast address")
	listCmd.Flags().BoolVarP(&outputJSON, "json", "j", false, "output list as JSON")
}
