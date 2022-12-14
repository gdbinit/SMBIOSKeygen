## SMBIOSKeygen

This is a port to Go of [macserial](https://github.com/acidanthera/OpenCorePkg/tree/master/Utilities/macserial) from [OpenCore](https://github.com/acidanthera/OpenCorePkg) project and [GenSMBIOS](https://github.com/corpnewt/GenSMBIOS) to generate valid SMBIOS serial numbers for OpenCore and other bootloaders.

It can retrieve the system information via IOKit using CGO but only for the machine where it's being built (due to CGO cross-compilation issues).

The pure Go code compiles (and tested) for macOS x64 and ARM64, Linux, and Windows (use the `windows` Makefile target to build it). The beauty of Go cross-compiling!

Use the `-k` command to generate all the needed information for OpenCore. The default model is `iMacPro1,1` but you can modify via options (`-m` in this case). All the available models can be listed with the `-l` command.

Motivation is my personal dislike of GenSMBIOS (and other scripts) downloading (unverified) software from the internet. Truth be told, it does its job and it's used by a lot of people so don't interpret this as a critic.

I just started liking Go a lot and lately have a lot of free time so it was maybe time to give back something to this great community (or just add another project to Github).

Enjoy,

fG!

## Other notes

The `scripts` folder contains an updated `update_generated.py` script to generate the Go version of `modelinfo_autogen` in case there are updates upstream and you want to merge them. There should be no updates since Apple Silicon models are not included here.

## References

https://www.insanelymac.com/forum/topic/303073-pattern-of-mlb-main-logic-board/
https://github.com/sickcodes/osx-serial-generator
https://dortania.github.io/OpenCore-Post-Install/universal/iservices.html

## macserial - original README

macserial is a tool that obtains and decodes Mac serial number and board identifier to provide more information about the production of your hardware. Works as a decent companion to [Apple Check Coverage](https://checkcoverage.apple.com) and [Apple Specs](http://support-sp.apple.com/sp/index?page=cpuspec&cc=HTD5) portal. Check the [format description](https://github.com/acidanthera/OpenCorePkg/blob/master/Utilities/macserial/FORMAT.md) for more details.

Should be built with a compiler supporting C99. Prebuilt binaries are available for macOS 10.4 and higher.

Run with `-h` argument to see all available arguments.
