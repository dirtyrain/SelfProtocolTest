package main

type FormData struct {
	name         string
	defaultValue string
}

var (
	Form0 = []FormData{
		{"Header", "c0"},
		{"Len", "00 00"},
		{"Type", "01"},
		{"ID Len", "01"},
		{"ID", "00"},
		{"Cmd", "ff 01 00 00"},
		{"Data", "01"},
		{"CRC", ""},
		{"Tail", "c1"},
		{"CMD(Hex)", ""},
	}
	Form1 = []FormData{
		{"Header", "c0"},
		{"Len", "00 00"},
		{"Type", "01"},
		{"ID Len", "01"},
		{"ID", "00"},
		{"Cmd", "ff 02 04 00"},
		{"Data", "01 02 03 04"},
		{"CRC", ""},
		{"Tail", "c1"},
		{"CMD(Hex)", ""},
	}
)
