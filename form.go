package main

type FormData struct {
	name         string
	defaultValue string
}

var (
	Form1 = []FormData{
		{"Header", "c0"},
		{"Len", "00 0a"},
		{"Type", "01"},
		{"ID Len", "01"},
		{"ID", "00"},
		{"Cmd", "ff 01 00 00"},
		{"Data", "01"},
		{"CRC", ""},
		{"Tail", "c1"},
		{"CMD(Hex)", ""},
	}
)
