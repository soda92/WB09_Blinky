#include "app_ble.h"
#include "main.h"
#include "app_conf.h"
#include "ble.h"
#include "ble_stack_user_cfg.h"
#include <stdio.h>

#define BLE_DYN_ALLOC_SIZE \
    BLE_STACK_TOTAL_BUFFER_SIZE( \
        CFG_NUM_RADIO_TASKS, \
        0, \
        CFG_BLE_NUM_GATT_ATTRIBUTES, \
        1, \
        CFG_BLE_MBLOCK_COUNT, \
        1, \
        0, \
        0, \
        0, \
        3, \
        0, \
        0, \
        0, \
        0, \
        0, \
        0, \
        0, \
        0, \
        0, \
        0, \
        0, \
        CFG_BLE_ISR0_FIFO_SIZE, \
        CFG_BLE_ISR1_FIFO_SIZE, \
        CFG_BLE_USER_FIFO_SIZE \
    )

static uint32_t a_BLE_Stack_Buffer[DIV_CEIL(BLE_DYN_ALLOC_SIZE, 4)];

void APP_BLE_Init(void)
{
    BLE_STACK_InitTypeDef BleStack_Init_Param = {
        .BLEStartRamAddress = (uint8_t*)a_BLE_Stack_Buffer,
        .TotalBufferSize = BLE_DYN_ALLOC_SIZE,
        .NumAttrRecords = CFG_BLE_NUM_GATT_ATTRIBUTES,
        .MaxNumOfClientProcs = 1,
        .NumOfRadioTasks = CFG_NUM_RADIO_TASKS,
        .NumOfEATTChannels = 0,
        .NumBlockCount = CFG_BLE_MBLOCK_COUNT,
        .ATT_MTU = 23,
        .MaxConnEventLength = 0xFFFFFFFF,
        .SleepClockAccuracy = 500,
        .NumOfAdvDataSet = 1,
        .NumOfSubeventsPAwR = 0,
        .MaxPAwRSubeventDataCount = 0,
        .NumOfAuxScanSlots = 0,
        .NumOfSyncSlots = 0,
        .FilterAcceptListSizeLog2 = 3,
        .L2CAP_MPS = 23,
        .L2CAP_NumChannels = 0,
        .CTE_MaxNumAntennaIDs = 0,
        .CTE_MaxNumIQSamples = 0,
        .NumOfSyncBIG = 0,
        .NumOfBrcBIG = 0,
        .NumOfSyncBIS = 0,
        .NumOfBrcBIS = 0,
        .NumOfCIG = 0,
        .NumOfCIS = 0,
        .ExtraLLProcedureContexts = 0,
        .isr0_fifo_size = CFG_BLE_ISR0_FIFO_SIZE,
        .isr1_fifo_size = CFG_BLE_ISR1_FIFO_SIZE,
        .user_fifo_size = CFG_BLE_USER_FIFO_SIZE
    };

    if (BLE_STACK_Init(&BleStack_Init_Param) != BLE_STATUS_SUCCESS) {
        printf("BLE Init Error\r\n");
        return;
    }

    // Configure Beacon Advertising
    uint8_t bd_addr[6] = {0x66, 0x55, 0x44, 0x33, 0x22, 0x11}; // Random address
    aci_hal_write_config_data(CONFIG_DATA_PUBADDR_OFFSET, CONFIG_DATA_PUBADDR_LEN, bd_addr);
    
    // Set Tx Power
    aci_hal_set_tx_power_level(0, 24); // 0 dBm

    // Set Advertising Data (Flags + Name "WB09")
    uint8_t adv_data[32];
    uint8_t index = 0;
    
    // Flags
    adv_data[index++] = 2;
    adv_data[index++] = AD_TYPE_FLAGS;
    adv_data[index++] = FLAG_BIT_LE_GENERAL_DISCOVERABLE_MODE | FLAG_BIT_BR_EDR_NOT_SUPPORTED;
    
    // Name
    adv_data[index++] = 1 + 4; // Length
    adv_data[index++] = AD_TYPE_COMPLETE_LOCAL_NAME;
    adv_data[index++] = 'W';
    adv_data[index++] = 'B';
    adv_data[index++] = '0';
    adv_data[index++] = '9';

    aci_gap_set_advertising_configuration(0, GAP_MODE_GENERAL_DISCOVERABLE,
                                          HCI_ADV_EVENT_PROP_CONNECTABLE | HCI_ADV_EVENT_PROP_SCANNABLE | HCI_ADV_EVENT_PROP_LEGACY,
                                          160, 160, // 100 ms
                                          HCI_ADV_CH_ALL,
                                          0, NULL, // Peer Addr
                                          HCI_ADV_FILTER_NONE,
                                          0, // Tx Power
                                          HCI_PHY_LE_1M, 0, HCI_PHY_LE_1M,
                                          0, 0);

    aci_gap_set_advertising_data(0, ADV_COMPLETE_DATA, index, adv_data);

    aci_gap_set_advertising_enable(ENABLE, 1, NULL);
    
    printf("BLE Beacon Started!\r\n");
}

void APP_BLE_Tick(void)
{
    BLE_STACK_Tick();
}

void BLEEVT_App_Notification(const hci_pckt *pckt)
{
  /* User can parse the packet here */
}
