package validation

import (
	"bytes"
	"encoding/json"

	log "k8s-webhook/pkg/logger"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
)

type ImageValidator struct{}

func sendSlackNotification(image string, namespace string, applicationName string) error {

	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	message := map[string]interface{}{
		"text": "Validation failed for Image: *" + image + "* in Namespace: *" + namespace +
			"* for Application: *" + applicationName + "*",
	}

	jsonValue, _ := json.Marshal(message)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Logger.Debug("Error sending message to Slack:", zap.Error(err))
	}
	defer resp.Body.Close()
	return nil
}

func (n ImageValidator) Validate(deployment *appsv1.Deployment) (validation, error) {
	log.Logger.Info("inside validation")
	for _, container := range deployment.Spec.Template.Spec.Containers {
		if !strings.HasPrefix(container.Image, "{Image Prefix}") {
			sendSlackNotification(container.Image, deployment.Namespace, deployment.Name)
			return validation{Valid: true, Reason: "valid name"}, nil
		}
	}
	return validation{Valid: true, Reason: "valid name"}, nil
}
