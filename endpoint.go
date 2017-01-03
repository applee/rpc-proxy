package proxy

type Endpoint struct {
	Host    string
	Port    int
	RPCType string
	IDL     []byte
}

func NewEndpoint(host string, port int) {

}
