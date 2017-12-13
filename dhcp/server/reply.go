package server

import (
	dhcp "github.com/krolaw/dhcp4"
)

type Reply interface {
	Packet(dhcp.Packet, dhcp.MessageType, *Handler, []byte) dhcp.Packet
	MergeOptions(dhcp.Options) dhcp.Options
}

type ACKReply struct {
	Lease   *Lease
	Options dhcp.Options
}

type NAKReply struct {
}

func (r *ACKReply) Packet(p dhcp.Packet, reqType dhcp.MessageType, h *Handler, req []byte) dhcp.Packet {
	var msgType dhcp.MessageType
	if reqType == dhcp.Discover {
		msgType = dhcp.Offer
	} else {
		msgType = dhcp.ACK
	}
	return dhcp.ReplyPacket(
		p,
		msgType,
		h.ServerIPAddr,
		r.Lease.IPAddr,
		h.Leases.Duration,
		r.MergeOptions(h.Options).SelectOrderOrAll(req),
	)
}

func (r *ACKReply) MergeOptions(options dhcp.Options) dhcp.Options {
	res := make(dhcp.Options)
	for k, v := range options {
		res[k] = v
	}
	for k, v := range r.Options {
		res[k] = v
	}
	return res
}

func (r *NAKReply) Packet(p dhcp.Packet, msgType dhcp.MessageType, h *Handler, _ []byte) dhcp.Packet {
	if msgType == dhcp.Discover {
		return nil
	}
	return dhcp.ReplyPacket(p, dhcp.NAK, h.ServerIPAddr, nil, 0, nil)
}

func (r *NAKReply) MergeOptions(_ dhcp.Options) dhcp.Options {
	return nil
}
