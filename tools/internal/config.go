package internal

const (
	// Default SDK Path (can be overridden via flags if needed, but hardcoded for now)
	SDKPath = "/home/soda/STM32Cube/Repository/STM32Cube_FW_WB0_V1.4.0"

	// Project Paths
	CoreDir        = "Core"
	LibraryDir     = "Library"
	DriversDir     = "Drivers"
	MiddlewaresDir = "Middlewares"
	STM32BLEDir    = "STM32_BLE"
	BuildDir       = "build"

	// Destination sub-directories for Library
	LibSrcDir = "Library/Src"
	LibIncDir = "Library/Inc"
)

// Default Template Path (BLE Beacon)
const DefaultTemplatePath = "/home/soda/STM32Cube/Repository/STM32Cube_FW_WB0_V1.4.0/Projects/NUCLEO-WB09KE/Applications/BLE/BLE_Beacon"
