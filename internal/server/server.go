package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Output struct {
	Parameters []map[string]string `json:"parameters"`
}
type GetParamsResponse struct {
	Output Output `json:"output"`
}

func Run() {
	token, ok := os.LookupEnv("API_TOKEN")

	if !ok || token == "" {
		log.Fatal("API_TOKEN environment variable is not set")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	mux.HandleFunc("POST /api/v1/getparams.execute", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Req", slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.String("remote_addr", r.RemoteAddr))

		authorization := r.Header.Get("Authorization")

		if authorization != "Bearer "+token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			slog.Warn("Request unauthorized")
			return
		}

		secrets, err := clientset.CoreV1().Secrets("").List(r.Context(), metav1.ListOptions{
			LabelSelector: "argocd.argoproj.io/secret-type: repository",
		})
		if err != nil {
			http.Error(w, "Error fetching secrets", http.StatusInternalServerError)
			slog.Error("Error fetching secrets", slog.String("error", err.Error()))
			return
		}

		repositories := make(map[string]string)

		for _, secret := range secrets.Items {
			if url, ok := secret.Data["url"]; ok {
				repositories[secret.Name] = string(url)
			} else {
				slog.Warn("Secret missing 'url' field", slog.String("secret", secret.Name))
			}
		}

		slog.Info("Fetched repositories", slog.Any("repositories", repositories))

		res := GetParamsResponse{
			Output: Output{
				Parameters: []map[string]string{repositories},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	})

}
