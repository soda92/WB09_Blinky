# STM32WB09 BLE Beacon & P2P Server

This project is a Bluetooth Low Energy (BLE) application for the **STM32WB09** microcontroller (specifically the **NUCLEO-WB09KE** board).

It functions as a **Dual-Mode Device**:
1.  **iBeacon**: Broadcasts presence information.
2.  **P2P Server**: Provides a custom GATT service to control an LED and receive button press notifications.

## Features

- **iBeacon Advertising**: Standard iBeacon format.
- **Custom P2P Service**:
  - **Control LED**: Write to a characteristic to toggle the on-board LED.
  - **Button Notification**: Receive notifications when the on-board button (SW1) is pressed.
- **Tooling**: Custom CLI (`wb09_tool`) for dependency management, building, flashing, and monitoring.
- **Build System**: CMake + Ninja.

## Hardware

- **Board**: NUCLEO-WB09KE (STM32WB09)
- **LED**: Blue LED (connected to `PB12` or `PA7` depending on board revision, mapped via BSP).
- **Button**: SW1 (User Button).

## BLE Specifications

| Entity | UUID | Properties | Description |
| :--- | :--- | :--- | :--- |
| **Service** | `0000FE40-CC7A-482A-984A-7F2ED5B3E58F` | N/A | P2P Server Service |
| **LED Char** | `0000FE41-8E22-4541-9D4C-21EDAE82ED19` | Write w/o Resp | Write 2 Bytes: `[ID] [State]`<br>`ID`: 0x01(Blue), 0x02(Green), 0x03(Red)<br>`State`: 0x01(ON), 0x00(OFF) |
| **Button Char**| `0000FE42-8E22-4541-9D4C-21EDAE82ED19` | Notify | Sends `0x01` on press |

**Note on Advertising**: Due to iBeacon packet size limits, the Device Name ("STM32") is broadcast in the **Scan Response** packet.

## Power Management

The device automatically enters **Low Power Mode (Sleep)** to save energy if:
1.  Advertising times out (default: 60 seconds).
2.  No connection is established.

**To Wake Up**: Press **SW1** (Button 1). This will restart advertising.

## Prerequisites

- **ARM GCC Toolchain** (`arm-none-eabi-gcc`)
- **CMake** & **Ninja**
- **Go** (for the management tool)
- **STM32CubeProgrammer** (CLI version)
- **STM32CubeWB0 SDK** (Installed at `~/STM32Cube/Repository/...` or similar)

## Getting Started

### 1. Build the Management Tool

The project includes a Go-based tool to handle everything. First, build it:

```bash
cd tools
go build -o ../wb09_tool .
cd ..
```

### 2. Initialize & Config (Optional)

If starting fresh or needing to restore libraries:

```bash
./wb09_tool init
./wb09_tool deps
```

You can configure project settings (like disabling Low Power Mode for debugging):

```bash
./wb09_tool config lpm disable
```

### 3. Build & Flash

To build the firmware and flash it to the board:

```bash
./wb09_tool flash
```

This will:
1.  Generate the build system using CMake/Ninja.
2.  Compile the code.
3.  Flash the binary (`WB09_Blinky.bin`) to address `0x10040000`.

### 4. Monitor Output

To view the UART logs (115200 baud):

```bash
./wb09_tool monitor
```

## Development

- **Formatting**: Run `./wb09_tool format` to format C/C++ and Go files.
- **VSCode**: The project is configured for VSCode with `compile_commands.json` support for IntelliSense.

## Project Structure

- `Core/`: Main application logic (`main.c`, interrupts).
- `STM32_BLE/App/`: BLE application logic (`app_ble.c`, `p2p_server_app.c`).
- `STM32_BLE/Target/`: Hardware interface for BLE.
- `Library/`: SDK libraries (copied by `wb09_tool deps`).
- `Drivers/`: HAL and BSP drivers.
- `tools/`: Source code for `wb09_tool`.
- `cmake/`: CMake toolchain files.
