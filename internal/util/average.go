package util

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func Average[T Number](values ...T) T {
	var total T

	for _, v := range values {
		total += v
	}

	return total / T(len(values))
}

func WeightedAverage[T Number](values, weights []T) T {
	if len(values) != len(weights) {
		panic("the number of values and averages differs")
	}
	var sum T

	for i, v := range values {
		sum += v * (weights[i] / 100)
	}

	return sum
}
