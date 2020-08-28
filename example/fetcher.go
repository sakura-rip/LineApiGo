// Copyright (c) 2020 @Ch31212y
// Version 1.1 bata
// LastUpdate 2020/08/28

package main

import (
	"fmt"

	"github.com/ch31212y/lineapigo"
)

func main() {
	bot := lineapigo.NewLineClient("LITE")
	bot.LoginWithQrCode()
	for {
		ops, err := bot.FetchOperations()
		if err == nil {
			for _, op := range ops {
				//  skip end of operation
				if op.Revision != -1 {
					fmt.Println(op)
					bot.SetRevision(op.Revision)
				}
			}
		} else {
			fmt.Println(err)
		}
	}
}
