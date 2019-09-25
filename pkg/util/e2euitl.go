package util

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func WaitForScheduler(c client.Client, namespace, name string, retryInterval, timeout time.Duration) error {
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		scheduler := &appsv1.Deployment{}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err = c.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, scheduler)
		if apierrors.IsNotFound(err) {
			fmt.Println("Cannot find scheduler")
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if scheduler.Status.ReadyReplicas == 1 {
			return true, nil
		}
		return false, nil
	})
	return err
}
