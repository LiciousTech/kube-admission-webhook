package admission

import (
	"encoding/json"
	"fmt"
	"k8s-webhook/pkg/mutation/webhooks/deployment"
	"net/http"

	"go.uber.org/zap"

	log "k8s-webhook/pkg/logger"
	validation "k8s-webhook/pkg/validation/webhooks/deployment"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

func MutateDeploymentReview(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, error) {
	log.Logger.Info("", zap.String("Received Mutate Deployement Request for :", request.Name))
	dp, err := Deployment(request)
	if err != nil {
		e := fmt.Sprintf("could not parse deployment in admission review request: %v", err)
		return reviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	m := deployment.NewMutationController()
	m.RegisterMutation(&deployment.ResourceMutator{})
	dp, modDp, err := m.ApplyMutation(dp)
	if err != nil {
		e := fmt.Sprintf("could not mutate deployment: %v", err)
		return reviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	patch, err := m.Patch(dp, modDp)
	if err != nil {
		e := fmt.Sprintf("could not mutate deployment: %v", err)
		return reviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}
	return patchReviewResponse(request.UID, patch)
}

func Deployment(request *admissionv1.AdmissionRequest) (*appsv1.Deployment, error) {
	if request.Kind.Kind != "Deployment" {
		return nil, fmt.Errorf("only deployments are supported here")
	}

	d := appsv1.Deployment{}
	if err := json.Unmarshal(request.Object.Raw, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func ValidateDeploymentReview(request *admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, error) {
	log.Logger.Info("Received Validation Deployement Request", zap.String("Service Name:", request.Name))
	deployment, err := Deployment(request)
	if err != nil {
		e := fmt.Sprintf("could not parse pod in admission review request: %v", err)
		return reviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}

	v := validation.NewValidationController()
	v.RegisterValidator(&validation.ImageValidator{})
	val, err := v.ApplyValidation(deployment)
	if err != nil {
		e := fmt.Sprintf("could not validate pod: %v", err)
		return reviewResponse(request.UID, false, http.StatusBadRequest, e), err
	}
	if !val.Valid {
		return reviewResponse(request.UID, false, http.StatusForbidden, val.Reason), nil
	}

	return reviewResponse(request.UID, true, http.StatusAccepted, "valid deployment"), nil
}

func reviewResponse(uid types.UID, allowed bool, httpCode int32,
	reason string) *admissionv1.AdmissionReview {
	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     uid,
			Allowed: allowed,
			Result: &metav1.Status{
				Code:    httpCode,
				Message: reason,
			},
		},
	}
}

func patchReviewResponse(uid types.UID, patch []byte) (*admissionv1.AdmissionReview, error) {
	patchType := admissionv1.PatchTypeJSONPatch

	return &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:       uid,
			Allowed:   true,
			PatchType: &patchType,
			Patch:     patch,
		},
	}, nil
}
