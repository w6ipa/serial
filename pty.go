package serial

/*
	Reads from the Read queue and sends to PTY port
    Reads from a pty port and write to Write Queue.
*/

import (
	"log"
	"syscall"
)

type PortPair struct {
	ptmx     int
	slave    int
	portName string
}

func NewPTY() (*PortPair, error) {

	master, err := syscall.Open("/dev/ptmx", syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_CLOEXEC, 0620)
	if err != nil {
		return nil, err
	}
	fd := uintptr(master)

	defer func() {
		if err != nil {
			_ = syscall.Close(master) // Best effort.
		}
	}()

	// Grant/unlock slave.
	if err := grantpt(fd); err != nil {
		panic(err)
	}
	if err := unlockpt(fd); err != nil {
		panic(err)
	}

	sname, err := ptsname(master)
	if err != nil {
		return nil, err
	}

	// Keep the pty open so that the other end can close/open at will without causing an EOF error
	x, err := syscall.Open(sname, syscall.O_RDWR|syscall.O_NOCTTY, 0620)
	if err != nil {
		log.Printf("[ERROR] Pty: Cannot open slave %s", err)
		return nil, err
	}

	return &PortPair{
		ptmx:     master,
		portName: sname,
		slave:    x,
	}, nil
}

func (p *PortPair) Close() {
	syscall.Close(p.slave)
	syscall.Close(p.ptmx)

}

func (p *PortPair) GetName() string {
	return p.portName
}
