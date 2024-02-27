package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/rivo/tview"
	"github.com/tarm/serial"
)

const (
	comNum     = "COM101"
	cmdHeader  = 0xC0
	cmdTailer  = 0xC1
	recvMaxLen = 1024
)

var (
	serialPort *serial.Port
	app        *tview.Application
	grid       *tview.Grid
	menu       *tview.List
	resultText *tview.Form
	mainpage   = newPrimitive("Main content")
	sideBar    = newPrimitive("Result")
	footer     = newPrimitive("Uart Log")

	appFormArr = [][]FormData{
		Form0,
		Form1,
		Form2,
		Form3,
		Form4,
		Form5,
		Form6,
		Form1,
	}

	appTextViewArr = []string{
		"Header",
		"Len",
		"Type",
		"ID Len",
		"ID",
		"Cmd",
		"Data",
		"CRC",
		"Tail",
	}

	// appMenuArr = []string{
	// 	"Get device id",
	// 	"Set device id",
	// 	"Get gpio value",
	// 	"Set gpio value",
	// 	"Get uart data",
	// 	"Send uart data",
	// 	"Get adc value",
	// }
)

func crc16(pData []byte, size int) uint16 {
	crc := uint16(0xFFFF)
	var i, j int

	for i = 0; i < size; i++ {
		crc ^= uint16(pData[i])
		for j = 0; j < 8; j++ {
			if (crc & 0x0001) > 0 {
				crc = (crc >> 1) ^ 0x8408
			} else {
				crc = crc >> 1
			}
		}
	}

	return crc
}

func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}

func receiveData(serialPort *serial.Port) {
	recvData := make([]byte, 0, recvMaxLen)
	buffer := make([]byte, 1)
	for {
		_, err := serialPort.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}

		if buffer[0] == 0xC1 {
			recvData = append(recvData, buffer[0])
			uartStr := ""
			for _, v := range recvData {
				uartStr += fmt.Sprintf("%02x ", v)
			}
			grid.AddItem(newPrimitive(uartStr), 2, 0, 1, 3, 0, 0, false)

			buffer[0] = 0

			d, err := ProtocolAnalysis(&recvData)
			if err != nil {
				fmt.Println(err)
			}

			t := reflect.TypeOf(d)
			var printStr string
			for i := 0; i < t.NumField(); i++ {
				fieldValue := reflect.ValueOf(d).Field(i).Interface()
				switch value := fieldValue.(type) {
				case byte:
					printStr = fmt.Sprintf("%02x", value)
				case uint16:
					printStr = fmt.Sprintf("%02x %02x", (value>>8)&0xFF, value&0xFF)
				case uint32:
					printStr = fmt.Sprintf("%02x %02x %02x %02x", (value>>24)&0xFF, (value>>16)&0xFF, (value>>8)&0xFF, value&0xFF)
				case []byte:
					printStr = ""
					for _, byteValue := range value {
						printStr += fmt.Sprintf("%02x ", byteValue)
					}
				}
				resultText.GetFormItem(i).(*tview.TextView).SetText(printStr)
			}

			grid.AddItem(resultText, 1, 2, 1, 1, 0, 0, false)
			app.ForceDraw()
			recvData = nil
			recvData = make([]byte, 0, recvMaxLen)
		} else {
			if len(recvData) > recvMaxLen {
				err = errors.New("beyond recv max len")
				fmt.Println(err)
				recvData = nil
				recvData = make([]byte, 0, recvMaxLen)
				continue
			}
			recvData = append(recvData, buffer[0])
		}
	}
}

func listCMD(formNum int) {
	var sendData []byte
	var err error
	form1 := tview.NewForm()
	for form1ArrI, form1ArrV := range appFormArr[formNum] {
		if form1ArrI == 9 {
			form1.AddTextView(form1ArrV.name, form1ArrV.defaultValue, 40, 3, true, false)
		} else if form1ArrI == 0 || form1ArrI == 1 || form1ArrI == 7 || form1ArrI == 8 {
			form1.AddTextView(form1ArrV.name, form1ArrV.defaultValue, 0, 1, true, false)
		} else {
			form1.AddInputField(form1ArrV.name, form1ArrV.defaultValue, 20, nil, nil)
		}
	}

	form1.AddButton("Send", func() {
		var errStr string
		var result tview.Primitive
		var sendStr string
		for i := 0; i < len(appFormArr[formNum])-3; i++ {
			if i == 0 || i == 1 {
				sendStr += form1.GetFormItem(i).(*tview.TextView).GetText(false)
			} else {
				sendStr += form1.GetFormItem(i).(*tview.InputField).GetText()
			}
		}

		sendStr = strings.ReplaceAll(sendStr, " ", "")
		sendData, err = hex.DecodeString(sendStr)
		if err != nil {
			errStr = fmt.Sprintf("Error: %s", err.Error())
		} else {
			sendStrLen := uint16(len(sendData))
			sendData[1] = byte(sendStrLen >> 8 & 0xFF)
			sendData[2] = byte(sendStrLen & 0xFF)
			sendStrLenStr := fmt.Sprintf("%02x %02x", sendStrLen>>8&0xFF, sendStrLen&0xFF)
			form1.GetFormItem(1).(*tview.TextView).SetText(sendStrLenStr)

			crcData := crc16(sendData, len(sendData))
			sendData = append(sendData, byte(crcData>>8&0xFF))
			sendData = append(sendData, byte(crcData&0xFF))
			crcDataStr := fmt.Sprintf("%02x %02x", byte(crcData>>8&0xFF), byte(crcData&0xFF))
			form1.GetFormItem(7).(*tview.TextView).SetText(crcDataStr)
			sendData = append(sendData, 0xC1)

			sendStr = ""
			for i := 0; i < len(appFormArr[formNum])-1; i++ {
				if i == 0 || i == 1 || i == 7 || i == 8 {
					sendStr += form1.GetFormItem(i).(*tview.TextView).GetText(false)
				} else {
					sendStr += form1.GetFormItem(i).(*tview.InputField).GetText()
				}
			}
			sendStr = strings.ReplaceAll(sendStr, " ", "")

			form1.GetFormItem(9).(*tview.TextView).SetText(sendStr)

			if _, err := serialPort.Write(sendData); err != nil {
				errStr = fmt.Sprintf("Error: %s", err.Error())
			}
		}

		if len(errStr) > 0 {
			result = newPrimitive(errStr)
			grid.AddItem(result, 1, 1, 1, 1, 0, 0, false)
		}
	}).
		AddButton("Back", func() {
			// app.Stop()
			app.SetFocus(menu)
			grid.AddItem(mainpage, 1, 1, 1, 1, 0, 0, false).
				AddItem(sideBar, 1, 2, 1, 1, 0, 0, false).
				AddItem(footer, 2, 0, 1, 3, 0, 0, false)
		})

	grid.AddItem(form1, 1, 1, 1, 1, 0, 0, false)

	app.SetFocus(form1)
}

func main() {
	var err error

	// Test()
	// return

	c := &serial.Config{Name: comNum, Baud: 115200}
	serialPort, err = serial.OpenPort(c)
	if err != nil {
		panic(err)
	}

	app = tview.NewApplication()

	// menu := newPrimitive("Menu")
	header := fmt.Sprintf("Open %s success\nTab for move\nEnter for select", comNum)
	grid = tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 0).
		SetBorders(true).
		AddItem(newPrimitive(header), 0, 0, 1, 3, 0, 0, false).
		AddItem(footer, 2, 0, 1, 3, 0, 0, false)

	resultText = tview.NewForm()
	for textViewIndex, textViewArrV := range appTextViewArr {
		if textViewIndex == 4 || textViewIndex == 6 {
			resultText.AddTextView(textViewArrV, "", 0, 3, true, false)
		} else {
			resultText.AddTextView(textViewArrV, "", 0, 1, true, false)
		}

	}

	menu = tview.NewList()
	// for i := 1; i <= len(appMenuArr); i++ {
	// 	title := fmt.Sprintf("List cmd %d", i)
	// 	menu.AddItem(title, appMenuArr[i-1], rune(strconv.Itoa(i)[0]), func() {
	// 		listCMD(i - 1)
	// 	})
	// }
	menu.AddItem("List cmd 0", "Get device id", '0', func() {
		listCMD(0)
	}).AddItem("List cmd 1", "Set device id", '1', func() {
		listCMD(1)
	}).AddItem("List cmd 2", "Get gpio value", '2', func() {
		listCMD(2)
	}).AddItem("List cmd 3", "Set gpio value", '3', func() {
		listCMD(3)
	}).AddItem("List cmd 4", "Get uart data", '4', func() {
		listCMD(4)
	}).AddItem("List cmd 5", "Send uart data", '5', func() {
		listCMD(5)
	}).AddItem("List cmd 6", "Get adc value", '6', func() {
		listCMD(6)
	}).AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	grid.AddItem(menu, 1, 0, 1, 1, 0, 0, true).
		AddItem(mainpage, 1, 1, 1, 1, 0, 0, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 0, false)

	go receiveData(serialPort)

	if err := app.SetRoot(grid, true).EnableMouse(false).Run(); err != nil {
		panic(err)
	}
}
