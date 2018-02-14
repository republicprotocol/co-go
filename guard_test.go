package do_test

import (
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/go-do"
)

type mockObject struct {
	GuardedObject

	notifications             []interface{}
	notificationReceived      bool
	notificationsLeftGuard    *Guard
	notificationReceivedGuard *Guard
}

func newMockObject(notifications []interface{}) *mockObject {
	obj := new(mockObject)
	obj.GuardedObject = NewGuardedObject()
	obj.notifications = notifications
	obj.notificationReceived = false
	obj.notificationsLeftGuard = obj.Guard(func() bool { return len(obj.notifications) > 0 })
	obj.notificationReceivedGuard = obj.Guard(func() bool { return obj.notificationReceived })
	return obj
}

func (obj *mockObject) Notify(notification interface{}) {
	obj.Enter(nil)
	defer obj.Exit()

	obj.notificationReceived = true
	obj.notifications = append(obj.notifications, notification)
}

func (obj *mockObject) Notification() interface{} {
	obj.Enter(obj.notificationsLeftGuard)
	defer obj.Exit()

	ret := obj.notifications[0]
	if len(obj.notifications) == 1 {
		obj.notifications = []interface{}{}
		return ret
	}
	obj.notifications = obj.notifications[1:]
	return ret
}

func (obj *mockObject) NotificationsWaiting() int {
	obj.EnterReadOnly(nil)
	defer obj.ExitReadOnly()

	return len(obj.notifications)
}

func (obj *mockObject) NotificationReceived() {
	obj.EnterReadOnly(obj.notificationReceivedGuard)
	defer obj.ExitReadOnly()
}

var _ = Describe("Mutual exclusion", func() {

	Context("when using a guarded object", func() {

		It("should express liveliness properties on initially closed guards", func() {
			for n := 0; n < 10; n++ {
				ins := int64(0)
				outs := int64(0)
				obj := newMockObject([]interface{}{})
				ret := Process(func() Option {
					CoBegin(func() Option {
						ps := make([]int, 1000)
						CoForAll(ps, func(i int) {
							obj.NotificationReceived()
						})
						return Ok(nil)
					}, func() Option {
						ps := make([]int, 1000)
						CoForAll(ps, func(i int) {
							obj.NotificationsWaiting()
						})
						return Ok(nil)
					}, func() Option {
						ps := make([]int, 1000)
						CoForAll(ps, func(i int) {
							obj.Notification()
							atomic.AddInt64(&ins, 1)
						})
						return Ok(nil)
					}, func() Option {
						ps := make([]int, 1000)
						CoForAll(ps, func(i int) {
							obj.Notify(i)
							atomic.AddInt64(&outs, 1)
						})
						return Ok(nil)
					})
					return Ok(outs - ins)
				})
				select {
				case <-time.Tick(time.Minute):
					panic("Deadlock detected")
				case val := <-ret:
					Ω(val.Ok).Should(Equal(int64(0)))
				}
			}
		})

		It("should express liveliness properties on initially open guards", func() {
			for n := 0; n < 10; n++ {
				ins := int64(0)
				outs := int64(0)
				obj := newMockObject([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
				ret := Process(func() Option {
					CoBegin(func() Option {
						ps := make([]int, 1000)
						CoForAll(ps, func(i int) {
							obj.NotificationReceived()
						})
						return Ok(nil)
					}, func() Option {
						ps := make([]int, 1000)
						CoForAll(ps, func(i int) {
							obj.NotificationsWaiting()
						})
						return Ok(nil)
					}, func() Option {
						ps := make([]int, 1000)
						CoForAll(ps, func(i int) {
							obj.Notification()
							atomic.AddInt64(&ins, 1)
						})
						return Ok(nil)
					}, func() Option {
						ps := make([]int, 1000)
						CoForAll(ps, func(i int) {
							obj.Notify(i)
							atomic.AddInt64(&outs, 1)
						})
						return Ok(nil)
					})
					return Ok(outs - ins)
				})
				select {
				case <-time.Tick(time.Minute):
					panic("Deadlock detected")
				case val := <-ret:
					Ω(val.Ok).Should(Equal(int64(0)))
				}
			}
		})

	})
})
