//go:build arm64

package utils

/*
	On all Linux systems (both 32-bit and 64-bit) except for
	arm64/aarch64, a utmp record is 384 bytes in length and the
	Session, Tv.Sec and Tv.Usec fields are all 32 bits in length.

	On arm64/aarch64, a utmp record is 400 bytes in length and the
	Session, Tv.Sec and Tv.Usec fields are all 64 bits in length.

	There are two versions of this file, one for arm64/aarch64 and
	one for all other architectures.
*/

/*
	var info Utmp
	log.Println(unsafe.Sizeof(info))
	size 384 and arm 400
*/

type TimeVal struct {
	Sec  int64
	USec int64
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
	Session  int64
	Time     TimeVal
	AddrV6   [16]byte
	Reserved [20]byte // Reserved member
	_        [4]byte  // _ to align to the next record
}
