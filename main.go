package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("App started")

	fd, err := unix.Socket(unix.AF_CAN, unix.SOCK_DGRAM, unix.CAN_J1939)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ifindex, err := getIfIndex(fd, "can0")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	destAddr := &unix.SockaddrCANJ1939{
		Ifindex: ifindex,
		PGN:     0x04600,
		Addr:    0x1A,
	}

	sourceAddr := &unix.SockaddrCANJ1939{
		Ifindex: ifindex,
		PGN:     0x04600,
		Addr:    0xFA,
	}

	err = unix.Bind(fd, sourceAddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	data := []byte{0x01, 0x02}

	err = unix.Sendto(fd, data, 0, destAddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("J1939 packet sended")
}

func getIfIndex(fd int, ifName string) (int, error) {
	ifNameRaw, err := unix.ByteSliceFromString(ifName)
	if err != nil {
		return 0, err
	}
	if len(ifNameRaw) > unix.IFNAMSIZ {
		return 0, fmt.Errorf("maximum ifname length is %d characters", unix.IFNAMSIZ)
	}

	type ifreq struct {
		Name  [unix.IFNAMSIZ]byte
		Index int
	}
	var ifReq ifreq
	copy(ifReq.Name[:], ifNameRaw)
	_, _, errno := unix.Syscall(unix.SYS_IOCTL,
		uintptr(fd),
		unix.SIOCGIFINDEX,
		uintptr(unsafe.Pointer(&ifReq)))
	if errno != 0 {
		return 0, fmt.Errorf("ioctl: %v", errno)
	}

	return ifReq.Index, nil
}
