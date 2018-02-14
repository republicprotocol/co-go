package do_test

import (
	"log"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"

	. "github.com/republicprotocol/go-do"
)

type mockObject struct {
	GuardedObject

	notifications          []interface{}
	notificationsLeftGuard *Guard
}

func newMockObject() *mockObject {
	obj := new(mockObject)
	obj.GuardedObject = NewGuardedObject()
	obj.notifications = []interface{}{}
	obj.notificationsLeftGuard = obj.Guard(func() bool { return len(obj.notifications) > 0 })
	return obj
}

func (obj *mockObject) Notify(notification interface{}) {
	obj.Enter(nil)
	defer obj.Exit()

	log.Println("Notify")
	obj.notifications = append(obj.notifications, notification)
}

func (obj *mockObject) Notification() interface{} {
	obj.Enter(obj.notificationsLeftGuard)
	defer obj.Exit()

	log.Println("Notification", len(obj.notifications))
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

var _ = Describe("Mutual exclusion", func() {

	Context("when using guarded objects", func() {

		It("should not produce any racing errors", func() {
			obj := newMockObject()
			for {
				log.Println("=================")
				CoBegin(func() Option {
					ps := make([]int, 3)
					CoForAll(ps, func(i int) {
						obj.Notification()
					})
					return Ok(nil)
				}, func() Option {
					ps := make([]int, 3)
					CoForAll(ps, func(i int) {
						obj.Notify(i)
					})
					return Ok(nil)
				})
			}
		})

	})

})
