package jhd1313m

import (
	"errors"
	"fmt"
	"time"
)

const (
	REG_RED   = 0x04
	REG_GREEN = 0x03
	REG_BLUE  = 0x02

	LCD_CLEARDISPLAY        = 0x01
	LCD_RETURNHOME          = 0x02
	LCD_ENTRYMODESET        = 0x04
	LCD_DISPLAYCONTROL      = 0x08
	LCD_CURSORSHIFT         = 0x10
	LCD_FUNCTIONSET         = 0x20
	LCD_SETCGRAMADDR        = 0x40
	LCD_SETDDRAMADDR        = 0x80
	LCD_ENTRYRIGHT          = 0x00
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTINCREMENT = 0x01
	LCD_ENTRYSHIFTDECREMENT = 0x00
	LCD_DISPLAYON           = 0x04
	LCD_DISPLAYOFF          = 0x00
	LCD_CURSORON            = 0x02
	LCD_CURSOROFF           = 0x00
	LCD_BLINKON             = 0x01
	LCD_BLINKOFF            = 0x00
	LCD_DISPLAYMOVE         = 0x08
	LCD_CURSORMOVE          = 0x00
	LCD_MOVERIGHT           = 0x04
	LCD_MOVELEFT            = 0x00
	LCD_2LINE               = 0x08
	LCD_CMD                 = 0x80
	LCD_DATA                = 0x40

	LCD_2NDLINEOFFSET = 0x40
)

var (
	ErrEncryptedBytes  = errors.New("Encrypted bytes")
	ErrNotEnoughBytes  = errors.New("Not enough bytes read")
	ErrNotReady        = errors.New("Device is not ready")
	ErrInvalidPosition = errors.New("Invalid position value")
)

//var _ gobot.Driver = (*JHD1313M1Driver)(nil)

// CustomLCDChars is a map of CGRAM characters that can be loaded
// into a LCD screen to display custom characters. Some LCD screens such
// as the Grove screen (jhd1313m1) isn't loaded with latin 1 characters.
// It's up to the developer to load the set up to 8 custom characters and
// update the input text so the character is swapped by a byte reflecting
// the position of the custom character to use.
// See SetCustomChar
var CustomLCDChars = map[string][8]byte{
	"é":       [8]byte{130, 132, 142, 145, 159, 144, 142, 128},
	"è":       [8]byte{136, 132, 142, 145, 159, 144, 142, 128},
	"ê":       [8]byte{132, 138, 142, 145, 159, 144, 142, 128},
	"à":       [8]byte{136, 134, 128, 142, 145, 147, 141, 128},
	"â":       [8]byte{132, 138, 128, 142, 145, 147, 141, 128},
	"á":       [8]byte{2, 4, 14, 1, 15, 17, 15, 0},
	"î":       [8]byte{132, 138, 128, 140, 132, 132, 142, 128},
	"í":       [8]byte{2, 4, 12, 4, 4, 4, 14, 0},
	"û":       [8]byte{132, 138, 128, 145, 145, 147, 141, 128},
	"ù":       [8]byte{136, 134, 128, 145, 145, 147, 141, 128},
	"ñ":       [8]byte{14, 0, 22, 25, 17, 17, 17, 0},
	"ó":       [8]byte{2, 4, 14, 17, 17, 17, 14, 0},
	"heart":   [8]byte{0, 10, 31, 31, 31, 14, 4, 0},
	"smiley":  [8]byte{0, 0, 10, 0, 0, 17, 14, 0},
	"frowney": [8]byte{0, 0, 10, 0, 0, 0, 14, 17},
}

// JHD1313M1Driver is a driver for the Jhd1313m1 LCD display which has two i2c addreses,
// one belongs to a controller and the other controls solely the backlight.
// This module was tested with the Seed Grove LCD RGB Backlight v2.0 display which requires 5V to operate.
// http://www.seeedstudio.com/wiki/Grove_-_LCD_RGB_Backlight
type JHD1313M1Driver struct {
	name       string
	connection Adaptor
	lcdAddress int
	rgbAddress int
}

// NewJHD1313M1Driver creates a new driver with specified i2c interface.
func NewJHD1313M1Driver(a Adaptor) *JHD1313M1Driver {
	return &JHD1313M1Driver{
		name:       "JHD1313M1",
		connection: a,
		lcdAddress: 0x3E,
		rgbAddress: 0x62,
	}
}

// Name returns the name the JHD1313M1 Driver was given when created.
func (h *JHD1313M1Driver) Name() string { return h.name }

// SetName sets the name for the JHD1313M1 Driver.
func (h *JHD1313M1Driver) SetName(n string) { h.name = n }

/* Connection returns the driver connection to the device.
func (h *JHD1313M1Driver) Connection() gobot.Connection {
	return h.connection.(gobot.Connection)
}
*/

func (h *JHD1313M1Driver) Close() error {
	return h.connection.I2cClose()
}

// Start starts the backlit and the screen and initializes the states.
func (h *JHD1313M1Driver) Start() error {
	if err := h.connection.I2cStart(h.lcdAddress); err != nil {
		return err
	}

	if err := h.connection.I2cStart(h.rgbAddress); err != nil {
		return err
	}

	time.Sleep(50000 * time.Microsecond)
	payload := []byte{LCD_CMD, LCD_FUNCTIONSET | LCD_2LINE}
	if err := h.connection.I2cWrite(h.lcdAddress, payload); err != nil {
		if err := h.connection.I2cWrite(h.lcdAddress, payload); err != nil {
			return err
		}
	}

	time.Sleep(100 * time.Microsecond)
	if err := h.connection.I2cWrite(h.lcdAddress, []byte{LCD_CMD, LCD_DISPLAYCONTROL | LCD_DISPLAYON}); err != nil {
		return err
	}

	time.Sleep(100 * time.Microsecond)
	if err := h.Clear(); err != nil {
		return err
	}

	if err := h.connection.I2cWrite(h.lcdAddress, []byte{LCD_CMD, LCD_ENTRYMODESET | LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT}); err != nil {
		return err
	}

	if err := h.setReg(0, 0); err != nil {
		return err
	}
	if err := h.setReg(1, 0); err != nil {
		return err
	}
	if err := h.setReg(0x08, 0xAA); err != nil {
		return err
	}

	if err := h.SetRGB(255, 255, 255); err != nil {
		return err
	}

	return nil
}

// SetRGB sets the Red Green Blue value of backlit.
func (h *JHD1313M1Driver) SetRGB(r, g, b int) error {
	if err := h.setReg(REG_RED, r); err != nil {
		return err
	}
	if err := h.setReg(REG_GREEN, g); err != nil {
		return err
	}
	return h.setReg(REG_BLUE, b)
}

// Clear clears the text on the lCD display.
func (h *JHD1313M1Driver) Clear() error {
	err := h.command([]byte{LCD_CLEARDISPLAY})
	return err
}

// Home sets the cursor to the origin position on the display.
func (h *JHD1313M1Driver) Home() error {
	err := h.command([]byte{LCD_RETURNHOME})
	// This wait fixes a race condition when calling home and clear back to back.
	time.Sleep(2 * time.Millisecond)
	return err
}

// Write displays the passed message on the screen.
func (h *JHD1313M1Driver) Write(message string) error {
	// This wait fixes an odd bug where the clear function doesn't always work properly.
	time.Sleep(1 * time.Millisecond)
	for _, val := range message {
		if val == '\n' {
			if err := h.SetPosition(16); err != nil {
				return err
			}
			continue
		}
		if err := h.connection.I2cWrite(h.lcdAddress, []byte{LCD_DATA, byte(val)}); err != nil {
			return err
		}
	}
	return nil
}

// SetPosition sets the cursor and the data display to pos.
// 0..15 are the positions in the first display line.
// 16..32 are the positions in the second display line.
func (h *JHD1313M1Driver) SetPosition(pos int) (err error) {
	if pos < 0 || pos > 31 {
		err = ErrInvalidPosition
		return
	}
	offset := byte(pos)
	if pos >= 16 {
		offset -= 16
		offset |= LCD_2NDLINEOFFSET
	}
	err = h.command([]byte{LCD_SETDDRAMADDR | offset})
	return
}

func (h *JHD1313M1Driver) Scroll(leftToRight bool) error {
	if leftToRight {
		return h.connection.I2cWrite(h.lcdAddress, []byte{LCD_CMD, LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVELEFT})
	}

	return h.connection.I2cWrite(h.lcdAddress, []byte{LCD_CMD, LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVERIGHT})
}

// Halt is a noop function.
func (h *JHD1313M1Driver) Halt() error { return nil }

func (h *JHD1313M1Driver) setReg(command int, data int) error {
	return h.connection.I2cWrite(h.rgbAddress, []byte{byte(command), byte(data)})
}

func (h *JHD1313M1Driver) command(buf []byte) error {
	return h.connection.I2cWrite(h.lcdAddress, append([]byte{LCD_CMD}, buf...))
}

// SetCustomChar sets one of the 8 CGRAM locations with a custom character.
// The custom character can be used by writing a byte of value 0 to 7.
// When you are using LCD as 5x8 dots in function set then you can define a total of 8 user defined patterns
// (1 Byte for each row and 8 rows for each pattern).
// Use http://www.8051projects.net/lcd-interfacing/lcd-custom-character.php to create your own
// characters.
// To use a custom character, write byte value of the custom character position as a string after
// having setup the custom character.
func (h *JHD1313M1Driver) SetCustomChar(pos int, charMap [8]byte) error {
	if pos > 7 {
		return fmt.Errorf("can't set a custom character at a position greater than 7")
	}
	location := uint8(pos)
	if err := h.command([]byte{LCD_SETCGRAMADDR | (location << 3)}); err != nil {
		return err
	}

	return h.connection.I2cWrite(h.lcdAddress, append([]byte{LCD_DATA}, charMap[:]...))
}

/*
import (
	"log"
	"sync"
	"time"
)

import "math"

import "errors"

const (
	REG_RED   = 0x04
	REG_GREEN = 0x03
	REG_BLUE  = 0x02

	LCD_CLEARDISPLAY        = 0x01
	LCD_RETURNHOME          = 0x02
	LCD_ENTRYMODESET        = 0x04
	LCD_DISPLAYCONTROL      = 0x08
	LCD_CURSORSHIFT         = 0x10
	LCD_FUNCTIONSET         = 0x20
	LCD_SETCGRAMADDR        = 0x40
	LCD_SETDDRAMADDR        = 0x80
	LCD_ENTRYRIGHT          = 0x00
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTINCREMENT = 0x01
	LCD_ENTRYSHIFTDECREMENT = 0x00
	LCD_DISPLAYON           = 0x04
	LCD_DISPLAYOFF          = 0x00
	LCD_CURSORON            = 0x02
	LCD_CURSOROFF           = 0x00
	LCD_BLINKON             = 0x01
	LCD_BLINKOFF            = 0x00
	LCD_DISPLAYMOVE         = 0x08
	LCD_CURSORMOVE          = 0x00
	LCD_MOVERIGHT           = 0x04
	LCD_MOVELEFT            = 0x00
	LCD_2LINE               = 0x08
	LCD_CMD                 = 0x80
	LCD_DATA                = 0x40

	LCD_2NDLINEOFFSET = 0x40
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var ACCUMULATION_ENERGY_ACTIVE = []byte{0xA5, 0x08, 0x41, 0x00, 0x1E, 0x4E, 0x08, 0x62}
var ACTIVE_POWER = []byte{0xA5, 0x08, 0x41, 0x00, 0x12, 0x4E, 0x04, 0x52}
var VOLTAGE_RMS = []byte{0xA5, 0x08, 0x41, 0x00, 0x06, 0x4E, 0x02, 0x44}
var FREQUENCY = []byte{0xA5, 0x08, 0x41, 0x00, 0x08, 0x4E, 0x02, 0x46}

var ACTIVE_ACCUMULATION_ENERGY = []byte{0xA5, 0x0A, 0x41, 0X00, 0xDC, 0X4D, 0X02, 0x01, 0X00, 0x1C}

var SAVE_TO_FLASH = []byte{0xA5, 0x04, 0x53, 0xFC}

type MCP39F521 struct {
	i2c *I2C
	sync.Mutex
}

func New() (*MCP39F521, error) {
	i, err := NewI2c(0x74, 1)
	check(err)
	return &MCP39F521{i2c: i}, err
}

func (mcp *MCP39F521) Close() {
	mcp.i2c.Close()
}

/*
ACTIVE ENERGY ACCUMULATION

func (mcp *MCP39F521) ActiveEnergyAccumulation() error {
	_, err := mcp.getVALUE(ACTIVE_ACCUMULATION_ENERGY, 10)
	if err != nil {
		return err
	}
	_, err = mcp.getVALUE(SAVE_TO_FLASH, 10)

	return err
}

/*
Retourne l'energie Active

func (mcp *MCP39F521) GetACTIVE_POWER() (float64, error) {
	val, err := mcp.getVALUE(ACTIVE_POWER, 10)
	return val, err
}

/*
Retourne l'Energie Active Cumulée

func (mcp *MCP39F521) GetACCUMULATION_ENERGY_ACTIVE() (float64, error) {
	val, err := mcp.getVALUE(ACCUMULATION_ENERGY_ACTIVE, 13)
	return val, err
}

/*
Retourne le voltage

func (mcp *MCP39F521) GetVOLTAGE_RMS() (float64, error) {
	val, err := mcp.getVALUE(VOLTAGE_RMS, 10)
	return val, err
}

/*
Retourne la Frequence

func (mcp *MCP39F521) GetFREQUENCY() (float64, error) {
	val, err := mcp.getVALUE(FREQUENCY, 10)
	return val, err
}

/*
Retourne une valeur Active seule

func (mcp *MCP39F521) getVALUE(request []byte, nbval int) (float64, error) {
	mcp.Lock()
	defer mcp.Unlock()

	p := make([]byte, nbval)
	mcp.i2c.Write(request)
	time.Sleep(9 * time.Millisecond)
	_, err := mcp.i2c.Read(p)
	test := 0
	for err != nil {
		_, err = mcp.i2c.Read(p)
		if test > 10000 {
			return 0.0, errors.New("Too much try")
		}
		test++
	}
	err, val := convertToValue(p)
	return val, err
}

/*
Fonction qui convertir le retour en valeur float
Ex: de retour
ACK, Nb de bit, value, checksum
[Ox06, 0x03, 0x00, 0x50, 0xFF]
Return 0x5000

func convertToValue(b []byte) (error, float64) {
	if b[0] == 6 || b[0] == 134 {
		sum := 0.0
		powerF := 0
		if int(b[1]) > len(b) {
			return errors.New("Size too long"), float64(0.0)
		}
		for i := 2; i < int(b[1])-1; i++ {
			sum = sum + float64(int(b[i]))*math.Pow(16, float64(powerF))
			powerF = powerF + 2
		}
		return nil, sum
	}
	return errors.New("NOK"), 0.0
}
*/
