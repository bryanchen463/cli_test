package gwscli

import (
	"crypto/tls"
	"log"

	"github.com/lxzan/gws"
)

type wsHandler struct {
	messages []string
	curIndex int
}

func (h *wsHandler) OnOpen(socket *gws.Conn) {
	socket.WriteMessage(gws.OpcodeText, []byte(h.messages[h.curIndex]))
	h.curIndex += 1
}
func (h *wsHandler) OnClose(socket *gws.Conn, err error) {}
func (h *wsHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	if h.messages[h.curIndex] != message.Data.String() {
		return
	}
	h.curIndex++
	if h.curIndex == len(h.messages) {
		socket.WriteClose(uint16(gws.OpcodeCloseConnection), []byte("finished"))
		socket.NetConn().Close()
		return
	}
	socket.WriteMessage(gws.OpcodeText, []byte(h.messages[h.curIndex]))
}
func (h *wsHandler) OnPong(socket *gws.Conn, payload []byte) {}
func (h *wsHandler) OnPing(socket *gws.Conn, payload []byte) {}

func Start(addr string, message []string) error {
	app, _, err := gws.NewClient(&wsHandler{messages: message}, &gws.ClientOption{TlsConfig: &tls.Config{InsecureSkipVerify: true}})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	app.ReadLoop()
	return nil
}
