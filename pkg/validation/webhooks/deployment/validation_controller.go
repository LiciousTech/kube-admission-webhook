package validation

import (
	log "k8s-webhook/pkg/logger"

	appsv1 "k8s.io/api/apps/v1"
)

type deploymentValidator interface {
	Validate(*appsv1.Deployment) (validation, error)
}

type ValidationController struct {
	validators []deploymentValidator
}

type validation struct {
	Valid  bool
	Reason string
}

func NewValidationController() *ValidationController {
	return &ValidationController{}
}

func (dm *ValidationController) RegisterValidator(validators deploymentValidator) {
	dm.validators = append(dm.validators, validators)
}

func (dm *ValidationController) ApplyValidation(rs *appsv1.Deployment) (validation, error) {
	log.Logger.Info("Applying All Validation")
	for _, v := range dm.validators {
		var err error
		vp, err := v.Validate(rs)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !vp.Valid {
			return validation{Valid: false, Reason: vp.Reason}, err
		}
	}
	return validation{Valid: true, Reason: "valid pod"}, nil

}
