package gobotDigispark

/*
#include <opendevice.h>
*/
//import "C"

/*
//int usbGetStringAscii(usb_dev_handle *dev, int index, char *buf, int buflen);
func usbGetStringAscii(dev *usb_dev_handle, index int, buf *int8, buflen int) int {
  return int(C.usbGetStringAscii(dev, index, buf, buflen))
}

//int usbOpenDevice(usb_dev_handle **device, int vendorID, char *vendorNamePattern, int productID, char *productNamePattern, char *serialNamePattern, FILE *printMatchingDevicesFp, FILE *warningsFp);
func usbOpenDevice(device **usb_dev_handle, vendorID int, vendorNamePattern *int8, productID int, productNamePattern *int8, serialNamePattern *int8, printMatchingDevicesFp *FILE, warningsFp *FILE) int {
  return int(usbOpenDevice(device, vendorID, vendorNamePattern, productID, productNamePattern, serialNamePattern, printMatchingDevicesFp, warningsFp))
}
*/
