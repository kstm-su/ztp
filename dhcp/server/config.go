package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	dhcp "github.com/krolaw/dhcp4"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Interface      *string
	Server_IP_Addr string `required:"true"`
	Options        OptionsDecoder
	Start_IP_Addr  string        `required:"true"`
	Lease_Range    int           `required:"true"`
	Lease_Duration time.Duration `default:"1h"`
}

type OptionsDecoder dhcp.Options

type IPDecoder net.IP

type CIDRDecoder []byte

type DomainNameDecoder []byte

func NewConfig() (*Config, error) {
	c := Config{
		Options: make(OptionsDecoder),
	}
	if err := envconfig.Process("dhcp", &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) Server(handler func(*Lease) Reply) *Server {
	return &Server{
		Handler: &Handler{
			ServerIPAddr: net.ParseIP(c.Server_IP_Addr).To4(),
			Options:      dhcp.Options(c.Options),
			Leases: Leases{
				StartIPAddr: net.ParseIP(c.Start_IP_Addr).To4(),
				Range:       c.Lease_Range,
				Duration:    c.Lease_Duration,
			},
			handlerFunc: handler,
		},
		Interface: c.Interface,
	}
}

func (o *OptionsDecoder) Decode(s string) error {
	labels := make(map[string]int)
	for i := 1; i < 0xff; i++ {
		key := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(dhcp.OptionCode(i).String(), "Option"), "Code("), ")")
		labels[key] = i
	}
	obj := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(s), &obj); err != nil {
		return err
	}
	options := make(dhcp.Options)
	for key, val := range obj {
		code, err := strconv.Atoi(key)
		if err != nil {
			var ok bool
			code, ok = labels[key]
			if !ok {
				continue
			}
		}
		var buf []byte
		switch code {
		case 1, 16, 28, 32, 50, 54, 78, 95: // net.IP
			s, ok := val.(string)
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = net.ParseIP(s).To4()
		case 3, 4, 5, 6, 7, 8, 9, 10, 11, 41, 42, 44, 45, 48, 49, 65, 68, 69, 70, 71, 72, 73, 74, 75, 76, 92, 112, 118, 138, 150: // []net.IP
			list, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = []byte{}
			for _, v := range list {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
				}
				buf = append(buf, net.ParseIP(strings.TrimSpace(s)).To4()...)
			}
		case 21, 33: // [][2]net.IP
			list, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = []byte{}
			for _, v := range list {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
				}
				pair := strings.SplitN(strings.TrimSpace(s), " ", 2)
				addr := net.ParseIP(pair[0]).To4()
				mask := net.ParseIP(pair[1]).To4()
				buf = append(buf, addr...)
				buf = append(buf, mask...)
			}
		case 121: // CIDR
			list, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = []byte{}
			for _, v := range list {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
				}
				pair := strings.SplitN(strings.TrimSpace(s), " ", 2)
				dst := CIDRDecoder{}
				if err := dst.Decode(pair[0]); err != nil {
					return err
				}
				gw := net.ParseIP(pair[1]).To4()
				buf = append(buf, []byte(dst)...)
				buf = append(buf, gw...)
			}
		case 119: // Domain Name
			list, ok := val.([]interface{})
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = []byte{}
			for _, v := range list {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
				}
				decoder := DomainNameDecoder{}
				if err := decoder.Decode(s); err != nil {
					return err
				}
				buf = append(buf, []byte(decoder)...)
			}
		case 2: // time.Duration (int32)
			s, ok := val.(string)
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			duration, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			i := int32(duration.Seconds())
			buf = []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
		case 24, 35, 38, 51, 58, 59, 91: // time.Duration (uint32)
			s, ok := val.(string)
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			duration, err := time.ParseDuration(s)
			if err != nil {
				return err
			}
			i := uint32(duration.Seconds())
			buf = []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
		case 19, 20, 27, 29, 30, 31, 34, 36, 39: // bool
			s, ok := val.(string)
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			if s == "1" || s == "yes" || s == "true" {
				buf = []byte{1}
			} else {
				buf = []byte{0}
			}
		case 23, 37, 46, 52, 53, 116: // byte
			n, ok := val.(byte)
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = []byte{byte(n)}
		case 13, 22, 26, 57, 93: // uint16
			n, ok := val.(uint16)
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = []byte{byte(n >> 8), byte(n)}
		case 25: // []uint16
			list, ok := val.([]uint16)
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = []byte{}
			for _, n := range list {
				buf = append(buf, []byte{byte(n >> 8), byte(n)}...)
			}
		default: // string
			s, ok := val.(string)
			if !ok {
				return fmt.Errorf("invalid format %s(%T)", dhcp.OptionCode(code).String(), val)
			}
			buf = []byte(s)
		}
		options[dhcp.OptionCode(code)] = buf
	}
	*o = OptionsDecoder(options)
	return nil
}

func (ip *IPDecoder) Decode(s string) error {
	*ip = IPDecoder(net.ParseIP(s))
	return nil
}

func (cidr *CIDRDecoder) Decode(s string) error {
	_, ip, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}
	mask, _ := ip.Mask.Size()
	size := (mask + 7) >> 3
	buf := make([]byte, size+1)
	buf[0] = byte(mask)
	for i := 0; i < size; i++ {
		buf[i+1] = []byte(ip.IP)[i]
	}
	*cidr = CIDRDecoder(buf)
	return nil
}

func (d *DomainNameDecoder) Decode(s string) error {
	buf := []byte{}
	for _, sub := range strings.Split(s, ".") {
		if len(sub) <= 0 {
			return fmt.Errorf(`"%s" is invalid domain name`, s)
		}
		buf = append(buf, byte(len(sub)))
		buf = append(buf, []byte(sub)...)
	}
	buf = append(buf, 0)
	*d = DomainNameDecoder(buf)
	return nil
}
