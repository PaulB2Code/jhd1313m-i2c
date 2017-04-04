package jhd1313m

import (
	"fmt"
	"os"
	"sync"
	"syscall"
)

//https://github.com/hybridgroup/gobot/blob/01cc6f73d6827cc508354e288057c28c804771bd/platforms/intel-iot/joule/joule_adaptor.go
//https://github.com/hybridgroup/gobot/blob/master/sysfs/i2c_device.go#L44

/*
Source From David Cheney
Date 25/07/2016
*/

const (
	i2c_SLAVE = 0x0703
)

// I2C represents a connection to an i2c device.
type I2C struct {
	rc *os.File
	sync.Mutex
}

// New opens a connection to an i2c device.
func NewI2c(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), i2c_SLAVE, uintptr(addr)); err != nil {
		return nil, err
	}
	return &I2C{rc: f}, nil
}

// Write sends buf to the remote i2c device. The interpretation of
// the message is implementation dependant.
func (i2c *I2C) Write(buf []byte) (int, error) {
	i2c.Lock()
	defer i2c.Unlock()
	return i2c.rc.Write(buf)
}

func (i2c *I2C) WriteByte(b byte) (int, error) {
	i2c.Lock()
	defer i2c.Unlock()

	var buf [1]byte
	buf[0] = b
	return i2c.rc.Write(buf[:])
}

func (i2c *I2C) Read(p []byte) (int, error) {
	i2c.Lock()
	defer i2c.Unlock()

	return i2c.rc.Read(p)
}

func (i2c *I2C) Close() error {
	i2c.Lock()
	defer i2c.Unlock()

	return i2c.rc.Close()
}

func ioctl(fd, cmd, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}
