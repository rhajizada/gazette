package tasks

import "github.com/pgvector/pgvector-go"

func vectorFromFloat64s(s []float64) pgvector.Vector {
	temp := make([]float32, len(s))
	for i, f := range s {
		temp[i] = float32(f)
	}
	v := pgvector.NewVector(temp)
	return v
}
