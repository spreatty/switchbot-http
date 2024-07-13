package main

import (
	"log"

	ble "tinygo.org/x/bluetooth"
)

var adapter = ble.DefaultAdapter
var botAddress = parseMac()

func parseMac() ble.Address {
	address, err := ble.ParseMAC(config.BotMac)
	if err != nil {
		log.Fatalln("Failed parsing MAC address", err)
	}
	return ble.Address{MACAddress: ble.MACAddress{MAC: address}}
}

func startBLE() {
	err := adapter.Enable()
	if err != nil {
		log.Fatalln("BLE error", err.Error())
	}
	preapreBLECache()
}

func preapreBLECache() {
	err := adapter.Scan(func(_ *ble.Adapter, res ble.ScanResult) {
		log.Println("Found device", res.Address.String())
		if res.Address != botAddress {
			log.Println("Discarded")
			return
		}
		adapter.StopScan()
		log.Println("BLE cache prepared")
	})
	if err != nil {
		log.Fatalln("Failed scanning")
	}
}
