# LineApiGo
Copyright (c) 2020 @Ch31212y
Version 1.1 bata
LastUpdate 2020/08/28


Unofficial line api

##Usage:
```Go
package main

import "github.com/ch31212y/lineapigo"

func main() {
    // Create Instanct
    bot := lineapigo.NewLineClient("LITE")
    //  Login with Qr Code
    bot.LoginWithQrCode()


    // Create Instanct
    bot1 := lineapigo.NewLineClient("IOS")
    // Login with AuthToken
    bot1.LoginWithAuthToken("token here")
}```
see example to get more imformation


have fun :)

Studying golang...