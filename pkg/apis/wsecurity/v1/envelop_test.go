package v1

import (
	"testing"
	"time"
)

func TestEnvelopPile_Clear(t *testing.T) {
	t.Run("Envelop", func(t *testing.T) {
		pp := new(EnvelopProfile)
		pc := new(EnvelopConfig)
		pPile := new(EnvelopPile)
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
