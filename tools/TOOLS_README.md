# WB09 Project Tool

This directory contains `wb09_tool`, a custom Go-based CLI utility designed to manage the lifecycle of the STM32WB09 project. It handles initialization, dependency management, building, flashing, and debugging.

## Installation

To build the tool, run the following command from this directory:

```bash
go build -o ../wb09_tool .
```

This will create the `wb09_tool` binary in the project root.

## Usage

```bash
./wb09_tool [command] [flags]
```

### Commands

#### 1. Initialization
```bash
./wb09_tool init
```
Initializes the project structure from the default template (BLE Beacon).
- **Preserves User Code**: Automatically detects and preserves code between `USER CODE BEGIN` and `USER CODE END` blocks in key files (`main.c`, `app_ble.c`, etc.).
- **Copies SDK**: Copies necessary Drivers and Middleware from the STM32CubeWB0 SDK.
- **Applies Configuration**: Applies settings from `wb09.json` (LPM, Traces, MAC).

#### 2. Dependency Management
```bash
./wb09_tool deps
```
Scans your source code (`Core/Src`, `STM32_BLE/App`, etc.) for `#include` statements and specific symbols.
- **Auto-Discovery**: Finds missing headers and C files in the SDK repository.
- **Auto-Copy**: Copies missing files to `Library/Inc` and `Library/Src` to ensure the project compiles without manually adding every file.
- **Patching**: Applies necessary patches to files like `osal.c` and `cpu_context_switch.s`.

#### 3. Build and Flash
```bash
./wb09_tool flash
```
Builds and flashes the project to the board.
- **Build System**: Uses **CMake** and **Ninja** for fast, incremental builds.
- **Auto-Configure**: Automatically configures the project if `build/build.ninja` is missing.
- **Flashing**: Uses `STM32_Programmer_CLI` to flash the binary (Address: `0x10040000`).

#### 4. Configuration
Manage project settings in `wb09.json`.
- **Low Power Mode**:
  ```bash
  ./wb09_tool config lpm [enable|disable]
  ```
  Disable LPM for easier debugging via ST-Link.
- **Debug Traces**:
  ```bash
  ./wb09_tool config trace [enable|disable]
  ```
  Enable/Disable application debug traces over UART.
- **MAC Address**:
  ```bash
  ./wb09_tool config mac [address]
  ```
  Sets the public MAC address (e.g., `0x112233445566`).

#### 5. Debugging & Monitoring
- **Serial Monitor**:
  ```bash
  ./wb09_tool monitor [port] [baud]
  ```
  Resets the board and captures serial output (default: `/dev/ttyACM0` @ 115200) for 10 seconds. Useful for capturing boot logs.
- **BLE Scan**:
  ```bash
  ./wb09_tool scan
  ```
  Scans for BLE devices for 10 seconds. Highlights the target device if it matches the expected MAC or Name ("STM32").

#### 6. Code Formatting
```bash
./wb09_tool format
```
Formats the codebase to ensure consistent style.
- **C/C++**: Uses `clang-format` (Google style, 2-space indent).
- **Go**: Uses `go fmt` for the tools directory.

#### 7. Cleanup
```bash
./wb09_tool clean
```
Removes all build artifacts (`build/`) and copied libraries (`Library/`, `Drivers/`, etc.).

## Project Structure
- `tools/`: Source code for this tool.
- `tools/cmd/`: Command implementations.
- `tools/internal/`: Shared logic (config, patching, utils).
- `tools/templates/`: Embedded templates (Makefile, etc.).
