/* Name: opendevice.c
 * Project: V-USB host-side library
 * Author: Christian Starkjohann
 * Creation Date: 2008-04-10
 * Tabsize: 4
 * Copyright: (c) 2008 by OBJECTIVE DEVELOPMENT Software GmbH
 * License: GNU GPL v2 (see License.txt), GNU GPL v3 or proprietary (CommercialLicense.txt)
 * This Revision: $Id: opendevice.c 740 2009-04-13 18:23:31Z cs $
 */

/*
General Description:
The functions in this module can be used to find and open a device based on
libusb or libusb-win32.
*/

#include <stdio.h>
#include "opendevice.h"

/* ------------------------------------------------------------------------- */

#define MATCH_SUCCESS			1
#define MATCH_FAILED			0
#define MATCH_ABORT				-1

/* private interface: match text and p, return MATCH_SUCCESS, MATCH_FAILED, or MATCH_ABORT. */
static int  _shellStyleMatch(char *text, char *p)
{
int last, matched, reverse;

    for(; *p; text++, p++){
        if(*text == 0 && *p != '*')
            return MATCH_ABORT;
        switch(*p){
        case '\\':
            /* Literal match with following character. */
            p++;
            /* FALLTHROUGH */
        default:
            if(*text != *p)
                return MATCH_FAILED;
            continue;
        case '?':
            /* Match anything. */
            continue;
        case '*':
            while(*++p == '*')
                /* Consecutive stars act just like one. */
                continue;
            if(*p == 0)
                /* Trailing star matches everything. */
                return MATCH_SUCCESS;
            while(*text)
                if((matched = _shellStyleMatch(text++, p)) != MATCH_FAILED)
                    return matched;
            return MATCH_ABORT;
        case '[':
            reverse = p[1] == '^';
            if(reverse) /* Inverted character class. */
                p++;
            matched = MATCH_FAILED;
            if(p[1] == ']' || p[1] == '-')
                if(*++p == *text)
                    matched = MATCH_SUCCESS;
            for(last = *p; *++p && *p != ']'; last = *p)
                if (*p == '-' && p[1] != ']' ? *text <= *++p && *text >= last : *text == *p)
                    matched = MATCH_SUCCESS;
            if(matched == reverse)
                return MATCH_FAILED;
            continue;
        }
    }
    return *text == 0;
}

/* public interface for shell style matching: returns 0 if fails, 1 if matches */
static int shellStyleMatch(char *text, char *pattern)
{
    if(pattern == NULL) /* NULL pattern is synonymous to "*" */
        return 1;
    return _shellStyleMatch(text, pattern) == MATCH_SUCCESS;
}

/* ------------------------------------------------------------------------- */

int usbGetStringAscii(usb_dev_handle *dev, int index, char *buf, int buflen)
{
char    buffer[256];
int     rval, i;

    if((rval = usb_get_string_simple(dev, index, buf, buflen)) >= 0) /* use libusb version if it works */
        return rval;
    if((rval = usb_control_msg(dev, USB_ENDPOINT_IN, USB_REQ_GET_DESCRIPTOR, (USB_DT_STRING << 8) + index, 0x0409, buffer, sizeof(buffer), 5000)) < 0)
        return rval;
    if(buffer[1] != USB_DT_STRING){
        *buf = 0;
        return 0;
    }
    if((unsigned char)buffer[0] < rval)
        rval = (unsigned char)buffer[0];
    rval /= 2;
    /* lossy conversion to ISO Latin1: */
    for(i=1;i<rval;i++){
        if(i > buflen)              /* destination buffer overflow */
            break;
        buf[i-1] = buffer[2 * i];
        if(buffer[2 * i + 1] != 0)  /* outside of ISO Latin1 range */
            buf[i-1] = '?';
    }
    buf[i-1] = 0;
    return i-1;
}

/* ------------------------------------------------------------------------- */

int usbOpenDevice(usb_dev_handle **device, int vendorID, char *vendorNamePattern, int productID, char *productNamePattern, char *serialNamePattern, FILE *printMatchingDevicesFp, FILE *warningsFp)
{
struct usb_bus      *bus;
struct usb_device   *dev;
usb_dev_handle      *handle = NULL;
int                 errorCode = USBOPEN_ERR_NOTFOUND;

    usb_find_busses();
    usb_find_devices();
    for(bus = usb_get_busses(); bus; bus = bus->next){
        for(dev = bus->devices; dev; dev = dev->next){  /* iterate over all devices on all busses */
            if((vendorID == 0 || dev->descriptor.idVendor == vendorID)
                        && (productID == 0 || dev->descriptor.idProduct == productID)){
                char    vendor[256], product[256], serial[256];
                int     len;
                handle = usb_open(dev); /* we need to open the device in order to query strings */
                if(!handle){
                    errorCode = USBOPEN_ERR_ACCESS;
                    if(warningsFp != NULL)
                        fprintf(warningsFp, "Warning: cannot open VID=0x%04x PID=0x%04x: %s\n", dev->descriptor.idVendor, dev->descriptor.idProduct, usb_strerror());
                    continue;
                }
                /* now check whether the names match: */
                len = vendor[0] = 0;
                if(dev->descriptor.iManufacturer > 0){
                    len = usbGetStringAscii(handle, dev->descriptor.iManufacturer, vendor, sizeof(vendor));
                }
                if(len < 0){
                    errorCode = USBOPEN_ERR_ACCESS;
                    if(warningsFp != NULL)
                        fprintf(warningsFp, "Warning: cannot query manufacturer for VID=0x%04x PID=0x%04x: %s\n", dev->descriptor.idVendor, dev->descriptor.idProduct, usb_strerror());
                }else{
                    errorCode = USBOPEN_ERR_NOTFOUND;
                    /* printf("seen device from vendor ->%s<-\n", vendor); */
                    if(shellStyleMatch(vendor, vendorNamePattern)){
                        len = product[0] = 0;
                        if(dev->descriptor.iProduct > 0){
                            len = usbGetStringAscii(handle, dev->descriptor.iProduct, product, sizeof(product));
                        }
                        if(len < 0){
                            errorCode = USBOPEN_ERR_ACCESS;
                            if(warningsFp != NULL)
                                fprintf(warningsFp, "Warning: cannot query product for VID=0x%04x PID=0x%04x: %s\n", dev->descriptor.idVendor, dev->descriptor.idProduct, usb_strerror());
                        }else{
                            errorCode = USBOPEN_ERR_NOTFOUND;
                            /* printf("seen product ->%s<-\n", product); */
                            if(shellStyleMatch(product, productNamePattern)){
                                len = serial[0] = 0;
                                if(dev->descriptor.iSerialNumber > 0){
                                    len = usbGetStringAscii(handle, dev->descriptor.iSerialNumber, serial, sizeof(serial));
                                }
                                if(len < 0){
                                    errorCode = USBOPEN_ERR_ACCESS;
                                    if(warningsFp != NULL)
                                        fprintf(warningsFp, "Warning: cannot query serial for VID=0x%04x PID=0x%04x: %s\n", dev->descriptor.idVendor, dev->descriptor.idProduct, usb_strerror());
                                }
                                if(shellStyleMatch(serial, serialNamePattern)){
                                    if(printMatchingDevicesFp != NULL){
                                        if(serial[0] == 0){
                                            fprintf(printMatchingDevicesFp, "VID=0x%04x PID=0x%04x vendor=\"%s\" product=\"%s\"\n", dev->descriptor.idVendor, dev->descriptor.idProduct, vendor, product);
                                        }else{
                                            fprintf(printMatchingDevicesFp, "VID=0x%04x PID=0x%04x vendor=\"%s\" product=\"%s\" serial=\"%s\"\n", dev->descriptor.idVendor, dev->descriptor.idProduct, vendor, product, serial);
                                        }
                                    }else{
                                        break;
                                    }
                                }
                            }
                        }
                    }
                }
                usb_close(handle);
                handle = NULL;
            }
        }
        if(handle)  /* we have found a deice */
            break;
    }
    if(handle != NULL){
        errorCode = 0;
        *device = handle;
    }
    if(printMatchingDevicesFp != NULL)  /* never return an error for listing only */
        errorCode = 0;
    return errorCode;
}

/* ------------------------------------------------------------------------- */
