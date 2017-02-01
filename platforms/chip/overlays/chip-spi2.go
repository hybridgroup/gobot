package overlays

const SPI2Overlay = `

/*
 * Copyright 2016, Robert Wolterman
 * This file is an amalgamation of stuff from Kolja Windeler, Maxime Ripard, and Renzo.
 *
 * This file is dual-licensed: you can use it either under the terms
 * of the GPL or the X11 license, at your option. Note that this dual
 * licensing only applies to this file, and not this project as a
 * whole.
 *
 *  a) This file is free software; you can redistribute it and/or
 *     modify it under the terms of the GNU General Public License as
 *     published by the Free Software Foundation; either version 2 of the
 *     License, or (at your option) any later version.
 *
 *     This file is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 * Or, alternatively,
 *
 *  b) Permission is hereby granted, free of charge, to any person
 *     obtaining a copy of this software and associated documentation
 *     files (the "Software"), to deal in the Software without
 *     restriction, including without limitation the rights to use,
 *     copy, modify, merge, publish, distribute, sublicense, and/or
 *     sell copies of the Software, and to permit persons to whom the
 *     Software is furnished to do so, subject to the following
 *     conditions:
 *
 *     The above copyright notice and this permission notice shall be
 *     included in all copies or substantial portions of the Software.
 *
 *     THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 *     EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 *     OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 *     NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 *     HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
 *     WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 *     FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *     OTHER DEALINGS IN THE SOFTWARE.
 */

/dts-v1/;
/plugin/;

/ {
    compatible = "nextthing,chip", "allwinner,sun5i-r8";

    /* activate the gpio for interrupt */
    fragment@0 {
        target-path = <&pio>;

        __overlay__ {
            chip_spi2_pins: spi2@0 {
                allwinner,pins = "PE1", "PE2", "PE3";
                allwinner,function = "spi2";
                allwinner,drive = "0"; //<SUN4I_PINCTRL_10_MA>;
                allwinner,pull = "0"; //<SUN4I_PINCTRL_NO_PULL>;
            };

            chip_spi2_cs0_pins: spi2_cs0@0 {
                allwinner,pins = "PE0";
                allwinner,function = "spi2";
                allwinner,drive = "0"; //<SUN4I_PINCTRL_10_MA>;
                allwinner,pull = "0"; //<SUN4I_PINCTRL_NO_PULL>;
            };
        };
    };

    /*
     * Enable our SPI device, with an spidev device connected
     * to it
    */
    fragment@1 {
        target = <&spi2>;

        __overlay__ {
            #address-cells = <1>;
            #size-cells = <0>;
            pinctrl-names = "default";
            pinctrl-0 = <&chip_spi2_pins>, <&chip_spi2_cs0_pins>;
            status = "okay";

            spi2@0 {
                compatible = "rohm,dh2228fv";
                reg = <0>;
                spi-max-frequency = <24000000>;
            };
        };
    };
};

`