package main

import (
	"errors"
	"log"
	"time"

	ble "tinygo.org/x/bluetooth"
)

const (
	StatusWriteError = iota
	StatusOK
	StatusError
	StatusBusy
	StatusVersionIncompatible
	StatusUnsupportedCommand
	StatusLowBattery
	StatusDeviceEncrypted
	StatusDeviceUnencrypted
	StatusPasswordError
	StatusUnsupportedEncryption
	StatusNoNearbyDevice
	StatusNoNetwork
)

type BotOpenOptions struct {
	ConnectTries                 int
	DiscoverServiceTries         int
	DiscoverCharacteristicsTries int
}

const retryCount = 3

var serviceFilter = makeServiceFilter("cba20d00-224d-11e6-9fb8-0002a5d5c51b")

var (
	ActionPress   = []byte{0x57, 0x01, 0x0}
	ActionPush    = []byte{0x57, 0x01, 0x3}
	ActionReturn  = []byte{0x57, 0x01, 0x4}
	ActionGetInfo = []byte{0x57, 0x02, 0x0}
)

const pingInterval = time.Minute + time.Second*30

type Bot struct {
	address    ble.Address
	device     *ble.Device
	service    *ble.DeviceService
	writeChar  *ble.DeviceCharacteristic
	notifyChar *ble.DeviceCharacteristic
	keepAlive  bool
	pingTicker *time.Ticker
}

func makeServiceFilter(serviceRawUUID string) []ble.UUID {
	serviceUUID, _ := ble.ParseUUID(serviceRawUUID)
	return []ble.UUID{serviceUUID}
}

func (bot *Bot) Press() (int, error) {
	log.Println("Pressing", bot.address)
	if bot.keepAlive {
		bot.pingTicker.Reset(pingInterval)
	}
	return bot.do(ActionPress)
}

func (bot *Bot) On() (int, error) {
	log.Println("Pushing", bot.address)
	if bot.keepAlive {
		bot.pingTicker.Reset(pingInterval)
	}
	return bot.do(ActionPush)
}

func (bot *Bot) Off() (int, error) {
	log.Println("Returning", bot.address)
	if bot.keepAlive {
		bot.pingTicker.Reset(pingInterval)
	}
	return bot.do(ActionReturn)
}

func (bot *Bot) do(action []byte) (int, error) {
	n, err := bot.writeChar.Write(action)
	if err != nil {
		log.Println("Failed writing characteristic", bot.address, err)
		return n, err
	}
	log.Println("Writing done. Code:", n, bot.address, err)
	return n, nil
}

func (bot *Bot) Open() error {
	if bot.keepAlive {
		log.Println("Bot connection is being kept alive", bot.address)
		return nil
	}

	if bot.writeChar != nil {
		log.Println("Checking if bot still connected", bot.address)
		_, err := bot.do(ActionGetInfo)
		if err == nil {
			log.Println("Bot still connected", bot.address)
			return nil
		}
		log.Println("Bot terminated connection, reconnecting", bot.address)
		bot.device.Disconnect()
	}

	err := bot.connect(retryCount)
	if err != nil {
		return err
	}

	err = bot.getService(retryCount)
	if err != nil {
		bot.device.Disconnect()
		return err
	}

	err = bot.getChars(retryCount)
	if err != nil {
		bot.device.Disconnect()
		return err
	}

	log.Println("Bot opened", bot.address)
	return nil
}

func (bot *Bot) StartKeepAlive() {
	if bot.keepAlive {
		return
	}
	bot.keepAlive = true
	bot.pingTicker = time.NewTicker(pingInterval)
	log.Println("Bot started keep-alive", bot.address)
	go func() {
		for range bot.pingTicker.C {
			log.Println("Ping", bot.address)
			bot.do(ActionGetInfo)
		}
	}()
}

func (bot *Bot) StopKeepAlive() {
	if !bot.keepAlive {
		return
	}
	bot.keepAlive = false
	bot.pingTicker.Stop()
	bot.pingTicker = nil
	log.Println("Bot stopped keep-alive", bot.address)
}

func (bot *Bot) connect(tries int) error {
	for try := 1; try <= tries; try++ {
		log.Println("Trying to connect", try, bot.address)
		device, err := adapter.Connect(bot.address, ble.ConnectionParams{})
		if err == nil {
			bot.device = device
			log.Println("Connected", bot.address)
			return nil
		}
		log.Println("Connecting error", bot.address, err.Error())
	}
	log.Println("Failed connecting", bot.address)
	return errors.New("failed connecting")
}

func (bot *Bot) getService(tries int) error {
	for try := 1; try <= tries; try++ {
		log.Println("Trying to discover services", try, bot.address)
		services, err := bot.device.DiscoverServices(serviceFilter)
		if err == nil {
			bot.service = &services[0]
			log.Println("Service discovered", bot.address)
			return nil
		}
		log.Println("Service discovering error", bot.address, err.Error())
	}
	log.Println("Failed discovering services", bot.address)
	return errors.New("failed discovering services")
}

func (bot *Bot) getChars(tries int) error {
	for try := 1; try <= tries; try++ {
		log.Println("Trying to discover characteristics", try, bot.address)
		chars, err := bot.service.DiscoverCharacteristics(nil)
		if err == nil && len(chars) == 2 {
			bot.notifyChar = &chars[0]
			bot.writeChar = &chars[1]
			log.Println("Characteristics discovered", bot.address)
			return nil
		}
		if err != nil {
			log.Println("Characteristics discovering error", bot.address, err.Error())
		} else {
			log.Println("Got 0 characteristics", bot.address)
		}
	}
	log.Println("Failed discovering characteristics", bot.address)
	return errors.New("failed discovering characteristics")
}
