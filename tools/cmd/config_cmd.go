package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"wb09_tool/internal"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage project configuration",
}

var lpmCmd = &cobra.Command{
	Use:   "lpm [enable|disable]",
	Short: "Manage Low Power Mode (disable for debugging)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		switch args[0] {
		case "disable":
			config.DisableLPM = true
			fmt.Println("LPM disabled.")
		case "enable":
			config.DisableLPM = false
			fmt.Println("LPM enabled.")
		default:
			fmt.Println("Invalid argument. Use: enable or disable")
			os.Exit(1)
		}

		if err := internal.SaveConfig(config); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		internal.ApplyConfig(config)
	},
}

var traceCmd = &cobra.Command{
	Use:   "trace [enable|disable]",
	Short: "Manage Debug Application Traces",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		switch args[0] {
		case "enable":
			config.EnableTrace = true
			fmt.Println("Traces enabled.")
		case "disable":
			config.EnableTrace = false
			fmt.Println("Traces disabled.")
		default:
			fmt.Println("Invalid argument. Use: enable or disable")
			os.Exit(1)
		}

		if err := internal.SaveConfig(config); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		internal.ApplyConfig(config)
	},
}

var macCmd = &cobra.Command{
	Use:   "mac [address]",
	Short: "Set Public MAC Address (e.g. 0x112233445566)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		config.MacAddress = args[0]
		fmt.Printf("MAC Address set to %s\n", config.MacAddress)

		if err := internal.SaveConfig(config); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		internal.ApplyConfig(config)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(lpmCmd)
	configCmd.AddCommand(traceCmd)
	configCmd.AddCommand(macCmd)
}
