//go:build iokit
// +build iokit

// use a iokit tag because otherwise we have cross compilation issues with the CGO
// so we just build this for whatever native platform it's being built on to avoid CGO cross compilation issues
package main

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOKitLib.h>

#define ARRAY_SIZE(arr) (sizeof(arr) / sizeof((arr)[0]))
#define SZUUID 16
#define PRIUUID "%02X%02X%02X%02X-%02X%02X-%02X%02X-%02X%02X-%02X%02X%02X%02X%02X%02X"
#define CASTUUID(uuid) (uuid)[0], (uuid)[1], (uuid)[2], (uuid)[3], (uuid)[4], (uuid)[5], (uuid)[6], \
  (uuid)[7], (uuid)[8], (uuid)[9], (uuid)[10], (uuid)[11], (uuid)[12], (uuid)[13], (uuid)[14], (uuid)[15]

static CFTypeRef get_ioreg_entry(const char *path, CFStringRef name, CFTypeID type) {
  CFTypeRef value = NULL;
  io_registry_entry_t entry = IORegistryEntryFromPath(kIOMasterPortDefault, path);
  if (entry) {
    value = IORegistryEntryCreateCFProperty(entry, name, kCFAllocatorDefault, 0);
    if (value) {
      if (CFGetTypeID(value) != type) {
        CFRelease(value);
        value = NULL;
        printf("%s in %s has wrong type!\n", CFStringGetCStringPtr(name, kCFStringEncodingMacRoman), path);
      }
    } else {
      printf("Failed to find to %s in %s!\n", CFStringGetCStringPtr(name, kCFStringEncodingMacRoman), path);
    }
    IOObjectRelease(entry);
  } else {
    printf("Failed to connect to %s!\n", path);
  }
  return value;
}


void get_system_info() {
  CFDataRef model    = get_ioreg_entry("IODeviceTree:/", CFSTR("model"), CFDataGetTypeID());
  CFDataRef board    = get_ioreg_entry("IODeviceTree:/", CFSTR("board-id"), CFDataGetTypeID());
  CFDataRef efiver   = get_ioreg_entry("IODeviceTree:/rom", CFSTR("version"), CFDataGetTypeID());
  CFStringRef serial = get_ioreg_entry("IODeviceTree:/", CFSTR("IOPlatformSerialNumber"), CFStringGetTypeID());
  CFStringRef hwuuid = get_ioreg_entry("IODeviceTree:/", CFSTR("IOPlatformUUID"), CFStringGetTypeID());
  CFDataRef smuuid   = get_ioreg_entry("IODeviceTree:/efi/platform", CFSTR("system-id"), CFDataGetTypeID());
  CFDataRef rom      = get_ioreg_entry("IODeviceTree:/options", CFSTR("4D1EDE05-38C7-4A6A-9CC6-4BCCA8B38C14:ROM"), CFDataGetTypeID());
  CFDataRef mlb      = get_ioreg_entry("IODeviceTree:/options", CFSTR("4D1EDE05-38C7-4A6A-9CC6-4BCCA8B38C14:MLB"), CFDataGetTypeID());

  CFDataRef   pwr[5] = {0};
  CFStringRef pwrname[5] = {
    CFSTR("Gq3489ugfi"),
    CFSTR("Fyp98tpgj"),
    CFSTR("kbjfrfpoJU"),
    CFSTR("oycqAZloTNDm"),
    CFSTR("abKPld1EcMni"),
  };

  for (size_t i = 0; i < ARRAY_SIZE(pwr); i++)
    pwr[i] = get_ioreg_entry("IOPower:/", pwrname[i], CFDataGetTypeID());

  if (model) {
    printf("%14s: %.*s\n", "Model", (int)CFDataGetLength(model), CFDataGetBytePtr(model));
    CFRelease(model);
  }

  if (board) {
    printf("%14s: %.*s\n", "Board ID", (int)CFDataGetLength(board), CFDataGetBytePtr(board));
    CFRelease(board);
  }

  if (efiver) {
    printf("%14s: %.*s\n", "FW Version", (int)CFDataGetLength(efiver), CFDataGetBytePtr(efiver));
    CFRelease(efiver);
  }

  if (hwuuid) {
    printf("%14s: %s\n", "Hardware UUID", CFStringGetCStringPtr(hwuuid, kCFStringEncodingMacRoman));
    CFRelease(hwuuid);
  }

  puts("");

  if (serial) {
    const char *cstr = CFStringGetCStringPtr(serial, kCFStringEncodingMacRoman);
    printf("%14s: %s\n", "Serial Number", cstr);
    CFRelease(serial);
    puts("");
  }

  if (smuuid) {
    if (CFDataGetLength(smuuid) == SZUUID) {
      const uint8_t *p = CFDataGetBytePtr(smuuid);
      printf("%14s: " PRIUUID "\n", "System ID", CASTUUID(p));
    }
    CFRelease(smuuid);
  }

  if (rom) {
    if (CFDataGetLength(rom) == 6) {
      const uint8_t *p = CFDataGetBytePtr(rom);
      printf("%14s: %02X%02X%02X%02X%02X%02X\n", "ROM", p[0], p[1], p[2], p[3], p[4], p[5]);
    }
    CFRelease(rom);
  }

  if (mlb) {
    printf("%14s: %.*s\n", "MLB", (int)CFDataGetLength(mlb), CFDataGetBytePtr(mlb));
    CFRelease(mlb);
  }

  puts("");

  for (size_t i = 0; i < ARRAY_SIZE(pwr); i++) {
    if (pwr[i]) {
      printf("%14s: ", CFStringGetCStringPtr(pwrname[i], kCFStringEncodingMacRoman));
      const uint8_t *p = CFDataGetBytePtr(pwr[i]);
      CFIndex sz = CFDataGetLength(pwr[i]);
      for (CFIndex j = 0; j < sz; j++)
        printf("%02X", p[j]);
      puts("");
      CFRelease(pwr[i]);
    }
  }

  puts("");
}
*/
import "C"

func GetSystemInfo() {
	C.get_system_info()
}
