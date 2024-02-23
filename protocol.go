package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type Data struct {
	Header byte   `json:"header"`
	Length uint16 `json:"length"`
	Typ    byte   `json:"typ"`
	IDLen  byte   `json:"idLen"`
	ID     []byte `json:"id"`
	Cmd    uint32 `json:"cmd"`
	Data   []byte `json:"data"`
	Crc16  uint16 `json:"crc16"`
	Tail   byte   `json:"tail"`
}

func ProtocolAnalysis(recvData *[]byte) (Data, error) {
	var d Data
	var err error

	reader := bytes.NewReader(*recvData)
	if d.Header, err = reader.ReadByte(); err != nil {
		fmt.Println("Read header failed: ", err)
		return d, err
	}

	if err = binary.Read(reader, binary.BigEndian, &d.Length); err != nil {
		fmt.Println("Read length failed: ", err)
		return d, err
	}

	if d.Typ, err = reader.ReadByte(); err != nil {
		fmt.Println("Read type failed: ", err)
		return d, err
	}

	if d.IDLen, err = reader.ReadByte(); err != nil {
		fmt.Println("Read idLen failed: ", err)
		return d, err
	}

	d.ID = make([]byte, d.IDLen)
	if _, err := reader.Read(d.ID); err != nil {
		fmt.Println("Read id failed: ", err)
		return d, err
	}

	if err := binary.Read(reader, binary.BigEndian, &d.Cmd); err != nil {
		fmt.Println("Read command failed", err)
		return d, err
	}

	dataLen := int(d.Length) - 1 - 1 - int(d.IDLen) - 4 - 2 - 1
	d.Data = make([]byte, dataLen)
	if _, err := reader.Read(d.Data); err != nil {
		fmt.Println("Read data failed: ", err)
		return d, err
	}

	if err = binary.Read(reader, binary.BigEndian, &d.Crc16); err != nil {
		fmt.Println("Read CRC16 failed: ", err)
		return d, err
	}

	if d.Tail, err = reader.ReadByte(); err != nil {
		fmt.Println("Read tail failed: ", err)
		return d, err
	}

	return d, nil
}

func Test() {
	// 定义一个字节数组，表示一个完整的协议数据包
	testData := []byte{0xc0, 0x00, 0x29, 0x06, 0x10, 0x68, 0x54, 0x04, 0x2d, 0x35, 0x37, 0x39, 0x48, 0x38, 0x4e, 0xf9, 0x80, 0x27, 0x27, 0x54, 0x54, 0xff, 0x01, 0x00, 0x00, 0x68, 0x54, 0x04, 0x2d, 0x35, 0x37, 0x39, 0x48, 0x38, 0x4e, 0xf9, 0x80, 0x27, 0x27, 0x54, 0x54, 0x33, 0xfb, 0xc1}

	d, err := ProtocolAnalysis(&testData)
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
			printStr = fmt.Sprintf("%02x %02x", (value >> 8 & 0xFF), (value & 0xFF))
		case uint32:
			printStr = fmt.Sprintf("%02x %02x %02x %02x", (value>>24)&0xFF, (value>>16)&0xFF, (value>>8)&0xFF, value&0xFF)
		case []byte:
			printStr = ""
			for _, byteValue := range value {
				printStr += fmt.Sprintf("%02x ", byteValue)
			}
		}
		fmt.Println(i, printStr)
	}

	// 遍历字节数组，根据协议字段，分别打印出各个字段的内容
	reader := bytes.NewReader(testData)

	// 读取协议头
	var header byte
	if header, err = reader.ReadByte(); err != nil {
		fmt.Println("Read header failed:", err)
		return
	}
	if header != 0xC0 {
		fmt.Println("Invalid header: ", header)
		return
	}

	// 读取协议长度
	var length uint16
	if err := binary.Read(reader, binary.BigEndian, &length); err != nil {
		fmt.Println("Read length failed:", err)
		return
	}

	// 读取协议类型
	var typ byte
	if typ, err = reader.ReadByte(); err != nil {
		fmt.Println("Read type failed:", err)
		return
	}

	// 读取ID长度
	var idLen byte
	if idLen, err = reader.ReadByte(); err != nil {
		fmt.Println("Read ID length failed:", err)
		return
	}

	// 计算数据部分的长度
	dataLen := int(length) - 1 - 1 - int(idLen) - 4 - 2 - 1

	// 读取ID
	var id []byte
	if idLen > 0 {
		id = make([]byte, idLen)
		if _, err := reader.Read(id); err != nil {
			fmt.Println("Read ID failed:", err)
			return
		}
	}

	// 读取命令
	var cmd uint32
	if err := binary.Read(reader, binary.BigEndian, &cmd); err != nil {
		fmt.Println("Read command failed:", err)
		return
	}

	// 读取数据
	data := make([]byte, dataLen)
	if _, err := reader.Read(data); err != nil {
		fmt.Println("Read data failed:", err)
		return
	}

	// 读取CRC16
	var crc16 uint16
	if err := binary.Read(reader, binary.BigEndian, &crc16); err != nil {
		fmt.Println("Read CRC16 failed:", err)
		return
	}

	// 读取协议尾
	var tail byte
	if tail, err = reader.ReadByte(); err != nil {
		fmt.Println("Read tail failed:", err)
		return
	}
	if tail != 0xC1 {
		fmt.Println("Invalid tail: ", tail)
		return
	}

	// 打印出各个字段的内容
	fmt.Println("Header:", header)
	fmt.Println("Length:", length)
	fmt.Println("Type:", typ)
	fmt.Println("ID Length:", idLen)
	fmt.Println("ID:", id)
	fmt.Println("Command:", cmd)
	fmt.Println("Data:", data)
	fmt.Println("CRC16:", crc16)
	fmt.Println("Tail:", tail)
}
