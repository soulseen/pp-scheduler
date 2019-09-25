package predicate

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/core/v1"
)

var _ = Describe("Test prioritize", func() {

	It("Prioritize should always true", func() {
		pod := &v1.Pod{}
		node := &v1.Node{}
		res, _ := AlwaysTrue(*pod, *node)
		Expect(res).To(Equal(true))
	})
})
