package utils

type UserGender uint8

const (
	Male UserGender = iota
	Female
)

type UserStatus uint8

const (
	NotExist UserStatus = iota
	New
	WaitForReport
	Done
)
