# codesys-sensor-client

This example just provides some functions to convert CAN data back into nummeric values. Your CAN library probably provides similar functionality (byte to float conversions IEEE 754 style).

- `ST_BYTE_TO_INT`: Convert 2 bytes into an int16 value.
- `ST_BYTE_TO_UINT`: Convert 2 bytes into an uint16 value.
- `ST_BYTE_TO_LREAL`: Convert 8 bytes into an float64 value.