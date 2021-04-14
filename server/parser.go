package server

type parseState struct {
	state   int
	as      int
	drop    int
	pa      pubArg
	argBuf  []byte
	msgBuf  []byte
	scratch [MAX_CONTROL_LINE_SIZE]byte
}

type pubArg struct {
	subject []byte
	reply   []byte
	szb     []byte
	size    int
}

const (
	OP_START = iota
	OP_C
	OP_CO
	OP_CON
	OP_CONN
	OP_CONNE
	OP_CONNEC
	OP_CONNECT
	CONNECT_ARG
	OP_P
	OP_PU
	OP_PUB
	PUB_ARG
	OP_PI
	OP_PIN
	OP_PING
	MSG_PAYLOAD
	MSG_END
	OP_S
	OP_SU
	OP_SUB
	SUB_ARG
	OP_U
	OP_UN
	OP_UNS
	OP_UNSU
	OP_UNSUB
	UNSUB_ARG
)
