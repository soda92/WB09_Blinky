package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"wb09_tool/internal"
)

var lpmCmd = &cobra.Command{
	Use:   "lpm [status|enable|disable]",
	Short: "Manage Low Power Mode (LPM) configuration",
	Long:  `Enable or disable Low Power Mode support to facilitate debugging. Disabling LPM ensures the debugger can always connect.`, 
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config, err := internal.LoadConfig()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		action := args[0]
		switch action {
		case "status":
			if config.DisableLPM {
				fmt.Println("LPM is currently DISABLED (Debugger friendly)")
			} else {
				fmt.Println("LPM is currently ENABLED (Power saving)")
			}
				case "disable":
					config.DisableLPM = true
					if err := internal.SaveConfig(config); err != nil {
						fmt.Printf("Error saving config: %v\n", err)
						return
					}
					internal.ApplyLPMSetting(true)
					fmt.Println("LPM disabled. Debugger should connect reliably.")
				case "enable":
					config.DisableLPM = false
					if err := internal.SaveConfig(config); err != nil {
						fmt.Printf("Error saving config: %v\n", err)
						return
					}
					internal.ApplyLPMSetting(false)
					fmt.Println("LPM enabled.")
				default:
					fmt.Println("Invalid argument. Use: status, enable, or disable")
					os.Exit(1)
				}
			},
		}
		
		func init() {
			rootCmd.AddCommand(lpmCmd)
		}
		
