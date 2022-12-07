package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

// autoregisters driver

func main() {
	fmt.Println("ごあいさつ")
	defer midi.CloseDriver()

	// 入力ポート
	inports := midi.GetInPorts()
	vt4_in_num := -1
	for i, v := range inports {
		// fmt.Println(v.String())
		if strings.Contains(v.String(), "VT-4") {
			vt4_in_num = i
		}
	}
	if vt4_in_num < 0 {
		fmt.Println("[INFO]VT-4 MIDI input がないです")
		return
	}
	vt4_in, err := midi.InPort(vt4_in_num)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("VT-4 in port num: ", vt4_in_num)

	// 出力ポート
	outports := midi.GetOutPorts()
	vt4_out_num := -1
	for i, v := range outports {
		if strings.Contains(v.String(), "VT-4") {
			vt4_out_num = i
		}
	}
	if vt4_out_num < 0 {
		fmt.Println("[INFO]VT-4 MIDI output がないです")
		return
	}
	vt4_out, err2 := midi.OutPort(vt4_out_num)
	if err2 != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("VT-4 out port num: ", vt4_out_num)

	vt4_out.Open()
	defer vt4_out.Close()

	sysex, _ := ioutil.ReadFile("VT-4_REQ_ID.syx")
	vt4_out.Send(sysex)

	vt4_in.Open()
	defer vt4_in.Close()

	var config drivers.ListenConfig
	config.SysEx = true
	stopFn, _ := vt4_in.Listen(onMsg, config)
	defer stopFn()
	start_tm := time.Now()
	for {
		time.Sleep(time.Millisecond)
		if time.Since(start_tm).Seconds() > 5 {
			break
		}
	}
	fmt.Println("ありがとうございました")
}

func onMsg(msg []byte, msec int32) {
	for _, v := range msg {
		fmt.Printf("%02X", v)
	}
	fmt.Println()
}
