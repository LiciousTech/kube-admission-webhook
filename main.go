package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	adm "k8s-webhook/pkg/admission"
	log "k8s-webhook/pkg/logger"

	admissionv1 "k8s.io/api/admission/v1"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	log.Logger.Info("Health Check")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func handleMutation(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Recieved Muation Request")

	in, err := parseRequest(*r)
	if err != nil {
		log.Logger.Error("Error:", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Logger.Info("Mutating ", zap.String("Service ", in.Request.Name))
	out, err := adm.MutateDeploymentReview(in.Request)
	if err != nil {
		log.Logger.Error("Could not generate admission response:", zap.Error(err))
		e := fmt.Sprintf("Could not generate admission response: %v", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	log.Logger.Debug("Admission Response ", zap.Any("Admission Review: ", out))
	w.Header().Set("Content-Type", "application/json")
	jout, err := json.Marshal(out)
	if err != nil {
		log.Logger.Error("Could not parse admission response", zap.Error(err))
		e := fmt.Sprintf("Could not parse admission response: %v", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", jout)

}
func handleValidation(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Recieved Validation Request")

	in, err := parseRequest(*r)
	if err != nil {
		log.Logger.Error("Error:", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out, err := adm.ValidateDeploymentReview(in.Request)
	if err != nil {
		log.Logger.Error("Could not generate admission response:", zap.Error(err))
		e := fmt.Sprintf("Could not generate admission response: %v", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jout, err := json.Marshal(out)
	if err != nil {
		log.Logger.Error("Could not parse admission response", zap.Error(err))
		e := fmt.Sprintf("Could not parse admission response: %v", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", jout)
}

func parseRequest(r http.Request) (*admissionv1.AdmissionReview, error) {
	log.Logger.Info("Parsing Request")
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}

	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	if len(body) == 0 {
		return nil, fmt.Errorf("admission request body is empty")
	}

	var a admissionv1.AdmissionReview

	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	if a.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}

	return &a, nil
}

func main() {

	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/mutate-deployment", handleMutation)
	http.HandleFunc("/validate-deployment", handleValidation)

	log.Logger.Info("Server is starting...")
	log.Logger.Info("Server is listening on port 8080...")

	cert := "/etc/k8s-webhook/tls/tls.crt"
	key := "/etc/k8s-webhook/tls/tls.key"
	err := http.ListenAndServeTLS(":8080", cert, key, nil)
	if err != nil {
		log.Logger.Info("Error:", zap.Error(err))
		fmt.Println("Error:", err)
	}

}
