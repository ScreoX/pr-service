package providers

import (
	"math/rand"
	"time"

	"pr-service/internal/app"
)

type random struct {
	random *rand.Rand
}

func NewRealRandom() app.RandomProvider {
	return &random{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *random) Shuffle(n int, swapFunc func(i, j int)) {
	r.random.Shuffle(n, swapFunc)
}

func (r *random) Intn(n int) int {
	return r.random.Intn(n)
}
