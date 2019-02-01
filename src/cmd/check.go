/*
 * @Description:
 * @Author: EagleXiang
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-01-02 12:42:49
 * @LastEditors: EagleXiang
 * @LastEditTime: 2019-01-31 21:11:36
 */

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	cipher "github.com/eaglexiang/go-cipher"
	simplecipher "github.com/eaglexiang/go-simplecipher"

	"../etcore"
	myet "../etcore/et"
)

// Check check命令
func Check(args []string) {
	if len(args) < 3 {
		fmt.Println("no cmd for et-check")
		return
	}
	theType := myet.ParseEtCheckType(args[2])
	switch theType {
	case myet.EtCheckPING:
		ping()
	case myet.EtCheckAuth:
		auth()
	case myet.EtCheckVERSION:
		version()
	default:
		fmt.Println("invalid check command")
	}
	os.Exit(0)
}

// ping 发送Ping请求并打印结果
func ping() {
	et := createET()
	var time int
	var success int
	timeSig := make(chan string)
	for i := 0; i < 10; i++ {
		go myet.SendEtCheckPingReq(et, timeSig)
	}
	for i := 0; i < 10; i++ {
		timeStr := <-timeSig
		theTime, err := strconv.ParseInt(timeStr, 10, 32)
		if err != nil {
			fmt.Println("PING ", err.Error())
		} else {
			time += int(theTime)
			success++
			fmt.Println("PING ", theTime, "ms")
		}
	}
	fmt.Println("--- ping statistics ---")
	fmt.Println("10 messages transmitted, ", success, " received, ", (10-success)*10, "% loss")
}

func auth() {
	et := createET()
	reply := myet.SendEtCheckAuthReq(et)
	fmt.Println(reply)
}

func version() {
	et := createET()
	reply := myet.SendEtCheckVersionReq(et)
	fmt.Println(reply)
}

func createET() *myet.ET {
	cipher.DefaultCipher = func() cipher.Cipher {
		cipherType := cipher.ParseCipherType(etcore.ConfigKeyValues["cipher"])
		switch cipherType {
		case cipher.SimpleCipherType:
			c := simplecipher.SimpleCipher{}
			c.SetKey(etcore.ConfigKeyValues["data-key"])
			return &c
		default:
			return nil
		}
	}
	et := myet.CreateET(
		myet.ProxyENABLE,
		etcore.ConfigKeyValues["ip-type"],
		etcore.ConfigKeyValues["head"],
		etcore.ConfigKeyValues["relayer"],
		etcore.ConfigKeyValues["location"],
		etcore.LocalUser,
		etcore.Users,
		time.Second*time.Duration(etcore.Timeout),
	)
	return et
}
