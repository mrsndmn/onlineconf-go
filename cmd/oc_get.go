package main

import (
	"fmt"
	"time"

	"github.com/onlineconf/onlineconf-go"
)

func main() {
	fmt.Println(onlineconf.GetString("/vkpay/host"))
	time.Sleep(2 * time.Second)
	fmt.Println(onlineconf.GetString("/vkpay/host"))
}
