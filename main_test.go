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
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
)

var (
	testEnv          env.Environment
	namespace        string
	ingressNamespace string
	kindClusterName  string
	istioNamespace   string
	istioHome        string
)

func TestMain(m *testing.M) {

	cfg, _ := envconf.NewFromFlags()
	testEnv = env.NewWithConfig(cfg)
	kindClusterName = envconf.RandomName("third-party", 16)
	namespace = envconf.RandomName("third-party", 16)
	istioNamespace = "istio-system"
	ingressNamespace = "ingress-nginx"
	istioHome = os.Getenv("ISTIO_HOME")

	testEnv.Setup(
		envfuncs.CreateKindCluster(kindClusterName),
		envfuncs.CreateNamespace(namespace),
		envfuncs.CreateNamespace(istioNamespace),
		envfuncs.CreateNamespace(ingressNamespace),
	)

	testEnv.Finish(
	// envfuncs.DeleteNamespace(istioNamespace),
	// envfuncs.DeleteNamespace(namespace),
	// envfuncs.DestroyKindCluster(kindClusterName),
	)
	os.Exit(testEnv.Run(m))
}
