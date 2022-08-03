package v1

import (
	"bytes"
	"context"
	"fmt"
	"time"
)

/*
type PileInterface interface {
	Append(*Pile)
	Marshal() string
	Pile(*ProfileInterface)
	Clear()
}

type ProfileInterface interface {
	DeepCopyInto(*ProfileInterface)
}
*/

type SessionProfile struct {
	Alert                 string
	Cancel                context.CancelFunc
	SecurityMonitorTicker *time.Ticker
	ReqTime               time.Time
	Req                   ReqProfile
	Resp                  RespProfile
	ReqBody               BodyProfile
	RespBody              BodyProfile
	Envelop               EnvelopProfile
}

type Pile struct {
	Req      ReqPile
	Resp     RespPile
	ReqBody  BodyPile
	RespBody BodyPile
	Envelop  EnvelopPile
	Pod      PodPile
}

func (p *Pile) Clear() {
	p.Req.Clear()
	p.Resp.Clear()
	p.ReqBody.Clear()
	p.RespBody.Clear()
	p.Envelop.Clear()
	p.Pod.Clear()
}

func (p *Pile) Append(a *Pile) {
	p.Req.Append(&a.Req)
	p.Resp.Append(&a.Resp)
	p.ReqBody.Append(&a.ReqBody)
	p.RespBody.Append(&a.RespBody)
	p.Envelop.Append(&a.Envelop)
	p.Pod.Append(&a.Pod)
}

func (p *Pile) Marshal() string {
	var description bytes.Buffer
	description.WriteString("{\n")
	description.WriteString(fmt.Sprintf("  Req: %s", p.Req.Marshal(1)))
	description.WriteString(fmt.Sprintf("  Resp: %s", p.Resp.Marshal(1)))
	description.WriteString(fmt.Sprintf("  ReqBody: %s", p.ReqBody.Marshal(1)))
	description.WriteString(fmt.Sprintf("  RespBody: %s", p.RespBody.Marshal(1)))
	description.WriteString(fmt.Sprintf("  Envelop: %s", p.Envelop.Marshal(1)))
	description.WriteString(fmt.Sprintf("  Pod: %s", p.Pod.Marshal(1)))
	description.WriteString("}\n")
	return description.String()
}

func (p *Pile) Pile(sp *SessionProfile, pp *PodProfile) {
	p.Req.Add(&sp.Req)
	p.Resp.Add(&sp.Resp)
	p.ReqBody.Add(&sp.ReqBody)
	p.RespBody.Add(&sp.RespBody)
	p.Envelop.Add(&sp.Envelop)
	p.Pod.Add(pp)
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SessionProfile) DeepCopyInto(out *SessionProfile) {
	*out = *in
	out.ReqTime = in.ReqTime
	in.Req.DeepCopyInto(&out.Req)
	in.Resp.DeepCopyInto(&out.Resp)
	in.ReqBody.DeepCopyInto(&out.ReqBody)
	in.RespBody.DeepCopyInto(&out.RespBody)
	in.Envelop.DeepCopyInto(&out.Envelop)
}
