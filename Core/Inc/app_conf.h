#ifndef APP_CONF_H
#define APP_CONF_H

#include <stdint.h>

#define CFG_BLE_NUM_RADIO_TASKS                         2

/* BLE Stack Configuration */
#define CFG_BLE_CONNECTION_ENABLED                      1
#define CFG_BLE_CONTROLLER_SCAN_ENABLED                 0
#define CFG_BLE_CONTROLLER_PRIVACY_ENABLED              0
#define CFG_BLE_SECURE_CONNECTIONS_ENABLED              0
#define CFG_BLE_CONTROLLER_DATA_LENGTH_EXTENSION_ENABLED 0
#define CFG_BLE_CONTROLLER_2M_CODED_PHY_ENABLED         0
#define CFG_BLE_CONTROLLER_EXT_ADV_SCAN_ENABLED         0
#define CFG_BLE_L2CAP_COS_ENABLED                       0
#define CFG_BLE_CONTROLLER_PERIODIC_ADV_ENABLED         0
#define CFG_BLE_CONTROLLER_PERIODIC_ADV_WR_ENABLED      0
#define CFG_BLE_CONTROLLER_CTE_ENABLED                  0
#define CFG_BLE_CONTROLLER_POWER_CONTROL_ENABLED        0
#define CFG_BLE_CONTROLLER_CHAN_CLASS_ENABLED           0
#define CFG_BLE_CONTROLLER_BIS_ENABLED                  0
#define CFG_BLE_CONNECTION_SUBRATING_ENABLED            0
#define CFG_BLE_CONTROLLER_CIS_ENABLED                  0

/* Memory Buffer Sizes (Minimal for Beacon) */
#define CFG_BLE_NUM_GATT_ATTRIBUTES                     10
#define CFG_BLE_NUM_GATT_SERVICES                       2
#define CFG_BLE_ATT_VALUE_ARRAY_SIZE                    64
#define CFG_BLE_NUM_LINK                                1
#define CFG_BLE_DATA_LENGTH_EXTENSION                   0
#define CFG_BLE_PREPARE_WRITE_LIST_SIZE                 0
#define CFG_BLE_MBLOCK_COUNT                            16

#define CFG_BLE_ISR0_FIFO_SIZE                          256
#define CFG_BLE_ISR1_FIFO_SIZE                          256
#define CFG_BLE_USER_FIFO_SIZE                          256

/* Define required by ble_stack_user_cfg.c */
#define BLESTACK_CONTROLLER_ONLY 0

#endif /* APP_CONF_H */
