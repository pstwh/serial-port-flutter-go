# serial_port_flutter

This Go package implements the host-side of the Flutter [serial_port_flutter](https://github.com/pstwh/serial_port_flutter) plugin.

## Usage

Import as:

```go
import serial_port_flutter "github.com/pstwh/serial_port_flutter/go"
```

Then add the following option to your go-flutter [application options](https://github.com/go-flutter-desktop/go-flutter/wiki/Plugin-info):

```go
flutter.AddPlugin(&serial_port_flutter.SerialPortFlutterPlugin{}),
```
