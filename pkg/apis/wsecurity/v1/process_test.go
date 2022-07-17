package v1

import (
	"testing"
	"time"
)

func TestProcessPile_Clear(t *testing.T) {
	t.Run("Process", func(t *testing.T) {
		pp := new(ProcessProfile)
		pc := new(ProcessConfig)
		pPile := new(ProcessPile)
		pp.Profile(time.Now(), time.Now(), time.Now())
		pp.Marshal(0)
		pPile.Clear()
		pPile.Add(pp)
		pPile.Marshal(0)
		pPile.Append(pPile)

		pc.AddTypicalVal()
		pc.Decide(pp)
		pc.Marshal(0)
		pc.Normalize()
		pc.Reconcile()
		pc.Merge(pc)
		pc.Learn(pPile)
		pc.Decide(pp)
		pPile.Clear()

	})

}
