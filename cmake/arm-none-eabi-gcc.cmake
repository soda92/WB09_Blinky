set(CMAKE_SYSTEM_NAME Generic)
set(CMAKE_SYSTEM_PROCESSOR cortex-m0plus)

set(CMAKE_C_COMPILER arm-none-eabi-gcc)
set(CMAKE_CXX_COMPILER arm-none-eabi-g++)
set(CMAKE_ASM_COMPILER arm-none-eabi-gcc)

set(CMAKE_OBJCOPY arm-none-eabi-objcopy)
set(CMAKE_SIZE arm-none-eabi-size)

set(CMAKE_TRY_COMPILE_TARGET_TYPE STATIC_LIBRARY)

set(CMAKE_C_FLAGS "-mcpu=cortex-m0plus -mthumb -Wall -w -fdata-sections -ffunction-sections" CACHE INTERNAL "C Compiler options")
set(CMAKE_CXX_FLAGS "-mcpu=cortex-m0plus -mthumb -Wall -fdata-sections -ffunction-sections" CACHE INTERNAL "C++ Compiler options")
set(CMAKE_ASM_FLAGS "-mcpu=cortex-m0plus -mthumb -x assembler-with-cpp" CACHE INTERNAL "ASM Compiler options")

set(CMAKE_EXE_LINKER_FLAGS "-mcpu=cortex-m0plus -mthumb -specs=nano.specs -Wl,--gc-sections" CACHE INTERNAL "Linker options")
