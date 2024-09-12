package inbound

import (
	"net"

	C "github.com/qqty012/clash2/constant"
	"github.com/qqty012/clash2/context"
	"github.com/qqty012/clash2/transport/socks5"
)

// NewHTTP receive normal http request and return HTTPContext
func NewHTTP(target socks5.Addr, rawSrc, rawDst net.Addr, conn net.Conn) *context.ConnContext {
	metadata := parseSocksAddr(target)
	metadata.NetWork = C.TCP
	metadata.Type = C.HTTP
	if ip, port, err := parseAddr(rawSrc.String()); err == nil {
		metadata.SrcIP = ip
		metadata.SrcPort = port
	}

	metadata.RawSrcAddr = rawSrc
	metadata.RawDstAddr = rawDst

	return context.NewConnContext(conn, metadata)
}
