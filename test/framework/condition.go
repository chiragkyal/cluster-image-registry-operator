package framework

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	operatorapi "github.com/openshift/api/operator/v1"

	"github.com/openshift/cluster-image-registry-operator/pkg/defaults"
)

func ConditionExistsWithStatusAndReason(te TestEnv, conditionType string, conditionStatus operatorapi.ConditionStatus, conditionReason string) {
	var errs []error

	// Wait for the image registry resource to have an updated condition
	err := wait.Poll(1*time.Second, AsyncOperationTimeout, func() (stop bool, err error) {
		errs = nil
		conditionExists := false

		// Get a fresh version of the image registry resource
		cr, err := te.Client().Configs().Get(
			context.Background(), defaults.ImageRegistryResourceName, metav1.GetOptions{},
		)
		if err != nil {
			if errors.IsNotFound(err) {
				errs = append(errs, err)
				return false, nil
			}
			return true, err
		}
		for _, condition := range cr.Status.Conditions {
			if condition.Type == conditionType {
				conditionExists = true
				if condition.Status != conditionStatus {
					errs = append(errs, fmt.Errorf("condition %s status should be \"%v\" but was %v instead", conditionType, conditionStatus, condition.Status))
				}
				if len(conditionReason) != 0 && condition.Reason != conditionReason {
					errs = append(errs, fmt.Errorf("condition %s reason should have been \"%s\" but was %s instead", conditionType, conditionReason, condition.Reason))
				}
			}
		}
		if !conditionExists {
			errs = append(errs, fmt.Errorf("condition %s was not found, but should have been. %#v", conditionType, cr.Status.Conditions))
		}
		if len(errs) != 0 {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		errs = append(errs, err)
	}

	for _, err := range errs {
		te.Errorf("%#v", err)
	}
}

func PrunerConditionExistsWithStatusAndReason(te TestEnv, conditionType string, conditionStatus operatorapi.ConditionStatus, conditionReason string) {
	var errs []error

	// Wait for the image registry resource to have an updated condition
	err := wait.Poll(1*time.Second, AsyncOperationTimeout, func() (stop bool, err error) {
		errs = nil
		conditionExists := false

		// Get a fresh version of the image registry resource
		cr, err := te.Client().ImagePruners().Get(
			context.Background(), defaults.ImageRegistryImagePrunerResourceName, metav1.GetOptions{},
		)
		if err != nil {
			if errors.IsNotFound(err) {
				errs = append(errs, err)
				return false, nil
			}
			return true, err
		}
		for _, condition := range cr.Status.Conditions {
			if condition.Type == conditionType {
				conditionExists = true
				if condition.Status != conditionStatus {
					errs = append(errs, fmt.Errorf("condition %s status should be \"%v\" but was %v instead", conditionType, conditionStatus, condition.Status))
				}
				if len(conditionReason) != 0 && condition.Reason != conditionReason {
					errs = append(errs, fmt.Errorf("condition %s reason should have been \"%s\" but was %s instead", conditionType, conditionReason, condition.Reason))
				}
			}
		}
		if !conditionExists {
			errs = append(errs, fmt.Errorf("condition %s was not found, but should have been. %#v", conditionType, cr.Status.Conditions))
		}
		if len(errs) != 0 {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		errs = append(errs, err)
	}

	for _, err := range errs {
		te.Errorf("%#v", err)
	}
}
