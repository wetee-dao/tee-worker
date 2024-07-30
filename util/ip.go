package util

import (
	"fmt"
	"math/big"

	"github.com/wetee-dao/go-sdk/pallet/types"
)

func GetUrlFromIp1(ip types.Ip1) string {
	url := ""
	if ip.Domain.IsSome {
		url = "/dns4/" + string(ip.Domain.AsSomeField0)
	} else if ip.Ipv4.IsSome {
		ipv4 := ip.Ipv4.AsSomeField0
		url = "/ip4/" + fmt.Sprintf("%d.%d.%d.%d",
			(ipv4>>24)&0xFF,
			(ipv4>>16)&0xFF,
			(ipv4>>8)&0xFF,
			ipv4&0xFF)
	} else if ip.Ipv6.IsSome {
		ipv6 := ip.Ipv6.AsSomeField0
		ipv6Int128 := big.NewInt(0)
		ipv6Int128.SetBytes(ipv6.Bytes())
		url = "/ip6/" + fmt.Sprintf("%04x:%04x:%04x:%04x:%04x:%04x:%04x:%04x",
			ipv6Int128.Rsh(ipv6Int128, 112).Uint64(),
			ipv6Int128.Rsh(ipv6Int128, 96).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 80).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 64).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 48).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 32).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 16).Uint64()&0xFFFF,
			ipv6Int128.Uint64()&0xFFFF)
	}
	return url
}

func GetUrlFromIp(ip types.Ip) string {
	d := ""
	if ip.Domain.IsSome {
		d = "/dns4/" + string(ip.Domain.AsSomeField0)
	} else if ip.Ipv4.IsSome {
		ipv4 := ip.Ipv4.AsSomeField0
		d = "/ip4/" + fmt.Sprintf("%d.%d.%d.%d",
			(ipv4>>24)&0xFF,
			(ipv4>>16)&0xFF,
			(ipv4>>8)&0xFF,
			ipv4&0xFF)
	} else if ip.Ipv6.IsSome {
		ipv6 := ip.Ipv6.AsSomeField0
		ipv6Int128 := big.NewInt(0)
		ipv6Int128.SetBytes(ipv6.Bytes())
		d = "/ip6/" + fmt.Sprintf("%04x:%04x:%04x:%04x:%04x:%04x:%04x:%04x",
			ipv6Int128.Rsh(ipv6Int128, 112).Uint64(),
			ipv6Int128.Rsh(ipv6Int128, 96).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 80).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 64).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 48).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 32).Uint64()&0xFFFF,
			ipv6Int128.Rsh(ipv6Int128, 16).Uint64()&0xFFFF,
			ipv6Int128.Uint64()&0xFFFF)
	}
	return d
}
