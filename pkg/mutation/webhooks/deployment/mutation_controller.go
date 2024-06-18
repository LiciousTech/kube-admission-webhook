package deployment

import (
	"encoding/json"

	log "k8s-webhook/pkg/logger"

	"github.com/wI2L/jsondiff"
	"go.uber.org/zap"

	appsv1 "k8s.io/api/apps/v1"
)

type deploymentMutator interface {
	Mutate(*appsv1.Deployment) (*appsv1.Deployment, error)
}

type MutationController struct {
	mutators []deploymentMutator
}

func NewMutationController() *MutationController {
	return &MutationController{}
}

func (dm *MutationController) RegisterMutation(mutator deploymentMutator) {
	dm.mutators = append(dm.mutators, mutator)
}

func (dm *MutationController) ApplyMutation(rs *appsv1.Deployment) (*appsv1.Deployment, *appsv1.Deployment, error) {
	log.Logger.Info("Applying All Mutations on", zap.String("Deployment Name", rs.Name))
	modDeployment := rs.DeepCopy()
	for _, mutator := range dm.mutators {
		var err error
		modDeployment, err = mutator.Mutate(modDeployment)
		if err != nil {
			return nil, nil, err
		}
	}
	return rs, modDeployment, nil

}

func (dm *MutationController) Patch(resource interface{}, modResource interface{}) ([]byte, error) {

	patch, err := jsondiff.Compare(resource, modResource)
	if err != nil {
		return nil, err
	}

	patchb, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}

	return patchb, nil

}
