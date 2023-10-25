package hbrelay

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/opensvc/om3/core/client"
	"github.com/opensvc/om3/core/hbtype"
	"github.com/opensvc/om3/core/omcrypto"
	"github.com/opensvc/om3/daemon/api"
	"github.com/opensvc/om3/daemon/ccfg"
	"github.com/opensvc/om3/daemon/hb/hbctrl"
	"github.com/opensvc/om3/util/plog"
)

type (
	// rx holds a hb unicast receiver
	rx struct {
		sync.WaitGroup
		ctx      context.Context
		id       string
		nodes    []string
		relay    string
		username string
		password string
		insecure bool
		timeout  time.Duration
		interval time.Duration
		lastAt   time.Time

		name   string
		log    plog.Logger
		cmdC   chan<- any
		msgC   chan<- *hbtype.Msg
		cancel func()
	}
)

// Id implements the Id function of the Receiver interface for rx
func (t *rx) Id() string {
	return t.id
}

// Stop implements the Stop function of the Receiver interface for rx
func (t *rx) Stop() error {
	t.log.Debugf("cancelling")
	t.cancel()
	for _, node := range t.nodes {
		t.cmdC <- hbctrl.CmdDelWatcher{
			HbId:     t.id,
			Nodename: node,
		}
	}
	t.Wait()
	t.log.Debugf("wait done")
	return nil
}

// Start implements the Start function of the Receiver interface for rx
func (t *rx) Start(cmdC chan<- any, msgC chan<- *hbtype.Msg) error {
	ctx, cancel := context.WithCancel(t.ctx)
	t.cmdC = cmdC
	t.msgC = msgC
	t.cancel = cancel
	ticker := time.NewTicker(t.interval)

	for _, node := range t.nodes {
		cmdC <- hbctrl.CmdAddWatcher{
			HbId:     t.id,
			Nodename: node,
			Ctx:      ctx,
			Timeout:  t.timeout,
		}
	}

	t.Add(1)
	go func() {
		defer t.Done()
		defer ticker.Stop()
		t.log.Infof("started")
		defer t.log.Infof("stopped")
		for {
			select {
			case <-ctx.Done():
				t.cancel()
				return
			case <-ticker.C:
				t.onTick()
			}
		}
	}()
	return nil
}

func (t *rx) onTick() {
	for _, node := range t.nodes {
		t.recv(node)
	}
}

func (t *rx) recv(nodename string) {
	cluster := ccfg.Get()
	cli, err := client.New(
		client.WithURL(t.relay),
		client.WithUsername(t.username),
		client.WithPassword(t.password),
		client.WithInsecureSkipVerify(t.insecure),
	)
	if err != nil {
		t.log.Errorf("recv: node %s new client: %s", nodename, err)
		return
	}

	params := api.GetRelayMessageParams{
		Nodename:  &nodename,
		ClusterId: &cluster.ID,
	}
	resp, err := cli.GetRelayMessageWithResponse(context.Background(), &params)
	if err != nil {
		t.log.Debugf("recv: node %s do request: %s", nodename, err)
		return
	} else if resp.StatusCode() != http.StatusOK {
		t.log.Debugf("unexpected get relay message %s status %s", nodename, resp.Status())
		return
	}
	if resp.JSON200 == nil {
		t.log.Debugf("recv: node %s data has no stored data", nodename)
		return
	}
	messages := resp.JSON200
	if len(messages.Messages) == 0 {
		t.log.Debugf("recv: node %s data has no stored data", nodename)
		return
	}
	c := messages.Messages[0]
	if c.UpdatedAt.IsZero() {
		t.log.Debugf("recv: node %s data has never been updated", nodename)
		return
	}
	if !t.lastAt.IsZero() && c.UpdatedAt == t.lastAt {
		t.log.Debugf("recv: node %s data has not change since last read", nodename)
		return
	}
	elapsed := time.Now().Sub(c.UpdatedAt)
	if elapsed > t.timeout {
		t.log.Debugf("recv: node %s data has not been updated for %s", nodename, elapsed)
		return
	}
	encMsg := omcrypto.NewMessage([]byte(c.Msg))
	b, msgNodename, err := encMsg.DecryptWithNode()
	if err != nil {
		t.log.Debugf("recv: decrypting node %s: %s", nodename, err)
		return
	}

	if nodename != msgNodename {
		t.log.Debugf("recv: node %s data was written by unexpected node %s: %s", nodename, msgNodename, err)
		return
	}

	msg := hbtype.Msg{}
	if err := json.Unmarshal(b, &msg); err != nil {
		t.log.Warnf("can't unmarshal msg from %s: %s", nodename, err)
		return
	}
	t.log.Debugf("recv: node %s", nodename)
	t.cmdC <- hbctrl.CmdSetPeerSuccess{
		Nodename: msg.Nodename,
		HbId:     t.id,
		Success:  true,
	}
	t.msgC <- &msg
	t.lastAt = c.UpdatedAt
}

func newRx(ctx context.Context, name string, nodes []string, relay, username, password string, insecure bool, timeout, interval time.Duration) *rx {
	id := name + ".rx"
	return &rx{
		ctx:      ctx,
		id:       id,
		nodes:    nodes,
		relay:    relay,
		username: username,
		password: password,
		insecure: insecure,
		timeout:  timeout,
		interval: interval,
		log: plog.Logger{
			Logger: plog.PkgLogger(ctx, "daemon/hb/hbrelay").With().
				Str("hb_func", "rx").
				Str("hb_name", name).
				Str("hb_id", id).
				Logger(),
			Prefix: "daemon: hb: relay: rx: " + name + ": ",
		},
	}
}
