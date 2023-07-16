package dialictl

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func ArgoCD() {
	//exec.Command("minikube", "update-context")
	_, err := exec.Command("kubectl", "create", "namespace", "argocd").Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("argocd namespace created")
	//exec.Command("minikube", "update-context")
	_, err2 := exec.Command("kubectl", "apply", "-n", "argocd", "-f", "https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml").Output()
	if err2 != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("argo is installed")
	// Wait until pods in the argocd namespace are ready
	for {
		if PodsReady("argocd") {
			break
		}
		fmt.Println("Waiting for argoCD-pods to be ready...")
		time.Sleep(30 * time.Second) // Sleep for seconds before checking again
	}
	fmt.Println("All pods in the argocd namespace are ready")
}
func CrossPlane() {
	//exec.Command("minikube", "update-context")
	_, err := exec.Command("helm", "repo", "update").Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, err2 := exec.Command("kubectl", "create", "namespace", "crossplane-system").Output()
	if err2 != nil {
		fmt.Println("Error:", err2)
		return
	}
	fmt.Println("crossplane-system namespace created")
	_, err3 := exec.Command("helm", "upgrade", "--install", "crossplane", "--namespace", "crossplane-system", "crossplane-stable/crossplane").Output()
	if err3 != nil {
		fmt.Println("Error:", err3)
		return
	}
	fmt.Println("crossplane is installed")
	time.Sleep(30 * time.Second)
	for {
		if PodsReady("crossplane-system") {
			break
		}
		fmt.Println("Waiting for crossplane-pods to be ready...")
		time.Sleep(30 * time.Second) // Sleep for seconds before checking again
	}
	fmt.Println("All pods in crossplane-system namespace are ready")
}

func Setup_EKS(provider string) {
	aroConfig := fmt.Sprintf(`apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
    name: my-cluster
    namespace: argocd
spec:
    project: default
    source:
        repoURL: https://github.com/younesELouafi/providers.git
        targetRevision: HEAD
        path: providers-%s
    destination:
        server: https://kubernetes.default.svc
        namespace: cluster-%s
    syncPolicy:
        automated:
            prune: true
            selfHeal: true
        syncOptions:
        - CreateNamespace=true`, provider, provider)

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(aroConfig)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//fmt.Println(string(output))
}
func PodsReady(namespace string) bool {
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	type PodStatus struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"conditions"`
		} `json:"status"`
	}

	type PodList struct {
		Items []PodStatus `json:"items"`
	}

	var pods PodList
	err = json.Unmarshal(output, &pods)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	for _, pod := range pods.Items {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == "Ready" && condition.Status != "True" {
				return false
			}
		}
	}

	return true
}

/*func ProvidersReady() bool {
	cmd := exec.Command("kubectl", "get", "providers", "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error here:", err)
		return false
	}

	type ProviderStatus struct {
		Status struct {
			Conditions []struct {
				Type   string `json:"type"`
				Status string `json:"status"`
			} `json:"conditions"`
		} `json:"status"`
	}

	type ProviderList struct {
		Items []ProviderStatus `json:"items"`
	}

	var providers ProviderList
	err = json.Unmarshal(output, &providers)
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	for _, provider := range providers.Items {
		for _, condition := range provider.Status.Conditions {
			if condition.Type == "Healthy" && condition.Status != "True" {
				return false
			}
		}
	}

	return true
}*/
