package groups

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/raft"

	"github.com/bocchi-the-cache/indeep/api"
	"github.com/bocchi-the-cache/indeep/internal/logs"
)

type streamLayerMux struct {
	localAddr api.NodeHost
	listener  *net.TCPListener
	groups    map[api.GroupID]chan net.Conn
	groupsMu  sync.RWMutex
}

func NewStreamLayerMux(localAddr api.NodeHost) (api.StreamLayerMux, error) {
	l, err := net.Listen("tcp", string(localAddr))
	if err != nil {
		return nil, err
	}
	m := &streamLayerMux{
		localAddr: localAddr,
		listener:  l.(*net.TCPListener),
		groups:    make(map[api.GroupID]chan net.Conn),
	}
	go m.acceptMux()
	return m, nil
}

func (m *streamLayerMux) acceptMux() {
	for {
		conn, err := m.listener.Accept()
		if err != nil {
			logs.S.Error("accept error", "err", err)
			continue
		}
		groupID, err := readDialHeader(conn)
		if err != nil {
			logs.S.Error("read dial header error", "err", err)
			continue
		}
		acceptor, ok := m.getAcceptor(groupID)
		if !ok {
			logs.S.Error("unknown acceptor", "groupID", groupID)
			continue
		}
		acceptor <- conn
	}
}

func (m *streamLayerMux) getAcceptor(groupID api.GroupID) (chan net.Conn, bool) {
	m.groupsMu.RLock()
	defer m.groupsMu.RUnlock()
	ch, ok := m.groups[groupID]
	return ch, ok
}

func (m *streamLayerMux) NetworkLayer(groupID api.GroupID) raft.StreamLayer {
	return &streamLayer{
		localAddr: m.localAddr,
		groupID:   groupID,
		mux:       m,
		acceptor:  m.newAcceptor(groupID),
	}
}

func (m *streamLayerMux) newAcceptor(groupID api.GroupID) chan net.Conn {
	m.groupsMu.Lock()
	defer m.groupsMu.Unlock()
	if ch, ok := m.groups[groupID]; ok {
		close(ch)
	}
	ch := make(chan net.Conn, 1)
	m.groups[groupID] = ch
	return ch
}

func (m *streamLayerMux) Close() error {
	m.groupsMu.Lock()
	defer m.groupsMu.Unlock()
	for _, ch := range m.groups {
		close(ch)
	}
	return m.listener.Close()
}

type streamLayer struct {
	localAddr api.NodeHost
	groupID   api.GroupID
	mux       *streamLayerMux
	acceptor  <-chan net.Conn
}

func readDialHeader(r io.Reader) (api.GroupID, error) {
	lenBuf := make([]byte, 1)
	if _, err := r.Read(lenBuf); err != nil {
		return "", err
	}
	l := int(lenBuf[0])
	if l <= 0 {
		return "", fmt.Errorf("invalid group ID length %q", l)
	}
	groupID := make([]byte, l)
	if _, err := r.Read(groupID); err != nil {
		return "", err
	}
	return api.GroupID(groupID), nil
}

func writeDialHeader(w io.Writer, groupID api.GroupID) error {
	_, err := w.Write(append([]byte{byte(len(groupID))}, []byte(groupID)...))
	return err
}

func trimGroupID(addr raft.ServerAddress) string {
	s := string(addr)
	return s[:strings.Index(s, "/")]
}

func (*streamLayer) Network() string  { return "tcp" }
func (s *streamLayer) String() string { return fmt.Sprintf("%s/%s", s.localAddr, s.groupID) }

func (s *streamLayer) Accept() (net.Conn, error) {
	return <-s.acceptor, nil
}

func (*streamLayer) Close() error { return nil }

func (s *streamLayer) Addr() net.Addr { return s }

func (s *streamLayer) Dial(addr raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout(s.Network(), trimGroupID(addr), timeout)
	if err != nil {
		return nil, err
	}
	if err := writeDialHeader(conn, s.groupID); err != nil {
		return nil, err
	}
	return conn, nil
}
