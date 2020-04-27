package serial_port_flutter

import (
	flutter "github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"go.bug.st/serial"
	"fmt"
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

	p.stop = make(chan bool)

	eventChannel := plugin.NewEventChannel(messenger, eventChannelName, plugin.StandardMethodCodec{})
	eventChannel.Handle(p)
	
	methodChannel := plugin.NewMethodChannel(messenger, methodChannelName, plugin.StandardMethodCodec{})
	methodChannel.HandleFunc("getPlatformVersion", p.handlePlatformVersion)
	methodChannel.HandleFunc("open", p.openDevice)
	methodChannel.HandleFunc("close", p.closeDevice)
	methodChannel.HandleFunc("write", p.writeDevice)
	methodChannel.HandleFunc("getAllDevices", p.getAllDevices)
	methodChannel.HandleFunc("getAllDevicesPath", p.getAllDevicesPath)
	
	return nil
}

func (p *SerialPortFlutterPlugin) OnListen(arguments interface{}, sink *plugin.EventSink) {
	buff := make([]byte, 100)
	sum := 0;
	for {
		select {
		case <-p.stop:
				return
		default:
			n, _ := p.Port.Read(buff)
			fmt.Println(buff)
			sum += n
			if sum > 100 {
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
	//argsMap := arguments.(map[interface{}]interface{})

	//if p.Port != nil {
		// devicePath := argsMap["devicePath"].(string)
		// if devicePath == "" {
		// 	return nil, nil
		// }
	
	//baudrate := argsMap["baudrate"].(int)

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open("/dev/ttyUSB0", mode)
	if err != nil {
		return false, err
	}

	p.Port = port
	
	return true, nil
	//}

	//return
}

func (p *SerialPortFlutterPlugin) closeDevice(arguments interface{}) (reply interface{}, err error) {
	if p.Port != nil {
		err := p.Port.Close()

		return nil, err
	}

	return true, nil
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
