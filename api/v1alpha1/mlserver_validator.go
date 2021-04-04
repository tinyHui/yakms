package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func validateServerReplicas(replicas int32, fldPath *field.Path) *field.Error {
	if replicas > 0 {
		return nil
	}

	return field.Invalid(fldPath, replicas, "value needs to > 0")
}

func (r *MLServer) validateServer() error {
	var allErrs field.ErrorList

	err := validateServerReplicas(r.Spec.Replicas, field.NewPath("spec").Child("replicas"))
	allErrs = append(allErrs, err)

	return apierrors.NewInvalid(
		GroupKind, r.Name, allErrs,
	)
}
