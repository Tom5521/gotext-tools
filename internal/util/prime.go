package util

type u32 = uint32

func NextPrime(n u32) u32 {
	if n <= 2 {
		return 2
	}

	if n%2 == 0 {
		n++
	}

	for {
		if IsPrime(n) {
			return n
		}
		n += 2
	}
}

func IsPrime(n u32) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := u32(3); i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}
