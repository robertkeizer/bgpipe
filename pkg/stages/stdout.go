package stages

import (
	"os"
	"sync"

	"github.com/bgpfix/bgpfix/msg"
	"github.com/bgpfix/bgpfix/pipe"
	"github.com/bgpfix/bgpipe/pkg/bgpipe"
)

type Stdout struct {
	*bgpipe.StageBase
	pool sync.Pool
}

func NewStdout(parent *bgpipe.StageBase) bgpipe.Stage {
	s := &Stdout{StageBase: parent}
	s.Descr = "print JSON representation to stdout"

	f := s.Flags
	f.Bool("last", true, "operate at the very end instead of here")
	// f.StringSlice("grep", []string{}, "print only given types")
	// f.StringSlice("filter", []string{}, "filter given types")
	return s
}

func (s *Stdout) Prepare() error {
	// TODO: grep /filter
	// for _, t := range s.K.Strings("grep") {
	// }

	po := &s.P.Options
	if s.K.Bool("last") {
		po.OnMsgLast(s.OnMsg, s.Dst)
	} else {
		po.OnMsg(s.OnMsg, s.Dst)
	}

	return nil
}

func (s *Stdout) OnMsg(m *msg.Msg) (action pipe.Action) {
	// TODO: yuck
	// we should collect all Callbacks in the parent, and
	// toggle an atomic on/off switch on Start
	if !s.Started() {
		return
	}

	// get from pool, marshal
	buf, _ := s.pool.Get().([]byte)
	buf = m.ToJSON(buf[:0])
	buf = append(buf, '\n')

	// write, re-use
	os.Stdout.Write(buf)
	s.pool.Put(buf)

	return
}

func (s *Stdout) Start() error {
	return nil
}
