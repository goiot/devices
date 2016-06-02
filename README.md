# Golang libraries/drivers for specific IoT devices

This repo contains a suite of libraries for IoT devices that are meant to be as dependency free 
and as idiomatic as possible.

These libraries rely on Go's new exp/io interfaces and are designed so they can be used as is or via 
existing Go IoT libraries/frameworks such as [GoBot](https://gobot.io/), [Embed](http://embd.kidoman.io/) and [hwio](https://github.com/mrmorphic/hwio).

If you are interested in helping, feel free to look at the open issues mentioning help needed.
If you have questions on how you implement some of the features, don't hesitate to ask. If you are trying to integrate
these libraries in your projects and have questions, please open an issue.

Note that there are a LOT of IoT devices and while we would love to have libs for all of them, we will need your help.

## Repo organization

Device libraries are organized by manufacturers and should use names that easy to google or identify.
Each device package contains a README file with references and details about the device (and usually a picture and links to datasheets). You will also find an examples folder with basic examples on how to use the library.

* **Grove** Refers to [Seedstudio's Grove](http://www.seeedstudio.com/wiki/Grove_System), a very user friendly collection of modules
using a common connector type. This is often easier to use that having to connect to a [breadboard](https://en.wikipedia.org/wiki/Breadboard).

## Test setup

Testing IoT devices is quite complicated, most of us use a [Raspberry Pi](https://www.raspberrypi.org/), connect the devices
directly or via [shield](http://www.dexterindustries.com/grovepi/) and run the examples to test. Yes, it's far from perfect :(
