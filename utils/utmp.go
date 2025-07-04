package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"time"
)

const (
	LineSize = 32
	NameSize = 32
	HostSize = 256
)

type ExitStatus struct {
	Termination int16
	Exit        int16
}

// Addr returns the IPv4 or IPv6 address of the login record.
func (r *Utmp) Addr() net.IP {
	ip := make(net.IP, 16)
	// no error checking: reading from r.AddrV6 cannot fail
	_ = binary.Read(bytes.NewReader(r.AddrV6[:]), binary.BigEndian, ip)
	if bytes.Equal(ip[4:], net.IPv6zero[4:]) {
		// IPv4 address, shorten the slice so that net.IP behaves correctly:
		ip = ip[:4]
	}
	return ip
}

func ReadUtmp(file string) ([]*UtmpSmart, error) {
	var us []*UtmpSmart

	f, err := os.Open(file)
	if err != nil {
		return us, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	for {
		u, readErr := readLine(f)
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			if errors.Is(readErr, io.ErrUnexpectedEOF) {
				log.Println("[warn] parse utmp:", readErr)
				continue
			}
			return nil, readErr
		}
		us = append(us, utmpConv(u))
	}

	sort.Slice(us, func(i, j int) bool {
		return us[i].Time.Before(us[j].Time)
	})

	return us, nil
}

func readLine(file io.Reader) (*Utmp, error) {
	u := new(Utmp)

	err := binary.Read(file, binary.LittleEndian, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

type UtmpSmart struct {
	Type    int
	Pid     int
	Device  string
	Id      string
	User    string
	Host    string
	Exit    ExitStatus
	Session int
	Time    time.Time
	Addr    string
}

func utmpConv(u *Utmp) *UtmpSmart {
	return &UtmpSmart{
		Type:   int(u.Type),
		Pid:    int(u.Pid),
		Device: string(u.Device[:getByteLen(u.Device[:])]),
		Id:     string(u.Id[:getByteLen(u.Id[:])]),
		User:   string(u.User[:getByteLen(u.User[:])]),
		Host:   string(u.Host[:getByteLen(u.Host[:])]),
		Exit: ExitStatus{
			Termination: u.Exit.Termination,
			Exit:        u.Exit.Exit,
		},
		Session: int(u.Session),
		Time:    time.Unix(int64(u.Time.Sec), 0).UTC(),
		Addr:    u.Addr().String(),
	}
}

// get byte \0 index
func getByteLen(byteArray []byte) int {
	n := bytes.IndexByte(byteArray[:], 0)
	if n == -1 {
		return 0
	}

	return n
}
