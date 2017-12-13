package server

import (
	"bytes"
	"errors"
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
)

type Leases struct {
	StartIPAddr net.IP
	Range       int
	Duration    time.Duration
	Table       []*Lease
}

type Lease struct {
	CHAddr net.HardwareAddr
	IPAddr net.IP
	Expiry time.Time
	leases *Leases
}

func (l *Leases) Initialize() {
	l.Table = make([]*Lease, l.Range)
}

func (l *Leases) Get(addr net.HardwareAddr) *Lease {
	if l.Table == nil {
		l.Initialize()
	}
	for i, lease := range l.Table {
		if lease != nil && lease.Expiry.After(time.Now()) {
			continue
		}
		lease = &Lease{
			CHAddr: addr,
			IPAddr: dhcp.IPAdd(l.StartIPAddr, i),
			Expiry: time.Now().Add(l.Duration),
			leases: l,
		}
		l.Table[i] = lease
		return lease
	}
	return nil
}

func (l *Leases) Delete(addr net.HardwareAddr) {
	for _, lease := range l.Table {
		if bytes.Compare(lease.CHAddr, addr) == 0 {
			lease = nil
			break
		}
	}
}

func (l *Lease) Find() error {
	tmp := l.leases.Get(l.CHAddr)
	if tmp == nil {
		return errors.New("No IP pool space available")
	}
	*l = *tmp
	return nil
}

func (l *Lease) Update() {
	if l.leases.Table == nil {
		l.leases.Initialize()
	}
	l.Expiry = time.Now().Add(l.leases.Duration)
	i := dhcp.IPRange(l.leases.StartIPAddr, l.IPAddr) - 1
	l.leases.Table[i] = l
}
