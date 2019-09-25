package e2e_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"testing"
	"time"

	"github.com/soulseen/pp-scheduler/pkg/util"
)

var (
	testClient    client.Client
	cfg           *rest.Config
	workspace     string
	testNamespace string
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = BeforeSuite(func() {
	//install deploy
	testNamespace = os.Getenv("TEST_NS")
	Expect(testNamespace).ShouldNot(BeEmpty())
	workspace = getWorkspace() + "/../.."
	cfg, err := config.GetConfig()
	Expect(err).ShouldNot(HaveOccurred(), "Error reading kubeconfig")
	c, err := client.New(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred(), "Error in creating client")
	testClient = c
	//waiting for scheduler up
	err = util.WaitForScheduler(c, "kubesphere-system", "ks-scheduler", 5*time.Second, 3*time.Minute)
	Expect(err).ShouldNot(HaveOccurred(), "timeout waiting for scheduler up: %s\n", err)

	fmt.Fprintf(GinkgoWriter, "ks-scheduler is up now")
})

var _ = AfterSuite(func() {
	cmd := exec.Command("kubectl", "delete", "-f", workspace+"/deploy/ks-scheduler.yaml")
	Expect(cmd.Run()).ShouldNot(HaveOccurred())
})

func getWorkspace() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}
