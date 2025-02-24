package util

// Seq is a generator function type that takes a yield function as a parameter.
// The yield function receives values of type V and returns a boolean indicating
// whether to continue generating values.
type (
	Seq[V any]     func(yield func(V) bool)
	Seq2[K, V any] func(yield func(K, V) bool)
)

// ROverSeq is an adaptation of the range-over-func feature from Go 1.23
// for earlier versions. It converts a Seq[T] generator into a Go channel.
// The generator runs in a separate goroutine and sends the generated values to the channel.
func ROverSeq[T any](seq Seq[T]) chan T {
	ch := make(chan T)

	go func() {
		defer close(ch) // Ensure the channel is closed when the generator completes.

		seq(func(v T) bool {
			ch <- v     // Send the value to the channel.
			return true // Continue the sequence.
		})
	}()

	return ch
}

// ROverSeq2 is an adaptation of the range-over-func feature from Go 1.23
// for earlier versions. It converts a Seq2[K, V] generator into a Go channel.
// The generator runs in a separate goroutine and sends only the values (not keys) to the channel.
func ROverSeq2[K, V any](seq2 Seq2[K, V]) chan V {
	ch := make(chan V)

	go func() {
		defer close(ch) // Ensure the channel is closed when the generator completes.

		seq2(func(_ K, v V) bool {
			ch <- v     // Send only the value to the channel.
			return true // Continue the sequence.
		})
	}()

	return ch
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// ROverNumber generates a sequence of numbers from 0 to x-1 and sends them to a channel.
// It runs in a separate goroutine and closes the channel when done.
func ROverNumber[T Number](x T) chan T {
	ch := make(chan T)

	go func() {
		defer close(ch) // Ensure the channel is closed when the loop completes.
		for i := T(0); i < x; i++ {
			ch <- i // Send the current number to the channel.
		}
	}()

	return ch
}
