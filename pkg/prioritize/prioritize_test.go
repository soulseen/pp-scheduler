package prioritize

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/soulseen/pp-scheduler/pkg/sqlite"
	"reflect"
)

var _ = Describe("Test prioritize", func() {

	It("Calculation Should right", func() {
		dbCilent := sqlite.KeyNodeCilent

		type insertData struct {
			key      string
			nodeName string
			count    int
		}
		testsData := []insertData{
			{"a", "node1", 1},
			{"b", "node1", 1},
			{"c", "node2", 10},
			{"d", "node1", 20},
		}

		for _, lb := range testsData {
			dbCilent.KeyNodeInsert(lb.key, lb.nodeName, lb.count)
		}

		type parseMarkData struct {
			name          string
			keys          []string
			nodeName      string
			expectedScore int
		}

		tests := []parseMarkData{
			{name: "a-b-c-d", keys: []string{"a", "b", "c", "d"}, nodeName: "node1", expectedScore: 12},
			{name: "a-b-c", keys: []string{"a", "b", "c"}, nodeName: "node1", expectedScore: 6},
			{name: "z-x-y", keys: []string{"z", "x", "y"}, nodeName: "node1", expectedScore: 1},
			{name: "z-x-a", keys: []string{"z", "x", "a"}, nodeName: "node1", expectedScore: 10},
			{name: "z-x-c", keys: []string{"z", "x", "c"}, nodeName: "node2", expectedScore: 11},
			{name: "z-x-a", keys: []string{"z", "x", "a"}, nodeName: "node2", expectedScore: 1},
		}

		for _, match := range tests {
			score, err := Calculation(match.keys, match.nodeName)
			Expect(err).NotTo(HaveOccurred(), "Cannot inser data to sqliteDB")
			Expect(score).To(Equal(match.expectedScore))
		}
	})

	It("Parse labels mark Should right", func() {
		type parseMarkData struct {
			labels   map[string]string
			expected []string
		}

		tests := []parseMarkData{
			{labels: map[string]string{"ks-pipeline": "jenkins-java-maven-1"}, expected: []string{"jenkins", "java", "maven", "1"}},
			{labels: map[string]string{"ks-pipeline": "jenk/ins-java-maven-1/"}, expected: []string{"jenk/ins", "java", "maven", "1/"}},
			{labels: map[string]string{"ks-pipeline": "jenkins"}, expected: []string{"jenkins"}},
		}

		for _, lb := range tests {
			keys := ParseMark(lb.labels)
			res := reflect.DeepEqual(keys, lb.expected)
			Expect(res).To(Equal(true))
		}
	})
})
