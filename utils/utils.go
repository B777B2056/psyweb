package utils

type UserGender uint8

const (
	Male UserGender = iota
	Female
)

func (e UserGender) String() string {
	ret := ""
	switch e {
	case Male:
		ret = "男"
	case Female:
		ret = "女"
	}
	return ret
}

type UserStatus uint8

const (
	NotExist UserStatus = iota
	New
	WaitForReport
	Done
)

type DiseaseSeverity uint8

const (
	Asymptomatic DiseaseSeverity = iota
	SlightSymptoms
	ModerateSymptoms
	SeriousSymptoms
)

func (e DiseaseSeverity) String() string {
	ret := ""
	switch e {
	case Asymptomatic:
		ret = "无症状"
	case SlightSymptoms:
		ret = "轻微"
	case ModerateSymptoms:
		ret = "中度"
	case SeriousSymptoms:
		ret = "严重"
	}
	return ret
}
