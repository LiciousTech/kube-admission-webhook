package deployment

import (
	"k8s.io/apimachinery/pkg/api/resource"

	log "k8s-webhook/pkg/logger"

	"go.uber.org/zap"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type ResourceMutator struct{}

func (dp *ResourceMutator) Mutate(rs *appsv1.Deployment) (*appsv1.Deployment, error) {
	log.Logger.Info("Applying Resource Mutation", zap.String("Deployment Name", rs.Name))
	rs.Spec.Template.Spec.Containers[0].Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"cpu":    resource.MustParse("100m"),
			"memory": resource.MustParse("512Mi"),
		},
		Limits: corev1.ResourceList{
			"cpu": resource.MustParse("1000m"),
		},
	}

	return rs, nil

}
