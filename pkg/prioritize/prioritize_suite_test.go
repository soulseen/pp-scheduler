package prioritize

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"testing"
)

func TestPrioritize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Prioritize Suite")
}

var _ = BeforeSuite(func() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	db := dir + "/test.db"
	os.Setenv("DATA_PATH", db)
})

var _ = AfterSuite(func() {
	testDB := os.Getenv("DATA_PATH")
	os.Remove(testDB)
})
