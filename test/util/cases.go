package test_util

func Case2[A, B any](aCol []A, bCol []B, callback func(a A, b B)) {
	for _, a := range aCol {
		for _, b := range bCol {
			callback(a, b)
		}
	}
}

func Case3[A, B, C any](aCol []A, bCol []B, cCol []C, callback func(a A, b B, c C)) {
	for _, item := range aCol {
		Case2(bCol, cCol, func(b B, c C) {
			callback(item, b, c)
		})
	}
}
