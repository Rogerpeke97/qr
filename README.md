## QR Code Generator

This QR code generator uses version 4L, 0 level masking and byte mode encoding. 
More versions, EC levels, auto masking as well as encoding modes, might be added in the future if I'm still interested in this :)

## Run
First you need to import the qr package and you could have an `example.go` file like:
```
package main

import (
    "github.com/Rogerpeke97/qr"
)

func main() {
    // will open the browser to visualize the qr code
    qr.GenQrCodeWithServer()
    //or
    //will return the coordinates to paint the qr code in a canvas or wherever
    qr.GenQrCode()
}
```
And then run `go run example.go`

[Screencast from 05-05-2024 05:14:18 PM.webm](https://github.com/Rogerpeke97/qr/assets/65107071/00e31144-bb38-45f3-83ab-551b9faa0681)
