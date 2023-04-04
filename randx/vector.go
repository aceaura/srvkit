package randx

import "math/rand"

type Vector []int64

type VectorPlayer struct {
	Index  int
	Vector Vector
}

var _ rand.Source = (*VectorPlayer)(nil)

func NewVectorPlayer() *VectorPlayer {
	return &VectorPlayer{}
}

func (v VectorPlayer) Int63() int64 {
	return v.Vector[v.Index%len(v.Vector)]
}

func (v *VectorPlayer) Seed(seed int64) {
}

type VectorRecorder struct {
	Rand   *rand.Rand
	Vector Vector
}

var _ rand.Source = (*VectorRecorder)(nil)

func NewVectorRecorder() *VectorRecorder {
	return &VectorRecorder{}
}

func (r *VectorRecorder) Int63() int64 {
	n := r.Rand.Int63()
	r.Vector = append(r.Vector, n)
	return n
}

func (r *VectorRecorder) Seed(seed int64) {
	r.Rand = rand.New(rand.NewSource(seed))
}
