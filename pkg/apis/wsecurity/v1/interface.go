package v1

type Pile struct {
	Req     ReqPile
	Resp    RespPile
	Process ProcessPile
}

func (p *Pile) Clear() {
	p.Req.Clear()
	p.Resp.Clear()
	p.Process.Clear()
}
