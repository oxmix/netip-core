//go:build !arm64

package utils

type TimeVal struct {
	Sec  int32
	USec int32
}

type Utmp struct {
	Type     int16
	_        [2]byte // _ alignment
	Pid      int32
	Device   [LineSize]byte
	Id       [4]byte
	User     [NameSize]byte
	Host     [HostSize]byte
	Exit     ExitStatus
	Session  int32
	Time     TimeVal
	AddrV6   [16]byte
	Reserved [20]byte // Reserved member
}
