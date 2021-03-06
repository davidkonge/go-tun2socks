package core

/*
#cgo CFLAGS: -I./src/include
#include "lwip/tcp.h"
#include "lwip/udp.h"
*/
import "C"
import (
	"log"
	"sync"
	"unsafe"
)

type LWIPStack interface {
	Write([]byte) (int, error)
}

// lwIP runs in a single thread, locking is needed in Go runtime.
var lwipMutex = &sync.Mutex{}

type lwipStack struct {
	tpcb *C.struct_tcp_pcb
	upcb *C.struct_udp_pcb
}

func NewLWIPStack() LWIPStack {
	tcpPCB := C.tcp_new()
	if tcpPCB == nil {
		panic("tcp_new return nil")
	}

	err := C.tcp_bind(tcpPCB, C.IP_ADDR_ANY, 0)
	switch err {
	case C.ERR_OK:
		break
	case C.ERR_VAL:
		log.Fatal("invalid PCB state")
	case C.ERR_USE:
		log.Fatal("port in use")
	default:
		C.memp_free(C.MEMP_TCP_PCB, unsafe.Pointer(tcpPCB))
		log.Fatal("unknown tcp_bind return value")
	}

	tcpPCB = C.tcp_listen_with_backlog(tcpPCB, C.TCP_DEFAULT_LISTEN_BACKLOG)

	// We can't call C function with Go functions as arguments here, it will
	// fail in compile time:
	// cannot use TCPAcceptFn (type func(unsafe.Pointer, *_Ctype_struct_tcp_pcb, _Ctype_schar) _Ctype_schar) as type *[0]byte in argument to func literal
	// I can't find other workarounds.
	// C.tcp_accept(tcpPCB, TCPAcceptFn)
	SetTCPAcceptCallback(tcpPCB)

	udpPCB := C.udp_new()
	if udpPCB == nil {
		panic("could not allocate udp pcb")
	}

	err = C.udp_bind(udpPCB, C.IP_ADDR_ANY, 0)
	if err != C.ERR_OK {
		log.Fatal("address already in use")
	}

	SetUDPRecvCallback(udpPCB, nil)

	return &lwipStack{
		tpcb: tcpPCB,
		upcb: udpPCB,
	}
}

func (s *lwipStack) Write(data []byte) (int, error) {
	return Input(data)
}

func init() {
	lwipInit()
}
