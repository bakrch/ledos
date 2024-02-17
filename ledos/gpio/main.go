package gpio

import "github.com/stianeikeland/go-rpio/v4"

func Main() {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}
}
