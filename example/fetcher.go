// Copyright (c) 2020 @sakura-rip
// Version 1.1 beta
// LastUpdate 2020/08/28

package main

import (
	"fmt"

	"github.com/sakura-rip/lineapigo"
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
