package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for Bluetooth devices",
	Run: func(cmd *cobra.Command, args []string) {
		adapter := bluetooth.DefaultAdapter
		err := adapter.Enable()
		if err != nil {
			fmt.Printf("Error enabling bluetooth adapter: %v\n", err)
			return
		}

		fmt.Println("Scanning for 10 seconds...")
		go func() {
			time.Sleep(10 * time.Second)
			adapter.StopScan()
			fmt.Println("Scan finished.")
		}()

		err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
			addr := device.Address.String()
			name := device.LocalName()

			// Normalize check
			if strings.Contains(strings.ToUpper(addr), "11:22:33:44:55:66") ||
				strings.Contains(strings.ToUpper(addr), "66:55:44:33:22:11") ||
				name == "STM32" {
				fmt.Printf(">>> FOUND TARGET: %s [%s] RSSI: %d\n", addr, name, device.RSSI)
			} else {
				fmt.Printf("Found: %s [%s] RSSI: %d\n", addr, name, device.RSSI)
			}
		})

		if err != nil {
			fmt.Printf("Scan error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
