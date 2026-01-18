/* USER CODE BEGIN Header */
/**
  ******************************************************************************
  * @file           : main.c
  * @brief          : Main program body
  ******************************************************************************
  * @attention
  *
  * Copyright (c) 2025 STMicroelectronics.
  * All rights reserved.
  *
  * This software is licensed under terms that can be found in the LICENSE file
  * in the root directory of this software component.
  * If no LICENSE file comes with this software, it is provided AS-IS.
  *
  ******************************************************************************
  */
/* USER CODE END Header */
/* Includes ------------------------------------------------------------------*/
#include "main.h"

/* Private includes ----------------------------------------------------------*/
/* USER CODE BEGIN Includes */
#include <stdio.h>
#include "stm32wb0x_ll_utils.h"
#include "stm32wb0x_ll_adc.h"
/* USER CODE END Includes */

/* Private typedef -----------------------------------------------------------*/
/* USER CODE BEGIN PTD */

/* USER CODE END PTD */

/* Private define ------------------------------------------------------------*/
/* USER CODE BEGIN PD */

/* USER CODE END PD */

/* Private macro -------------------------------------------------------------*/
/* USER CODE BEGIN PM */

/* USER CODE END PM */

/* Private variables ---------------------------------------------------------*/

COM_InitTypeDef BspCOMInit;
ADC_HandleTypeDef hadc1;

/* USER CODE BEGIN PV */

/* USER CODE END PV */

/* Private function prototypes -----------------------------------------------*/
void SystemClock_Config(void);
void PeriphCommonClock_Config(void);
static void MX_GPIO_Init(void);
static void MX_ADC1_Init(void);
/* USER CODE BEGIN PFP */

/* USER CODE END PFP */

/* Private user code ---------------------------------------------------------*/
/* USER CODE BEGIN 0 */

/* USER CODE END 0 */

/**
  * @brief  The application entry point.
  * @retval int
  */
int main(void)
{

  /* USER CODE BEGIN 1 */

  /* USER CODE END 1 */

  /* MCU Configuration--------------------------------------------------------*/

  /* Reset of all peripherals, Initializes the Flash interface and the Systick. */
  HAL_Init();

  /* USER CODE BEGIN Init */

  /* USER CODE END Init */

  /* Configure the system clock */
  SystemClock_Config();

  /* Configure the peripherals common clocks */
  PeriphCommonClock_Config();

  /* USER CODE BEGIN SysInit */

  /* USER CODE END SysInit */

  /* Initialize all configured peripherals */
  MX_GPIO_Init();
  MX_ADC1_Init();
  /* USER CODE BEGIN 2 */

  /* USER CODE END 2 */

  /* Initialize leds */
  BSP_LED_Init(LED_BLUE);
  BSP_LED_Init(LED_GREEN);
  BSP_LED_Init(LED_RED);

  /* Initialize USER push-button, will be used to trigger an interrupt each time it's pressed.*/
  BSP_PB_Init(B1, BUTTON_MODE_EXTI);
  BSP_PB_Init(B2, BUTTON_MODE_EXTI);
  BSP_PB_Init(B3, BUTTON_MODE_EXTI);

  /* Initialize COM1 port (115200, 8 bits (7-bit data + 1 stop bit), no parity */
  BspCOMInit.BaudRate   = 115200;
  BspCOMInit.WordLength = COM_WORDLENGTH_8B;
  BspCOMInit.StopBits   = COM_STOPBITS_1;
  BspCOMInit.Parity     = COM_PARITY_NONE;
  BspCOMInit.HwFlowCtl  = COM_HWCONTROL_NONE;
  if (BSP_COM_Init(COM1, &BspCOMInit) != BSP_ERROR_NONE)
  {
    Error_Handler();
  }

  /* Infinite loop */
  /* USER CODE BEGIN WHILE */
  while (1)
  {
    /* Blue LED ON */
    BSP_LED_On(LED_BLUE);
    BSP_LED_Off(LED_GREEN);
    BSP_LED_Off(LED_RED);
    HAL_Delay(500);

    /* Green LED ON */
    BSP_LED_Off(LED_BLUE);
    BSP_LED_On(LED_GREEN);
    BSP_LED_Off(LED_RED);
    HAL_Delay(500);

    /* Red LED ON */
    BSP_LED_Off(LED_BLUE);
    BSP_LED_Off(LED_GREEN);
    BSP_LED_On(LED_RED);
    HAL_Delay(500);

    /* USER CODE END WHILE */

    /* USER CODE BEGIN 3 */
  }
  /* USER CODE END 3 */
}

/**
  * @brief System Clock Configuration
  * @retval None
  */
void SystemClock_Config(void)
{
  RCC_ClkInitTypeDef RCC_ClkInitStruct = {0};

  /** Configure the SYSCLKSource and SYSCLKDivider
  */
  RCC_ClkInitStruct.SYSCLKSource = RCC_SYSCLKSOURCE_HSI;
  RCC_ClkInitStruct.SYSCLKDivider = RCC_RC64MPLL_DIV1;

  if (HAL_RCC_ClockConfig(&RCC_ClkInitStruct, FLASH_WAIT_STATES_1) != HAL_OK)
  {
    Error_Handler();
  }
}

/**
  * @brief Peripherals Common Clock Configuration
  * @retval None
  */
void PeriphCommonClock_Config(void)
{
  RCC_PeriphCLKInitTypeDef PeriphClkInitStruct = {0};

  /** Initializes the peripherals clock
  */
  PeriphClkInitStruct.PeriphClockSelection = RCC_PERIPHCLK_SMPS;
  PeriphClkInitStruct.SmpsDivSelection = RCC_SMPSCLK_DIV4;

  if (HAL_RCCEx_PeriphCLKConfig(&PeriphClkInitStruct) != HAL_OK)
  {
    Error_Handler();
  }
}

/**
  * @brief ADC1 Initialization Function
  * @param None
  * @retval None
  */
static void MX_ADC1_Init(void)
{

  /* USER CODE BEGIN ADC1_Init 0 */

  /* USER CODE END ADC1_Init 0 */

  ADC_ChannelConfTypeDef ConfigChannel = {0};

  /* USER CODE BEGIN ADC1_Init 1 */

  /* USER CODE END ADC1_Init 1 */

  /** Common config
  */
  hadc1.Instance = ADC1;
  hadc1.Init.ConversionType = ADC_CONVERSION_WITH_DS;
  hadc1.Init.SequenceLength = 1;
  hadc1.Init.SamplingMode = ADC_SAMPLING_AT_START;
  hadc1.Init.SampleRate = ADC_SAMPLE_RATE_16;
  hadc1.Init.InvertOutputMode = ADC_DATA_INVERT_NONE;
  hadc1.Init.Overrun = ADC_NEW_DATA_IS_LOST;
  hadc1.Init.ContinuousConvMode = DISABLE;
  hadc1.Init.DownSamplerConfig.DataWidth = ADC_DS_DATA_WIDTH_12_BIT;
  hadc1.Init.DownSamplerConfig.DataRatio = ADC_DS_RATIO_1;
  if (HAL_ADC_Init(&hadc1) != HAL_OK)
  {
    Error_Handler();
  }

  /** Configure Regular Channel
  */
  ConfigChannel.Channel = ADC_CHANNEL_TEMPSENSOR;
  ConfigChannel.Rank = ADC_RANK_1;
  ConfigChannel.VoltRange = ADC_VIN_RANGE_1V2;
  ConfigChannel.CalibrationPoint.Number = ADC_CALIB_NONE;
  ConfigChannel.CalibrationPoint.Gain = 0x00;
  ConfigChannel.CalibrationPoint.Offset = 0x00;
  if (HAL_ADC_ConfigChannel(&hadc1, &ConfigChannel) != HAL_OK)
  {
    Error_Handler();
  }
  /* USER CODE BEGIN ADC1_Init 2 */

  /* USER CODE END ADC1_Init 2 */

}

/**
  * @brief GPIO Initialization Function
  * @param None
  * @retval None
  */
static void MX_GPIO_Init(void)
{
  GPIO_InitTypeDef GPIO_InitStruct = {0};
  /* USER CODE BEGIN MX_GPIO_Init_1 */

  /* USER CODE END MX_GPIO_Init_1 */

  /* GPIO Ports Clock Enable */
  __HAL_RCC_GPIOB_CLK_ENABLE();
  __HAL_RCC_GPIOA_CLK_ENABLE();

  /*Configure GPIO pin : SPI3_SCK_Pin */
  GPIO_InitStruct.Pin = SPI3_SCK_Pin;
  GPIO_InitStruct.Mode = GPIO_MODE_AF_PP;
  GPIO_InitStruct.Pull = GPIO_NOPULL;
  GPIO_InitStruct.Speed = GPIO_SPEED_FREQ_VERY_HIGH;
  GPIO_InitStruct.Alternate = GPIO_AF4_SPI3;
  HAL_GPIO_Init(SPI3_SCK_GPIO_Port, &GPIO_InitStruct);

  /*Configure GPIO pins : SPI3_MISO_Pin SPI3_NSS_Pin SPI3_MOSI_Pin */
  GPIO_InitStruct.Pin = SPI3_MISO_Pin|SPI3_NSS_Pin|SPI3_MOSI_Pin;
  GPIO_InitStruct.Mode = GPIO_MODE_AF_PP;
  GPIO_InitStruct.Pull = GPIO_NOPULL;
  GPIO_InitStruct.Speed = GPIO_SPEED_FREQ_VERY_HIGH;
  GPIO_InitStruct.Alternate = GPIO_AF3_SPI3;
  HAL_GPIO_Init(GPIOA, &GPIO_InitStruct);

  /*Configure GPIO pins : I2C1_SDA_Pin I2C1_SCL_Pin */
  GPIO_InitStruct.Pin = I2C1_SDA_Pin|I2C1_SCL_Pin;
  GPIO_InitStruct.Mode = GPIO_MODE_AF_OD;
  GPIO_InitStruct.Pull = GPIO_NOPULL;
  GPIO_InitStruct.Speed = GPIO_SPEED_FREQ_VERY_HIGH;
  GPIO_InitStruct.Alternate = GPIO_AF0_I2C1;
  HAL_GPIO_Init(GPIOB, &GPIO_InitStruct);

  /**/
  HAL_PWREx_DisableGPIOPullUp(PWR_GPIO_B, PWR_GPIO_BIT_3|PWR_GPIO_BIT_7|PWR_GPIO_BIT_6);

  /**/
  HAL_PWREx_DisableGPIOPullUp(PWR_GPIO_A, PWR_GPIO_BIT_8|PWR_GPIO_BIT_9|PWR_GPIO_BIT_11);

  /**/
  HAL_PWREx_DisableGPIOPullDown(PWR_GPIO_B, PWR_GPIO_BIT_3|PWR_GPIO_BIT_7|PWR_GPIO_BIT_6);

  /**/
  HAL_PWREx_DisableGPIOPullDown(PWR_GPIO_A, PWR_GPIO_BIT_8|PWR_GPIO_BIT_9|PWR_GPIO_BIT_11);

  /* USER CODE BEGIN MX_GPIO_Init_2 */

  /* USER CODE END MX_GPIO_Init_2 */
}

/* USER CODE BEGIN 4 */
void HAL_GPIO_EXTI_Callback(GPIO_TypeDef* GPIOx, uint16_t GPIO_Pin)
{
  if (GPIO_Pin == B1_PIN && GPIOx == B1_GPIO_PORT)
  {
    BSP_PB_Callback(B1);
  }
  else if (GPIO_Pin == B2_PIN && GPIOx == B2_GPIO_PORT)
  {
    BSP_PB_Callback(B2);
  }
  else if (GPIO_Pin == B3_PIN && GPIOx == B3_GPIO_PORT)
  {
    BSP_PB_Callback(B3);
  }
}

void BSP_PB_Callback(Button_TypeDef Button)
{
  switch (Button)
  {
    case B1:
      printf("Button B1 pressed!\r\n");
      
      ADC_ChannelConfTypeDef sConfig = {0};
      sConfig.Rank = ADC_RANK_1;
      sConfig.VoltRange = ADC_VIN_RANGE_1V2;
      sConfig.CalibrationPoint.Number = ADC_CALIB_NONE;
      sConfig.CalibrationPoint.Gain = 0;
      sConfig.CalibrationPoint.Offset = 0;

      // Measure Temperature
      sConfig.Channel = ADC_CHANNEL_TEMPSENSOR;
      if (HAL_ADC_ConfigChannel(&hadc1, &sConfig) == HAL_OK)
      {
        HAL_ADC_Start(&hadc1);
        if (HAL_ADC_PollForConversion(&hadc1, 100) == HAL_OK)
        {
          uint32_t rawTemp = HAL_ADC_GetValue(&hadc1);
          int32_t temperature = __LL_ADC_CALC_TEMPERATURE(rawTemp, LL_ADC_DS_DATA_WIDTH_12_BIT);
          printf("Temperature: %ld C\r\n", temperature);
        }
      }

      // Measure Vbat
      sConfig.Channel = ADC_CHANNEL_VBAT;
      if (HAL_ADC_ConfigChannel(&hadc1, &sConfig) == HAL_OK)
      {
        HAL_ADC_Start(&hadc1);
        if (HAL_ADC_PollForConversion(&hadc1, 100) == HAL_OK)
        {
          uint32_t rawVbat = HAL_ADC_GetValue(&hadc1);
          // Calc voltage based on 1.2V ref, multiply by 3 for divider
          uint32_t volt = __LL_ADC_CALC_DATA_TO_VOLTAGE(LL_ADC_VIN_RANGE_1V2, rawVbat, LL_ADC_DS_DATA_WIDTH_12_BIT);
          printf("Vbat: %ld mV\r\n", volt * 3);
        }
      }
      break;
    case B2:
      printf("Button B2 pressed!\r\n");
      break;
    case B3:
      printf("Button B3 pressed!\r\n");
      printf("UID: 0x%08lX%08lX\r\n", LL_GetUID_Word1(), LL_GetUID_Word0());
      printf("Flash Size Reg: 0x%lX\r\n", LL_GetFlashSize());
      break;
    default:
      break;
  }
}
/* USER CODE END 4 */

/**
  * @brief  This function is executed in case of error occurrence.
  * @retval None
  */
void Error_Handler(void)
{
  /* USER CODE BEGIN Error_Handler_Debug */
  /* User can add his own implementation to report the HAL error return state */
  __disable_irq();
  while (1)
  {
  }
  /* USER CODE END Error_Handler_Debug */
}
#ifdef USE_FULL_ASSERT
/**
  * @brief  Reports the name of the source file and the source line number
  *         where the assert_param error has occurred.
  * @param  file: pointer to the source file name
  * @param  line: assert_param error line source number
  * @retval None
  */
void assert_failed(uint8_t *file, uint32_t line)
{
  /* USER CODE BEGIN 6 */
  /* User can add his own implementation to report the file name and line number,
     ex: printf("Wrong parameters value: file %s on line %d\r\n", file, line) */
  /* USER CODE END 6 */
}
#endif /* USE_FULL_ASSERT */
