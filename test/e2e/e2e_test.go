package e2e_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"time"
)

const (
	timeout = time.Second * 25
)

var _ = Describe("", func() {

	It("In normal, should scheduler to same node.", func() {
		//create pod
		podKey1, err := createPod("/test/config/same_node/a-b-c-d.yaml")
		fmt.Println("waiting for pod1 up....")
		time.Sleep(time.Duration(20) * time.Second)
		pod1 := checkStatus(podKey1)

		podKey2, err := createPod("/test/config/same_node/a-c-e.yaml")
		fmt.Println("waiting for pod2 up....")
		time.Sleep(time.Duration(20) * time.Second)
		pod2 := checkStatus(podKey2)

		podKey3, err := createPod("/test/config/same_node/a-d.yaml")
		fmt.Println("waiting for pod3 up....")
		time.Sleep(time.Duration(20) * time.Second)
		pod3 := checkStatus(podKey3)

		podKey4, err := createPod("/test/config/same_node/d-x-y.yaml")
		fmt.Println("waiting for pod4 up....")
		time.Sleep(time.Duration(20) * time.Second)
		pod4 := checkStatus(podKey4)

		Expect(err).NotTo(HaveOccurred(), "Cannot create pod with yamls")

		pod1Node := pod1.Spec.NodeName
		pod2Node := pod2.Spec.NodeName
		pod3Node := pod3.Spec.NodeName
		pod4Node := pod4.Spec.NodeName

		fmt.Println(pod1Node, pod2Node, pod3Node, pod4Node)

		res := checkNode(pod1Node, pod2Node, pod3Node, pod4Node)

		Expect(res).To(Equal(true))

		Expect(testClient.Delete(context.TODO(), pod1)).NotTo(HaveOccurred())
		Expect(testClient.Delete(context.TODO(), pod2)).NotTo(HaveOccurred())
		Expect(testClient.Delete(context.TODO(), pod3)).NotTo(HaveOccurred())
		Expect(testClient.Delete(context.TODO(), pod4)).NotTo(HaveOccurred())

	})
})

func checkNode(args ...string) (bool) {
	var tmp string
	tmp = args[0]
	for _,key := range args {
		if tmp == key {
			tmp = key
			continue
		} else {
			return false
		}
	}
	return true
}

func createPod(yamlPath string) (types.NamespacedName, error) {
	pod_1 := &corev1.Pod{}
	reader, err := os.Open(workspace + yamlPath)
	Expect(err).NotTo(HaveOccurred(), "Cannot read sample yaml")
	err = yaml.NewYAMLOrJSONDecoder(reader, 10).Decode(pod_1)
	Expect(err).NotTo(HaveOccurred(), "Cannot unmarshal yaml")
	err = testClient.Create(context.TODO(), pod_1)
	Expect(err).NotTo(HaveOccurred())

	var podKey = types.NamespacedName{
		Name:      pod_1.Name,
		Namespace: pod_1.Namespace,
	}
	instance := &corev1.Pod{}
	err = testClient.Get(context.TODO(), podKey, instance)

	return podKey, err
}

func checkStatus(podKey types.NamespacedName) (*corev1.Pod) {
	pod := &corev1.Pod{}
	Eventually(func() error {
		err := testClient.Get(context.TODO(), podKey, pod)
		if err != nil {
			return err
		}
		if pod.Status.Phase == corev1.PodRunning {
			return nil
		}
		return fmt.Errorf("Failed")
	}, time.Minute*5, time.Second*10).Should(Succeed())

	return pod
}
