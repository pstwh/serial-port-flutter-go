package serial_port_flutter

import (
	"go.bug.st/serial"
	"github.com/pkg/errors"

	flutter "github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
)

const methodChannelName = "serial_port"
const eventChannelName = methodChannelName+"/event"

// SerialPortFlutterPlugin implements flutter.Plugin and handles method.
type SerialPortFlutterPlugin struct{
	Port serial.Port
	stop chan bool
}

var _ flutter.Plugin = &SerialPortFlutterPlugin{} // compile-time type check

// InitPlugin initializes the plugin.
func (p *SerialPortFlutterPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	
	methodChannel := plugin.NewMethodChannel(messenger, methodChannelName, plugin.StandardMethodCodec{})
	methodChannel.HandleFunc("getPlatformVersion", p.handlePlatformVersion)
	methodChannel.HandleFunc("open", p.openDevice)
	methodChannel.HandleFunc("close", p.closeDevice)
	methodChannel.HandleFunc("write", p.writeDevice)
	methodChannel.HandleFunc("getAllDevices", p.getAllDevices)
	methodChannel.HandleFunc("getAllDevicesPath", p.getAllDevicesPath)

	p.stop = make(chan bool)

	eventChannel := plugin.NewEventChannel(messenger, eventChannelName, plugin.StandardMethodCodec{})
	eventChannel.Handle(p)
	
	return nil
}

func (p *SerialPortFlutterPlugin) OnListen(arguments interface{}, sink *plugin.EventSink) {
	buff := make([]byte, 128)
	
	for {
		select {
		case <-p.stop:
				return
		default:
			n, err := p.Port.Read(buff)
			
			if n == 0 || err != nil {
				sink.EndOfStream()
				break
			}
		
			sink.Success(buff)
		}
	}
}

func (p *SerialPortFlutterPlugin) OnCancel(arguments interface{}) {
	p.stop <- true
}

func (p *SerialPortFlutterPlugin) handlePlatformVersion(arguments interface{}) (reply interface{}, err error) {
	return "go-flutter " + flutter.PlatformVersion, nil
}

func (p *SerialPortFlutterPlugin) openDevice(arguments interface{}) (reply interface{}, err error) {
	argsMap := arguments.(map[interface{}]interface{})

	devicePath := argsMap["devicePath"].(string)
	if devicePath == "" {
		return false, errors.New("Device Path could not be null")
	}
	
	baudrate := int(argsMap["baudrate"].(int32))
	if baudrate == -1 {
		return false, errors.New("Baud Rate could not be null")
	}

	mode := &serial.Mode{
		BaudRate: baudrate,
	}

	port, err := serial.Open(devicePath, mode)
	if err != nil {
		return false, err
	}

	p.Port = port
	
	return true, nil
}

func (p *SerialPortFlutterPlugin) closeDevice(arguments interface{}) (reply interface{}, err error) {
	if p.Port != nil {
		err := p.Port.Close()
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return
}

func (p *SerialPortFlutterPlugin) writeDevice(arguments interface{}) (reply interface{}, err error) {
	argsMap := arguments.(map[interface{}]interface{})

	data := argsMap["data"].([]byte)

	_, err = p.Port.Write(data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *SerialPortFlutterPlugin) getAllDevices(arguments interface{}) (reply interface{}, err error) {
	ports, err := serial.GetPortsList()
	
	y := make([]interface{}, len(ports))
    for i, v := range ports {
        y[i] = v
    }

	return y, err;
}

func (p *SerialPortFlutterPlugin) getAllDevicesPath(arguments interface{}) (reply interface{}, err error) {
	return "go-flutter " + flutter.PlatformVersion, nil
}
