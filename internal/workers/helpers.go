package workers

import "github.com/pgvector/pgvector-go"

func float64sToVector(s []float64) pgvector.Vector {
	temp := make([]float32, len(s))
	for i, f := range s {
		temp[i] = float32(f)
	}
	v := pgvector.NewVector(temp)
	return v
}

func vectorToFloat64s(s pgvector.Vector) []float64 {
	temp := make([]float64, len(s.Slice()))
	for i, f := range s.Slice() {
		temp[i] = float64(f)
	}
	return temp
}
