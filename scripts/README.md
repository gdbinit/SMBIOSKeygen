The `update_generated.py` is an updated version of `https://github.com/acidanthera/OpenCorePkg/blob/master/AppleModels/update_generated.py` that generates the Go version of `modelinfo_autogen.h` called `modelinfo_autogen.go`.

This file contains hardware information that is necessary to generate the serial numbers.

The new function is called `export_db_macserial_go`.

The script depends on the data available in that OpenCorePkg folder so it should be copied there and then copy the generated Go file to SMBIOSKeygen folder.
