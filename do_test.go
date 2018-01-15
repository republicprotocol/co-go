package do_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/go-do"
)

var _ = Describe("Concurrency", func() {

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

	Context("when using a process", func() {
		It("should write the return value to a channel", func() {
			ret := <-Process(func() Return {
				return Value(1 + 2)
			})
			Ω(ret.Value).Should(Equal(3))
		})

		It("should write the error to a channel", func() {
			ret := <-Process(func() Return {
				return Err(errors.New("this is an error"))
			})
			Ω(ret.Err).Should(HaveOccurred())
		})
	})

})
