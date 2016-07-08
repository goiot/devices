# Go libraries/drivers for IoT devices

[![GoDoc](http://godoc.org/github.com/goiot/devices?status.svg)](http://godoc.org/github.com/goiot/devices)
[![Build Status](https://travis-ci.org/goiot/devices.svg?branch=master)](https://travis-ci.org/goiot/devices)

This repo contains a suite of libraries for IoT devices/sensors/actuators. The suite is meant to be as dependency free
and as idiomatic as possible.

If you are interested in helping, feel free to look at the open issues mentioning help needed.
If you have questions on how you implement some of the features, don't hesitate to ask. If you are trying to integrate
these libraries in your projects and have questions, please open an issue.

Note that there are a LOT of IoT devices and while we would love to have libs for all of them, we will need your help.

## Supported devices

### [Grove](http://www.seeedstudio.com/wiki/Grove_System)

* [3 Axis Digital Accelerometer](https://github.com/goiot/devices/tree/master/accel3xdigital)
* [LCD RGB Backlight](https://github.com/goiot/devices/tree/master/lcdrgbbacklight)
* [OLED 96 x 96](https://github.com/goiot/devices/tree/master/oled96x96)

### [Adafruit](https://www.adafruit.com/)

* [DotStar RGB LED (APA102)](https://github.com/goiot/devices/tree/master/dotstar)
* [Monochrome 0.96" 128x64 OLED graphic display (SSD1306)](https://github.com/goiot/devices/tree/master/ssd1306)

### [Pimoroni](https://shop.pimoroni.com/)

* [PiGlow](https://github.com/goiot/devices/tree/master/piglow)

### Generic

The following libraries can be used on multiple devices that use the same underlying component.
Often various manufacturers create their own version of a device using the same component.
If you have device that doesn't have a driver listed above, look at the main component used and see
it it matches one of the ones mentioned below.

* [APA102 LED strip](https://github.com/goiot/devices/tree/master/dotstar)
* [SSD1306 OLED](https://github.com/goiot/devices/tree/master/ssd1306)

## Repo organization

Device libraries are organized by manufacturers and should use names that easy to google or identify.
Each device package contains a README file with references and details about the device (and usually a picture and links to datasheets). You will also find an examples folder with basic examples on how to use the library.

## Test setup

Testing IoT devices is quite complicated, most of us use a [Raspberry Pi](https://www.raspberrypi.org/), connect the devices
directly or via [shield](http://www.dexterindustries.com/grovepi/) and run the examples to test. Yes, it's far from perfect :(

## More information / Advanced topics

Checkout the [wiki](https://github.com/goiot/devices/wiki) for more info.
