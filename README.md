# Switchbot HTTP

## Credits

Huge thanks to https://github.com/OpenWonderLabs/SwitchBotAPI-BLE

## About

A service that allows for simple and fast controlling of a Switchbot Bot over HTTP. The Bot must stay in pressing mode. Use Switchbot app or NRF connect app to find out your bot's MAC address.

If you need to control few bots, run few instances of this service on different ports.

Designed to run on the same network with downstream service, thus TLS support is absent.

The service doesn't provide actual reply from Switchbot. The reason is there are issues with WinRT callbacks.

## Build and run

Run like normal Go project.
```
go run .
```

Build like normal Go project.
```
go build .
```

## Endpoints

- POST /press - press, i.e. push and return back
- POST /on - push and stay in that position
- POST /off - just return back
- POST /open - connect to the bot and enable keep-alive mode
- POST /close - disable keep-alive mode

## Keep-alive mode

An action requires going through following steps:
1. Connect to the bot
2. Read services
3. Read characteristics
4. Perform an action

This may take some time and may not be suitable for certain apps. The keep-alive mode will keep the bot awake by reading it status every 1.5 minutes. In keep-alive mode an action is performed in one step and the bot acts instantly.

## Battery considerations

The keep-alive mode will drain the battery significantly faster. To keep the bot operational I recommend to replace the battery with a 5V to 3V USB adapter connected to a PC or 5V power adapter. Example: https://www.aliexpress.com/item/1005007018356955.html