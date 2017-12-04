package server

import (
	"log"
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
)

type Handler struct {
	ServerIPAddr net.IP
	Options      dhcp.Options
	Leases       *Leases
	handlerFunc  func(*Lease) Reply
}

var processing = map[string]struct{}{}

func (h *Handler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (replyPacket dhcp.Packet) {
	if _, ok := processing[p.CHAddr().String()]; ok {
		log.Printf("%s: ignore %s\n", p.CHAddr().String(), msgType.String())
		return
	}
	processing[p.CHAddr().String()] = struct{}{}
	defer delete(processing, p.CHAddr().String())
	log.Println("message type: ", msgType)
	switch msgType {
	case dhcp.Discover:
		lease := &Lease{
			CHAddr: p.CHAddr(),
			Expiry: time.Now().Add(h.Leases.Duration),
			leases: h.Leases,
		}
		if h.handlerFunc != nil {
			if reply := h.handlerFunc(lease); reply != nil {
				replyPacket = reply.Packet(p, msgType, h, options[dhcp.OptionParameterRequestList])
				break
			}
		}
		lease = h.Leases.Get(p.CHAddr())
		if lease == nil {
			break
		}
		replyPacket = dhcp.ReplyPacket(
			p,
			dhcp.Offer,
			h.ServerIPAddr,
			lease.IPAddr,
			h.Leases.Duration,
			h.Options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]),
		)
	case dhcp.Request:
		if addr, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(addr).Equal(h.ServerIPAddr) {
			break
		}
		req := net.IP(options[dhcp.OptionRequestedIPAddress]).To4()
		if req == nil {
			req = net.IP(p.CIAddr()).To4()
		}
		if len(req) != 4 || req.Equal(net.IPv4zero) {
			replyPacket = dhcp.ReplyPacket(p, dhcp.NAK, h.ServerIPAddr, nil, 0, nil)
			break
		}
		i := dhcp.IPRange(h.Leases.StartIPAddr, req) - 1
		if i < 0 || h.Leases.Range <= i {
			replyPacket = dhcp.ReplyPacket(p, dhcp.NAK, h.ServerIPAddr, nil, 0, nil)
			break
		}
		lease := h.Leases.Table[i]
		if lease == nil {
			replyPacket = dhcp.ReplyPacket(p, dhcp.NAK, h.ServerIPAddr, nil, 0, nil)
			break
		}
		lease.Expiry = time.Now().Add(h.Leases.Duration)
		if h.handlerFunc != nil {
			if reply := h.handlerFunc(lease); reply != nil {
				replyPacket = reply.Packet(p, msgType, h, options[dhcp.OptionParameterRequestList])
				break
			}
		}
		replyPacket = dhcp.ReplyPacket(
			p,
			dhcp.ACK,
			h.ServerIPAddr,
			req,
			h.Leases.Duration,
			h.Options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]),
		)
	case dhcp.Release, dhcp.Decline:
		h.Leases.Delete(p.CHAddr())
	}
	if len(replyPacket) > 0 {
		replyPacket.SetSIAddr(h.ServerIPAddr)
	}
	return
}
