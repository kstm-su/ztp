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
	Leases       Leases
	Change       func(*Lease) Reply
}

var processing = map[dhcp.MessageType]map[string]struct{}{}

func init() {
	processing[dhcp.Discover] = map[string]struct{}{}
	processing[dhcp.Request] = map[string]struct{}{}
}

func (h *Handler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (replyPacket dhcp.Packet) {
	if _, ok := processing[msgType][p.CHAddr().String()]; ok {
		log.Printf("%s: ignore %s\n", p.CHAddr().String(), msgType.String())
		return
	}
	processing[msgType][p.CHAddr().String()] = struct{}{}
	switch msgType {
	case dhcp.Discover:
		lease := h.Leases.Get(p.CHAddr())
		if lease == nil {
			break
		}
		if h.Change != nil {
			if reply := h.Change(lease); reply != nil {
				replyPacket = reply.Packet(p, msgType, h, options[dhcp.OptionParameterRequestList])
				break
			}
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
		if h.Change != nil {
			if reply := h.Change(lease); reply != nil {
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
	delete(processing[msgType], p.CHAddr().String())
	replyPacket.SetSIAddr(h.ServerIPAddr)
	return
}
