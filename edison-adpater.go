package jhd1313m

import "os"

func writeFile(path string, data []byte) (i int, err error) {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 200)
	var i = 0
	i, err = file.Read(buf)
	if i == 0 {
		return buf, err
	}
	return buf[:i], err
}

type mux struct {
	pin   int
	value int
}
type sysfsPin struct {
	pin          int
	resistor     int
	levelShifter int
	pwmPin       int
	mux          []mux
}

// Adaptor represents a Gobot Adaptor for an Intel Edison
type Adaptor struct {
	name   string
	board  string
	pinmap map[string]sysfsPin
	//tristate DigitalPin
	//digitalPins map[int]sysfs.DigitalPin
	//pwmPins     map[int]*pwmPin
	i2cDevice I2cDevice
	connect   func(e *Adaptor) (err error)
}

// changePinMode writes pin mode to current_pinmux file
func changePinMode(pin, mode string) (err error) {
	_, err = writeFile(
		"/sys/kernel/debug/gpio_debug/gpio"+pin+"/current_pinmux",
		[]byte("mode"+mode),
	)
	return
}

// NewAdaptor returns a new Edison Adaptor
func NewAdaptor() *Adaptor {
	return &Adaptor{
		name:  "Edison",
		board: "arduino",
		//	pinmap: arduinoPinMap,
	}
}

// I2cStart initializes i2c device for addresss
func (e *Adaptor) I2cStart(address int) (err error) {
	if e.i2cDevice != nil {
		return
	}

	// most board use I2C bus 1
	bus := "/dev/i2c-1"

	// except for Arduino which uses bus 6
	if e.board == "arduino" {
		bus = "/dev/i2c-6"
	}

	e.i2cDevice, err = NewI2cDevice(bus, address)
	return
}

// I2cWrite writes data to i2c device
func (e *Adaptor) I2cWrite(address int, data []byte) (err error) {
	if err = e.i2cDevice.SetAddress(address); err != nil {
		return err
	}
	_, err = e.i2cDevice.Write(data)
	return
}

// I2cRead returns size bytes from the i2c device
func (e *Adaptor) I2cRead(address int, size int) (data []byte, err error) {
	data = make([]byte, size)
	if err = e.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = e.i2cDevice.Read(data)
	return
}

// I2cClose Close i2c device
func (e *Adaptor) I2cClose() (err error) {
	return e.i2cDevice.Close()
}
