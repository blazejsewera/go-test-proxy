---
name: implementing-monitor
description: Implementing a new monitor that writes details about the request and/or response to the console.
---

# Implementing a monitor

Implement the `monitor.Monitor` interface.
Place the implementation in the `monitor` package.
Then, wire the new monitor up in `builder.WithMonitor(monitor.Combine(/* oldMonitors... */, newMonitor))`.

A monitor should always write to `stdout` and not to `stderr`.
Look at the struct `monitor.request` and its implementation for reference.
Use `*colorfmt.Fmt` when writing to `stdout`.
Inject the `*colorfmt.Fmt` dependency in `main.go` from `cfg`, just like `monitor.NewRequestMonitor`.

Take note that when running `gotestproxy`,
you have to remember to enable the new monitor with the `-monitors` arg,
using the name you set in the `Name` method.
