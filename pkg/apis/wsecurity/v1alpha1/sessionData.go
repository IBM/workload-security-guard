package v1alpha1

import "fmt"

//////////////////// SessionDataProfile ////////////////

// Exposes ValueProfile interface
type SessionDataProfile struct {
	Req      ReqProfile     `json:"req"`
	Resp     RespProfile    `json:"resp"`
	ReqBody  BodyProfile    `json:"reqbody"`
	RespBody BodyProfile    `json:"respbody"`
	Envelop  EnvelopProfile `json:"envelop"`
	Pod      PodProfile     `json:"pod"`
}

// This is Profile is never used by Guard - for testing only
func (profile *SessionDataProfile) Profile(args ...interface{}) {
	// never used
	profile.Req.Profile(args[0], args[1])
	profile.Resp.Profile(args[2])
	profile.ReqBody.Profile(args[3])
	profile.RespBody.Profile(args[4])
	profile.Envelop.Profile(args[5], args[6], args[7])
	profile.Pod.Profile()
}

//////////////////// SessionDataPile ////////////////

// Exposes ValuePile interface
type SessionDataPile struct {
	Req      ReqPile     `json:"req"`
	Resp     RespPile    `json:"resp"`
	ReqBody  BodyPile    `json:"reqbody"`
	RespBody BodyPile    `json:"respbody"`
	Envelop  EnvelopPile `json:"envelop"`
	Pod      PodPile     `json:"pod"`
}

func (pile *SessionDataPile) Add(valProfile ValueProfile) {
	profile := valProfile.(*SessionDataProfile)

	pile.Req.Add(&profile.Req)
	pile.Resp.Add(&profile.Resp)
	pile.ReqBody.Add(&profile.ReqBody)
	pile.RespBody.Add(&profile.RespBody)
	pile.Envelop.Add(&profile.Envelop)
	pile.Pod.Add(&profile.Pod)
}

func (pile *SessionDataPile) Clear() {
	pile.Req.Clear()
	pile.Resp.Clear()
	pile.ReqBody.Clear()
	pile.RespBody.Clear()
	pile.Envelop.Clear()
	pile.Pod.Clear()
}

func (pile *SessionDataPile) Merge(otherValPile ValuePile) {
	otherPile := otherValPile.(*SessionDataPile)
	pile.Req.Merge(&otherPile.Req)
	pile.Resp.Merge(&otherPile.Resp)
	pile.ReqBody.Merge(&otherPile.ReqBody)
	pile.RespBody.Merge(&otherPile.RespBody)
	pile.Envelop.Merge(&otherPile.Envelop)
	pile.Pod.Merge(&otherPile.Pod)
}

//////////////////// SessionDataConfig ////////////////

// Exposes ValueConfig interface
type SessionDataConfig struct {
	Active   bool          `json:"active"`   // If not active, criteria ignored
	Req      ReqConfig     `json:"req"`      // Request criteria for blocking/allowing
	Resp     RespConfig    `json:"resp"`     // Response criteria for blocking/allowing
	ReqBody  BodyConfig    `json:"reqbody"`  // Request body criteria for blocking/allowing
	RespBody BodyConfig    `json:"respbody"` // Response body criteria for blocking/allowing
	Envelop  EnvelopConfig `json:"envelop"`  // Envelop criteria for blocking/allowing
	Pod      PodConfig     `json:"pod"`      // Pod criteria for blocking/allowing
}

// This is Decide is never used by Guard - for testing only
func (config *SessionDataConfig) Decide(valProfile ValueProfile) string {
	// never used
	profile := valProfile.(*SessionDataProfile)
	if ret := config.Req.Decide(&profile.Req); ret != "" {
		return fmt.Sprintf("Req: %s", ret)
	}
	if ret := config.Resp.Decide(&profile.Resp); ret != "" {
		return fmt.Sprintf("Resp: %s", ret)
	}
	if ret := config.ReqBody.Decide(&profile.ReqBody); ret != "" {
		return fmt.Sprintf("ReqBody: %s", ret)
	}
	if ret := config.RespBody.Decide(&profile.RespBody); ret != "" {
		return fmt.Sprintf("RespBody: %s", ret)
	}
	if ret := config.Envelop.Decide(&profile.Envelop); ret != "" {
		return fmt.Sprintf("Envelop: %s", ret)
	}
	if ret := config.Pod.Decide(&profile.Pod); ret != "" {
		return fmt.Sprintf("Pod: %s", ret)
	}

	return ""
}

func (config *SessionDataConfig) Learn(valPile ValuePile) {
	pile := valPile.(*SessionDataPile)

	config.Req.Learn(&pile.Req)
	config.Resp.Learn(&pile.Resp)
	config.ReqBody.Learn(&pile.ReqBody)
	config.RespBody.Learn(&pile.RespBody)
	config.Envelop.Learn(&pile.Envelop)
	config.Pod.Learn(&pile.Pod)
}

func (config *SessionDataConfig) Fuse(otherValConfig ValueConfig) {
	otherConfig := otherValConfig.(*SessionDataConfig)

	config.Active = config.Active || otherConfig.Active
	config.Req.Fuse(&otherConfig.Req)
	config.Resp.Fuse(&otherConfig.Resp)
	config.ReqBody.Fuse(&otherConfig.ReqBody)
	config.RespBody.Fuse(&otherConfig.RespBody)
	config.Envelop.Fuse(&otherConfig.Envelop)
	config.Pod.Fuse(&otherConfig.Pod)
}
