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
	Form2 = []FormData{
		{"Header", "c0"},
		{"Len", "00 00"},
		{"Type", "01"},
		{"ID Len", "01"},
		{"ID", "00"},
		{"Cmd", "01 01 01 00"},
		{"Data", "00"},
		{"CRC", ""},
		{"Tail", "c1"},
		{"CMD(Hex)", ""},
	}
	Form3 = []FormData{
		{"Header", "c0"},
		{"Len", "00 00"},
		{"Type", "01"},
		{"ID Len", "01"},
		{"ID", "00"},
		{"Cmd", "02 01 01 00"},
		{"Data", "00"},
		{"CRC", ""},
		{"Tail", "c1"},
		{"CMD(Hex)", ""},
	}
	Form4 = []FormData{
		{"Header", "c0"},
		{"Len", "00 00"},
		{"Type", "01"},
		{"ID Len", "01"},
		{"ID", "00"},
		{"Cmd", "03 00 01 00"},
		{"Data", "05"},
		{"CRC", ""},
		{"Tail", "c1"},
		{"CMD(Hex)", ""},
	}
	Form5 = []FormData{
		{"Header", "c0"},
		{"Len", "00 00"},
		{"Type", "01"},
		{"ID Len", "01"},
		{"ID", "00"},
		{"Cmd", "03 00 02 00"},
		{"Data", "00"},
		{"CRC", ""},
		{"Tail", "c1"},
		{"CMD(Hex)", ""},
	}
)
