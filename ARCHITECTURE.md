# Application Architecture & Core Functions

This document outlines the software architecture of the **WB09_Blinky** project, highlighting the core functions responsible for system initialization, BLE management, and application logic.

## 1. System Initialization (`Core`)

The `Core/` directory contains the main entry point and hardware configuration.

### `Core/Src/main.c`
- **`main()`**: The absolute entry point.
  - Initializes the HAL (Hardware Abstraction Layer).
  - Configures the System Clock (`SystemClock_Config`).
  - Initializes peripherals (GPIO, UART, etc.).
  - Calls `MX_APPE_Init` to start the BLE application.
  - Enters the main infinite loop calling `MX_APPE_Process`.

### `Core/Src/app_entry.c`
- **`MX_APPE_Init()`**: Initializes the Application Environment.
  - Sets up the Random Number Generator (`HW_RNG_Init`) and Security (`HW_AES_Init`, `HW_PKA_Init`).
  - Calls `APP_BLE_Init()` to start the BLE stack.
  - Initializes Low Power Manager (`UTIL_LPM_Init`).
  - Initializes Buttons (`Button_Init`) and UART for debug/control (`RxUART_Init`).
- **`MX_APPE_Process()`**: The main background task runner.
  - Calls `UTIL_SEQ_Run()` which executes scheduled tasks (from `stm32_seq.c`).
- **`APPE_ButtonXAction()`**: Weak functions called when buttons are pressed, intended to be overridden or used to trigger sequencer tasks.

## 2. BLE Application Logic (`STM32_BLE`)

The `STM32_BLE/` directory handles the Bluetooth Low Energy stack and custom services.

### `STM32_BLE/App/app_ble.c`
- **`APP_BLE_Init()`**: Initializes the BLE Stack.
  - Calls `BleStack_Init()`.
  - Configures the device address (Public or Static Random).
  - Calls `P2P_SERVER_APP_Init()` to start the custom service.
  - Starts Advertising (`Adv_Request`).
- **`BLEEVT_App_Notification()`**: Global BLE Event Handler.
  - Receives events from the stack (Connection Complete, Disconnection, etc.).
  - Manages Advertising state (restarting on disconnection).

### `STM32_BLE/App/p2p_server_app.c`
This file implements the logic for the **P2P Server** (Peer-to-Peer) custom service.
- **`P2P_SERVER_APP_Init()`**: Registers the P2P service tasks.
  - Registers `P2P_SERVER_Switch_c_SendNotification` task to send button press updates.
- **`P2P_SERVER_Notification()`**: Handles events *specific* to the P2P Service.
  - **`P2P_SERVER_LED_C_WRITE_NO_RESP_EVT`**: Called when a client writes to the LED Characteristic.
    - Turns the Blue LED ON (`0x01`) or OFF (`0x00`).
  - **`P2P_SERVER_SWITCH_C_NOTIFY_ENABLED_EVT`**: Called when a client subscribes to Button notifications.
- **`P2P_SERVER_Switch_c_SendNotification()`**:
  - Updates the Button Characteristic value.
  - Sends a notification to the connected client if subscribed.

### `STM32_BLE/App/p2p_server.c`
This file defines the GATT Service and Characteristics (UUIDs, Properties).
- **`P2P_SERVER_Init()`**:
  - Adds the **P2P Service** (`0000FE40-CC7A-482A-984A-7F2ED5B3E58F`) to the GATT database.
  - Adds **LED Characteristic** (Write w/o Response).
  - Adds **Button Characteristic** (Notify).
- **`P2P_SERVER_NotifyValue()`**: Helper to send notifications to the stack.

## 3. Data Flow

1.  **System Start**: `main()` -> `MX_APPE_Init()` -> `APP_BLE_Init()` -> `P2P_SERVER_Init()`.
2.  **Advertising**: `APP_BLE_Init()` starts advertising. The device becomes visible as "STM32".
3.  **Connection**: When a client connects, `BLEEVT_App_Notification` handles `HCI_LE_CONNECTION_COMPLETE_SUBEVT_CODE`.
4.  **LED Control**:
    - Client writes to LED Characteristic.
    - Stack -> `P2P_SERVER_Notification()` -> `P2P_SERVER_LED_C_WRITE_NO_RESP_EVT` -> `BSP_LED_On/Off`.
5.  **Button Press**:
    - User presses Button 1.
    - Interrupt -> `BSP_PB_Callback` -> `Button_TriggerActions`.
    - `UTIL_SEQ_SetTask(TASK_BUTTON_1)` schedules the task.
    - `APPE_Button1Action` (in `app_entry.c`) calls `P2P_SERVER_Switch_c_SendNotification`.
    - `P2P_SERVER_Switch_c_SendNotification` sends the new value to the client.
