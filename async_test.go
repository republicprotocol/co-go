package async_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/go-async"
)

var _ = Describe("Async", func() {

	Context("when using a for all loop", func() {
		It("should apply the function to all items", func() {
			xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			ForAll(xs, func(i int) {
				xs[i] *= 2
			})
			for i := range xs {
				Ω(xs[i]).Should(Equal(i * 2))
			}
		})
	})

})
