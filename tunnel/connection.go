package tunnel

import (
	"errors"
	"net"
	"time"

	N "github.com/qqty012/clash2/common/net"
	"github.com/qqty012/clash2/common/pool"
	"github.com/qqty012/clash2/component/resolver"
	C "github.com/qqty012/clash2/constant"
)

func handleUDPToRemote(packet C.UDPPacket, pc net.PacketConn, metadata *C.Metadata) error {
	pc = unwrapPacket(pc)

	defer packet.Drop()

	// local resolve UDP dns
	if !metadata.Resolved() {
		ip, err := resolver.ResolveIP(metadata.Host)
		if err != nil {
			return err
		}
		metadata.DstIP = ip
	}

	addr := metadata.UDPAddr()
	if addr == nil {
		return errors.New("udp addr invalid")
	}

	if _, err := pc.WriteTo(packet.Data(), addr); err != nil {
		return err
	}
	// reset timeout
	pc.SetReadDeadline(time.Now().Add(udpTimeout))

	return nil
}

func handleUDPToLocal(packet C.UDPPacket, pc net.PacketConn, key string, fAddr net.Addr) {
	pc = unwrapPacket(pc)

	buf := pool.Get(pool.UDPBufferSize)
	defer pool.Put(buf)
	defer natTable.Delete(key)
	defer pc.Close()

	for {
		pc.SetReadDeadline(time.Now().Add(udpTimeout))
		n, from, err := pc.ReadFrom(buf)
		if err != nil {
			return
		}

		if fAddr != nil {
			from = fAddr
		}

		_, err = packet.WriteBack(buf[:n], from)
		if err != nil {
			return
		}
	}
}

func handleSocket(ctx C.ConnContext, outbound net.Conn) {
	left := unwrap(ctx.Conn())
	right := unwrap(outbound)

	if relayHijack(left, right) {
		return
	}

	N.Relay(left, right)
}
