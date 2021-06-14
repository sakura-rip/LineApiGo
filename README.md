# LineApiGo
Copyright (c) 2020 @sakura-rip\
Version 1.1 bata\
LastUpdate 2020/08/28



Unofficial line api



## Usage:
```Go
package main

import "github.com/sakura-rip/lineapigo"

func main() {
    // Create Instance
    bot := lineapigo.NewLineClient("LITE")
    //  Login with Qr Code
    bot.LoginWithQrCode()


    // Create Instance
    bot1 := lineapigo.NewLineClient("IOS")
    // Login with AuthToken
    bot1.LoginWithAuthToken("token here")
}
```
see example to get more information


have fun :)

Studying golang...

GOたのちい
