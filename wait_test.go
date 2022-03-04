/*
Copyright 2021 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package helm

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
	"sigs.k8s.io/e2e-framework/third_party/helm"
)

var curDir, _ = os.Getwd()

var applyList = []string{
	"wait --for=condition=ready pod --selector=app=nginx-nginx-ingress -n ingress-nginx --timeout=160s",
	"patch svc nginx-nginx-ingress -n ingress-nginx --patch \"$(cat istio-filter-lab/patches/patch_ingresscontroller.yaml)\"",
	"apply -f istio-filter-lab/httpbin.yaml",
}

func executeCommands(kubeconfig string, cmdList []string, t *testing.T) {
	log.Printf("Using kubeconfig file %s\n", kubeconfig)
	for _, v := range cmdList {
		output, err := executeKubectl(kubeconfig, v)
		if err != nil {
			t.Fatal(fmt.Sprintf("failed to invoke %s using config %s due to an error", v, kubeconfig), err)
		}
		fmt.Printf("%s\n", output)
	}
}

func executeKubectl(configFile string, resource string) (string, error) {

	cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl --kubeconfig=%s %s", configFile, resource))
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return "", err
	}

	return out.String(), nil
}

func TestLocalHelmChartWorkflow(t *testing.T) {
	feature := features.New("Local Helm chart workflow").
		Setup(func(ctx context.Context, t *testing.T, config *envconf.Config) context.Context {
			manager := helm.New(config.KubeconfigFile())

			executeCommands(config.KubeconfigFile(), []string{
				fmt.Sprintf("label namespace %s istio-injection=enabled --overwrite", ingressNamespace),
				"label namespace default istio-injection=enabled --overwrite",
				"label node --all=true ingress-ready=true",
			}, t)
			err := manager.RunRepo(helm.WithArgs("add", "nginx-stable", "https://helm.nginx.com/stable"))
			if err != nil {
				t.Fatal("failed to add nginx helm chart repo")
			}
			err = manager.RunRepo(helm.WithArgs("update"))
			if err != nil {
				t.Fatal("failed to upgrade helm repo")
			}
			err = manager.RunInstall(helm.WithName("nginx"), helm.WithNamespace(ingressNamespace), helm.WithReleaseName("nginx-stable/nginx-ingress"))
			if err != nil {
				t.Fatal("failed to install nginx Helm chart")
			}
			err = manager.RunInstall(helm.WithName("istio-base"), helm.WithNamespace(istioNamespace), helm.WithChart(filepath.Join(istioHome, "manifests", "charts", "base")), helm.WithWait(), helm.WithTimeout("10m"))
			if err != nil {
				t.Fatal("failed to invoke helm install operation istio-base due to an error", err)
			}
			err = manager.RunInstall(helm.WithName("istiod"), helm.WithNamespace(istioNamespace), helm.WithChart(filepath.Join(istioHome, "manifests", "charts", "istio-control", "istio-discovery")), helm.WithWait(), helm.WithTimeout("10m"))
			if err != nil {
				t.Fatal("failed to invoke helm install operation istiod due to an error", err)
			}
			err = manager.RunInstall(helm.WithName("istio-ingress"), helm.WithNamespace(istioNamespace), helm.WithChart(filepath.Join(istioHome, "manifests", "charts", "gateways", "istio-ingress")), helm.WithWait(), helm.WithTimeout("10m"))
			if err != nil {
				t.Fatal("failed to invoke helm install operation istio-ingress due to an error", err)
			}

			log.Printf("Using kubeconfig file %s\n", config.KubeconfigFile())
			executeCommands(config.KubeconfigFile(), []string{
				// "wait --for=condition=ready pod --selector=app=nginx-nginx-ingress -n ingress-nginx --timeout=160s",
				// "patch svc nginx-nginx-ingress -n ingress-nginx --patch \"$(cat istio-filter-lab/patches/patch_ingresscontroller.yaml)\"",
				"apply -f istio-filter-lab/httpbin.yaml",
			}, t)

			return ctx
		}).Feature()

	testEnv.Test(t, feature)
}
