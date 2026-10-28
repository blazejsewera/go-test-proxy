---
name: implementing-monitor
description: Implementing a new monitor that writes details about the request and/or response to the console.
---

# Implementing a monitor

Implement the `monitor.Monitor` interface.
Place the implementation in the `monitor` package.
Then, wire the new monitor up in `builder.WithMonitor(monitor.Combine([oldMonitors...], newMonitor))`.

Take note that when running `gotestproxy`,
you have to remember to enable the new monitor using the name you set in the `Name` method.

Look at the struct `monitor.request` and its implementation for reference.
Use the `colorfmt` package when writing to the console, if you want to have a pretty, colorful output.
