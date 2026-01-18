cmake_minimum_required(VERSION 3.16)
project(WB09_Blinky C ASM)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Source Files
file(GLOB_RECURSE SOURCES
    "Core/Src/*.c"
    "Library/Src/*.c"
    "STM32_BLE/App/*.c"
    "STM32_BLE/Target/*.c"
    "Drivers/STM32WB0x_HAL_Driver/Src/*.c"
    "Drivers/BSP/STM32WB0x-nucleo/*.c"
    "Middlewares/ST/STM32_BLE/evt_handler/src/*.c"
    "Middlewares/ST/STM32_BLE/stack/config/*.c"
)

# Add explicit startup file and ASM sources from Library
list(APPEND SOURCES 
    "${CMAKE_SOURCE_DIR}/startup_stm32wb09.s"
)

file(GLOB LIB_ASM_SOURCES "Library/Src/*.s")
list(APPEND SOURCES ${LIB_ASM_SOURCES})

# Exclude templates
list(FILTER SOURCES EXCLUDE REGEX ".*_template\.c$")

# Include Directories
set(INCLUDE_DIRS
    Core/Inc
    Library/Inc
    STM32_BLE/App
    STM32_BLE/Target
    Drivers/STM32WB0x_HAL_Driver/Inc
    Drivers/STM32WB0x_HAL_Driver/Inc/Legacy
    Drivers/BSP/STM32WB0x-nucleo
    Drivers/CMSIS/Device/ST/STM32WB0X/Include
    Drivers/CMSIS/Include
    Middlewares/ST/STM32_BLE
    Middlewares/ST/STM32_BLE/stack/include
    Middlewares/ST/STM32_BLE/evt_handler/inc
    Middlewares/ST/STM32_BLE/cryptolib/Inc
)

# Definitions
add_definitions(
    -DUSE_NUCLEO_32
    -DNUCLEO_WB09KE
    -DUSE_HAL_DRIVER
    -DSTM32WB09
)

# Linker Script
set(LINKER_SCRIPT ${CMAKE_SOURCE_DIR}/stm32wb09_flash.ld)

add_executable(${PROJECT_NAME} ${SOURCES})
set_target_properties(${PROJECT_NAME} PROPERTIES SUFFIX ".elf")

target_include_directories(${PROJECT_NAME} PRIVATE ${INCLUDE_DIRS})

target_link_options(${PROJECT_NAME} PRIVATE
    -T${LINKER_SCRIPT}
    -Wl,-Map=${CMAKE_BINARY_DIR}/${PROJECT_NAME}.map
    -Wl,--cref
)

# Precompiled Libraries
target_link_libraries(${PROJECT_NAME} PRIVATE
    ${CMAKE_SOURCE_DIR}/Middlewares/ST/STM32_BLE/stack/lib/stm32wb0x_ble_stack.a
    ${CMAKE_SOURCE_DIR}/Middlewares/ST/STM32_BLE/cryptolib/libcrypto.a
    c
    m
    nosys
)

# Generate Hex and Bin
add_custom_command(TARGET ${PROJECT_NAME} POST_BUILD
    COMMAND ${CMAKE_OBJCOPY} -O ihex $<TARGET_FILE:${PROJECT_NAME}> ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.hex
    COMMAND ${CMAKE_OBJCOPY} -O binary -S $<TARGET_FILE:${PROJECT_NAME}> ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.bin
    COMMAND ${CMAKE_SIZE} $<TARGET_FILE:${PROJECT_NAME}>
    COMMAND ${CMAKE_COMMAND} -E copy_if_different ${CMAKE_BINARY_DIR}/compile_commands.json ${CMAKE_SOURCE_DIR}/compile_commands.json
    COMMENT "Building hex, bin files and copying compile_commands.json"
)

# Flash Target
set(STM32_PROG "/home/soda/STMicroelectronics/STM32Cube/STM32CubeProgrammer/bin/STM32_Programmer_CLI")
add_custom_target(flash
    COMMAND ${STM32_PROG} -c port=SWD -w ${CMAKE_BINARY_DIR}/${PROJECT_NAME}.bin 0x10040000 -v -rst
    DEPENDS ${PROJECT_NAME}
    COMMENT "Flashing via STM32CubeProgrammer..."
    USES_TERMINAL
)
