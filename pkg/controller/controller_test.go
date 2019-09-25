package controller

//import (
//	"reflect"
//	"testing"
//)
//
//func TestProcessNextItem(t *testing.T) {
//
//	type parseMarkData struct {
//		labels   map[string]string
//		expected []string
//	}
//
//	tests := []parseMarkData{
//		{labels: map[string]string{"ks-pipeline": "jenkins-java-maven-1"}, expected: []string{"jenkins", "java", "maven", "1"}},
//		{labels: map[string]string{"ks-pipeline": "jenk/ins-java-maven-1/"}, expected: []string{"jenk/ins", "java", "maven", "1/"}},
//	}
//
//	for _, lb := range tests {
//		if res := processItem(lb.labels); ! reflect.DeepEqual(res, lb.expected) {
//			t.Errorf(lb.labels["ks-pipeline"])
//			t.Errorf("Expected key %s, but got %s", lb.expected, res)
//		}
//	}
//}
