# DotStar RGB LED

[![GoDoc](http://godoc.org/github.com/goiot/devices/dotstar?status.png)](http://godoc.org/github.com/goiot/devices/dotstar)

[Manufacturer info](https://www.adafruit.com/product/2238)

DotStar LEDs are 5050-sized LEDs with an embedded microcontroller inside the LED. You can set the color/brightness of each LED to 24-bit color (8 bits each red green and blue). Each LED acts like a shift register, reading incoming color data on the input pins, and then shifting the previous color data out on the output pin. By sending a long string of data, you can control an infinite number of LEDs, just tack on more or cut off unwanted LEDs at the end.

![Adafruit DotStar](https://cdn-shop.adafruit.com/product-videos/320x240/2238-06.jpg)

##Datasheets:

* [APA102 driver Datasheet](https://cdn-shop.adafruit.com/datasheets/APA102.pdf)