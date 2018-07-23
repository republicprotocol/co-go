package co_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/co-go"
)

var _ = Describe("Channels", func() {

	Context("when merging", func() {

		It("should close the output channel when the done channel is closed", func() {
			done := make(chan struct{})
			ins := make(chan (chan int))
			out := make(chan int)

			go Merge(done, ins, out)

			close(done)
			<-out
		})

		It("should merge one input", func() {
			done := make(chan struct{})
			ins := make(chan (chan int))
			out := make(chan int)

			go Merge(done, ins, out)
			go func() {
				defer close(ins)
				in := make(chan int)
				ins <- in
				for i := 0; i < 10; i++ {
					in <- i
				}
			}()
			for i := 0; i < 10; i++ {
				Expect(<-out).Should(Equal(i))
			}

			close(done)
			<-out
		})

		It("should merge multiple inputs", func() {
			done := make(chan struct{})
			ins := make(chan (chan int))
			out := make(chan int)

			go Merge(done, ins, out)
			go func() {
				defer close(ins)
				for n := 0; n < 10; n++ {
					in := make(chan int)
					ins <- in
					for i := 0; i < 10; i++ {
						in <- i
					}
				}
			}()
			js := map[int]int{}
			for i := 0; i < 10*10; i++ {
				j := <-out
				Expect(j).Should(BeNumerically(">=", 0))
				Expect(j).Should(BeNumerically("<", 10))
				js[j]++
			}
			for i := 0; i < 10; i++ {
				Expect(js[i]).Should(Equal(10))
			}

			close(done)
			<-out
		})

		Context("when using incompatible types", func() {
			It("should panic", func() {
				Expect(func() {
					Merge(make(chan struct{}), make(chan (chan float32)), make(chan int))
				}).Should(Panic())

				Expect(func() {
					Merge(make(chan struct{}), make(chan int), make(chan int))
				}).Should(Panic())

				Expect(func() {
					Merge(make(chan struct{}), make(chan float32), make(chan int))
				}).Should(Panic())

				Expect(func() {
					Merge(make(chan struct{}), make(chan (chan float32)), 0)
				}).Should(Panic())

				Expect(func() {
					Merge(make(chan struct{}), 0, make(chan int))
				}).Should(Panic())

				Expect(func() {
					Merge(make(chan struct{}), 0, 0)
				}).Should(Panic())
			})
		})
	})

	Context("when forwarding", func() {

		It("should close the output channel when the done channel is closed", func() {
			done := make(chan struct{})
			in := make(chan int)
			out := make(chan int)

			go Forward(done, in, out)

			close(done)
			<-out
		})

		It("should forward the input", func() {
			done := make(chan struct{})
			in := make(chan int)
			out := make(chan int)

			go Forward(done, in, out)
			go func() {
				defer close(in)
				for i := 0; i < 10; i++ {
					in <- i
				}
			}()
			for i := 0; i < 10; i++ {
				Expect(<-out).Should(Equal(i))
			}

			close(done)
			<-out
		})

		Context("when using incompatible types", func() {
			It("should panic", func() {
				Expect(func() {
					Forward(make(chan struct{}), make(chan (chan float32)), make(chan int))
				}).Should(Panic())

				Expect(func() {
					Forward(make(chan struct{}), make(chan float32), make(chan int))
				}).Should(Panic())

				Expect(func() {
					Forward(make(chan struct{}), make(chan (chan float32)), 0)
				}).Should(Panic())

				Expect(func() {
					Forward(make(chan struct{}), 0, make(chan int))
				}).Should(Panic())

				Expect(func() {
					Forward(make(chan struct{}), 0, 0)
				}).Should(Panic())
			})
		})
	})
})
