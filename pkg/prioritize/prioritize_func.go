package prioritize

import (
	log "github.com/golang/glog"
	"github.com/soulseen/pp-scheduler/pkg/sqlite"
	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	"strings"
)

const DEFAULT_PIPELINE_LABEL = "ks-pipeline"

var (
	TotalScore float64 = 10
)

func Pipeline(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
	var priorityList schedulerapi.HostPriorityList
	priorityList = make([]schedulerapi.HostPriority, len(nodes))
	keys := ParseMark(pod.Labels)
	if keys == nil {
		return zeroScore(nodes), nil
	}

	for i, node := range nodes {
		score, err := Calculation(keys, node.Name)
		if err != nil {
			panic(err)
		}

		priorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}
	return &priorityList, nil
}

func ParseMark(labels map[string]string) []string {
	label, exists := labels[DEFAULT_PIPELINE_LABEL]
	if exists == false {
		log.Info("Not exists default label: ", DEFAULT_PIPELINE_LABEL)
		return nil
	}
	keys := strings.Split(label, "-")

	return keys

}

func Calculation(keys []string, nodeName string) (int, error) {
	count := len(keys)
	scoreSegment := TotalScore / float64(count)
	var score int
	score = 1

	for i := count - 1; i >= 0; i-- {
		row, err := sqlite.KeyNodeCilent.KeyNodeSearch(keys[i], nodeName)
		if err != nil {
			return 0, err
		}
		if len(row) != 0 {
			score = int(float64(i)*scoreSegment + scoreSegment + float64(row[0].Count/10))
			log.Info("Calculation result: key-", keys[i], " score-", score)
			break
		}
	}
	return score, nil
}

func zeroScore(nodes []v1.Node) *schedulerapi.HostPriorityList {
	var priorityList schedulerapi.HostPriorityList
	priorityList = make([]schedulerapi.HostPriority, len(nodes))
	for i, node := range nodes {
		priorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: 0,
		}
	}
	return &priorityList

}
