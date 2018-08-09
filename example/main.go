package main

import (
	"fmt"
	"time"

	"github.com/dzen-it/ubm"
)

func main() {
	client := ubm.NewClientInMemory()
	client.SetTrigger("a", "b", time.Second*10)

	fmt.Println("UBM trigger state 1:", client.TriggerStatus("1", "a"))

	client.AddAction("1", "a")

	fmt.Println("UBM trigger state 2:", client.TriggerStatus("1", "a"))

	time.Sleep(time.Second * 10)

	fmt.Println("UBM trigger state 3:", client.TriggerStatus("1", "a"))

	client.AddAction("1", "b")

	fmt.Println("UBM trigger state 4:", client.TriggerStatus("1", "a"))

	time.Sleep(time.Second * 2)

	client.AddAction("1", "b")

	time.Sleep(time.Second * 10)

	fmt.Println("UBM trigger state 5:", client.TriggerStatus("1", "a"))

}
