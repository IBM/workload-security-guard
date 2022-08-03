package v1

import (
	"bytes"
	"fmt"
	"strings"
)

// Decission process
// If request profile allowed by ReqConfig: - Main Criteria
//        <Allow> + Log and gather statistics
// If Consult.Active and did not cross Consult.RequestsPerMinuete
//         If request profile allowed by Guard:  - Secondary Criteria
//                <Allow> + Log and gather statistics
// Log and gather statistics about request not allowed
// If ForceAllow
//          <Allow>		// used for example when ReqConfig is not ready
// <Block>
type Ctrl struct { // If guard needs to be consulted but is unavaliable => block
	Alert              bool   `json:"alert"`   // If true, use criteria to identify alerts
	Block              bool   `json:"block"`   // If true, block on alert.
	Learn              bool   `json:"learn"`   // If true, and no alert idetified, report piles
	Force              bool   `json:"force"`   // If true, learning is done even when alert idetified, report piles
	Auto               bool   `json:"auto"`    // If true, use learned criteria rather than configured criteria
	Consult            bool   `json:"consult"` // False means never consult guard, all decissions are local
	RequestsPerMinuete uint16 `json:"rpm"`     // Maximum rpm allows for consulting guard
}

type Criteria struct {
	Active   bool          `json:"active"`   // If not active, criteria ignored
	Req      ReqConfig     `json:"req"`      // Request criteria for blocking/allowing
	Resp     RespConfig    `json:"resp"`     // Response criteria for blocking/allowing
	ReqBody  BodyConfig    `json:"reqbody"`  // Request body criteria for blocking/allowing
	RespBody BodyConfig    `json:"respbody"` // Response body criteria for blocking/allowing
	Envelop  EnvelopConfig `json:"envelop"`  // Envelop criteria for blocking/allowing
	Pod      PodConfig     `json:"pod"`      // Pod criteria for blocking/allowing
}

type WsGate struct {
	Configured *Criteria `json:"configured"`        // configrued criteria
	Learned    *Criteria `json:"learned,omitempty"` // Learned citeria
	Control    *Ctrl     `json:"control"`           // Control
}

func (g *WsGate) AutoActivate() {
	g.Control.Auto = true
	g.Control.Learn = true
	g.Control.Force = true
	g.Control.Alert = true
}

func (g *WsGate) Reconcile() {
	if g.Configured != nil {
		g.Configured.Reconcile()
	}
	if g.Learned != nil {
		g.Learned.Reconcile()
	}
	if g.Control == nil {
		g.Control = new(Ctrl)
	}
	g.Control.Reconcile()
}

func (g *WsGate) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	if g.Configured != nil {
		description.WriteString(fmt.Sprintf("  Configrued: %s", g.Configured.Marshal(depth+1)))
		description.WriteString(shift)
	}
	if g.Learned != nil {
		description.WriteString(fmt.Sprintf("  Learned: %s", g.Learned.Marshal(depth+1)))
		description.WriteString(shift)
	}
	description.WriteString(fmt.Sprintf("  Control: %s", g.Control.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (c *Criteria) Learn(p *Pile) {
	c.Req.Learn(&p.Req)
	c.Resp.Learn(&p.Resp)
	c.ReqBody.Learn(&p.ReqBody)
	c.RespBody.Learn(&p.RespBody)
	c.Envelop.Learn(&p.Envelop)
	c.Pod.Learn(&p.Pod)
	//fmt.Printf("Criteria %v\n", c)
}
func (c *Criteria) Merge(mc *Criteria) {
	c.Active = c.Active || mc.Active
	c.Req.Merge(&mc.Req)
	c.Resp.Merge(&mc.Resp)
	c.ReqBody.Merge(&mc.ReqBody)
	c.RespBody.Merge(&mc.RespBody)
	c.Envelop.Merge(&mc.Envelop)
	c.Pod.Merge(&mc.Pod)
	//fmt.Printf("Merged Criteria %v\n", c)
}

func (c *Criteria) Reconcile() {
	c.Req.Reconcile()
	c.Resp.Reconcile()
	c.ReqBody.Reconcile()
	c.RespBody.Reconcile()
	//c.Envelop.Reconcile()
	c.Pod.Reconcile()
}

func (c *Criteria) Normalize() {
	c.Req.Normalize()
	c.Resp.Normalize()

	c.Envelop.Normalize()
	//c.Pod.Normalize()
}
func (c *Criteria) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Active: %v\n", c.Active))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Req: %s", c.Req.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Resp: %s", c.Resp.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  ReqBody: %s", c.ReqBody.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  RespBody: %s", c.RespBody.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Envelop: %s", c.Envelop.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Pod: %s", c.Pod.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (c *Ctrl) Reconcile() {

}

func (c *Ctrl) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Alert: %v", c.Alert))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Block: %v", c.Block))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Learn: %v", c.Learn))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Force: %v", c.Force))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Auto: %v", c.Auto))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Consult: %v", c.Consult))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  RPM: %v", c.RequestsPerMinuete))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}
