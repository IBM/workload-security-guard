package v1

import (
	"bytes"
	"fmt"
	"strings"
)

// Decission process
// If request profile allowed by ReqConfig: - Main Critiria
//        <Allow> + Log and gather statistics
// If Consult.Active and did not cross Consult.RequestsPerMinuete
//         If request profile allowed by Guard:  - Secondary Critiria
//                <Allow> + Log and gather statistics
// Log and gather statistics about request not allowed
// If ForceAllow
//          <Allow>		// used for example when ReqConfig is not ready
// <Block>
type Ctrl struct { // If guard needs to be consulted but is unavaliable => block
	Alert              bool   `json:"alert"`   // If true, use critiria to identify alerts
	Block              bool   `json:"block"`   // If true, block on alert.
	Learn              bool   `json:"learn"`   // If true, and no alert idetified, report piles
	Auto               bool   `json:"auto"`    // If true, use learned critiria rather than configured critiria
	Consult            bool   `json:"consult"` // False means never consult guard, all decissions are local
	RequestsPerMinuete uint16 `json:"rpm"`     // Maximum rpm allows for consulting guard
}

type Critiria struct {
	Active  bool          `json:"active"`  // If not active, Critiria ignored
	Req     ReqConfig     `json:"req"`     // Request critiria for blocking/allowing
	Resp    RespConfig    `json:"resp"`    // Response critiria for blocking/allowing
	Process ProcessConfig `json:"process"` // Processing critiria for blocking/allowing
}

type WsGate struct {
	Configured *Critiria   `json:"configured"`        // configrued critiria
	Learned    []*Critiria `json:"learned,omitempty"` // Learned citiria
	Control    Ctrl        `json:"control"`           // Control
}

func (g *WsGate) Reconcile() {
	if g.Configured != nil {
		g.Configured.Reconcile()
	}
	if g.Learned != nil {
		for _, learned := range g.Learned {
			learned.Reconcile()
		}
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
		for i, leanred := range g.Learned {
			description.WriteString(fmt.Sprintf("  Learned[%d]: %s", i, leanred.Marshal(depth+1)))
			description.WriteString(shift)
		}
	}
	description.WriteString(fmt.Sprintf("  Control: %s", g.Control.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}

func (c *Critiria) Learn(p *Pile) {
	c.Req.Learn(&p.Req)
	c.Resp.Learn(&p.Resp)
	c.Process.Learn(&p.Process)
	fmt.Printf("Critiria %v\n", c)
}
func (c *Critiria) Merge(mc *Critiria) {
	c.Active = c.Active || mc.Active
	c.Req.Merge(&mc.Req)
	c.Resp.Merge(&mc.Resp)
	c.Process.Merge(&mc.Process)
	fmt.Printf("Merged Critiria %v\n", c)
}

func (c *Critiria) Reconcile() {
	c.Req.Reconcile()
	c.Resp.Reconcile()
	c.Process.Reconcile()
}

func (c *Critiria) Normalize() {
	c.Req.Normalize()
	c.Resp.Normalize()
	c.Process.Normalize()
}
func (c *Critiria) Marshal(depth int) string {
	var description bytes.Buffer
	shift := strings.Repeat("  ", depth)
	description.WriteString("{\n")
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Req: %s", c.Req.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Resp: %s", c.Resp.Marshal(depth+1)))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Process: %s", c.Process.Marshal(depth+1)))
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
	description.WriteString(fmt.Sprintf("  Auto: %v", c.Auto))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  Consult: %v", c.Consult))
	description.WriteString(shift)
	description.WriteString(fmt.Sprintf("  RPM: %v", c.RequestsPerMinuete))
	description.WriteString(shift)
	description.WriteString("}\n")
	return description.String()
}
