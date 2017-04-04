package jhd1313m

import (
	"fmt"
	"testing"
	"time"
)

func TestDataCollection(t *testing.T) {
	fmt.Println("Start Programme")
	mcp39F521 := NewJHD1313M1Driver(Adaptor{})
	mcp39F521.Start()
	mcp39F521.SetRGB(40, 100, 10)

	mcp39F521.Write("PaulB2Code \n say hello to you!!")
	time.Sleep(4 * time.Second)

	mcp39F521.SetRGB(40, 100, 100)
	mcp39F521.Clear()
	time.Sleep(4 * time.Second)

	mcp39F521.Write("and wish you \n a good use!!")
	time.Sleep(4 * time.Second)

}
