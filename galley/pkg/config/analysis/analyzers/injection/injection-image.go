// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package injection

import (
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"

	"istio.io/istio/galley/pkg/config/analysis"
	"istio.io/istio/galley/pkg/config/analysis/msg"
	"istio.io/istio/pkg/config/resource"
	"istio.io/istio/pkg/config/schema/collection"
	"istio.io/istio/pkg/config/schema/collections"
)

// ImageAnalyzer checks the image of auto-injection configured with the running proxies on pods.
type ImageAnalyzer struct{}

var _ analysis.Analyzer = &ImageAnalyzer{}

const sidecarInjectorConfigName = "istio-sidecar-injector"

// injectionConfigMap is a snippet of the sidecar injection ConfigMap
type injectionConfigMap struct {
	Global global `json:"global"`
}

type global struct {
	Hub   string `json:"hub"`
	Tag   string `json:"tag"`
	Proxy proxy  `json:"proxy"`
}

type proxy struct {
	Image string `json:"image"`
}

// Metadata implements Analyzer.
func (a *ImageAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "injection.ImageAnalyzer",
		Description: "Checks the image of auto-injection configured with the running proxies on pods",
		Inputs: collection.Names{
			collections.K8SCoreV1Namespaces.Name(),
			collections.K8SCoreV1Pods.Name(),
			collections.K8SCoreV1Configmaps.Name(),
		},
	}
}

// Analyze implements Analyzer.
func (a *ImageAnalyzer) Analyze(c analysis.Context) {
	var proxyImage string

	// TODO: when multiple injector configmaps exist, we may need to assess them respectively.
	c.ForEach(collections.K8SCoreV1Configmaps.Name(), func(r *resource.Instance) bool {
		if r.Metadata.FullName.Name.String() == sidecarInjectorConfigName {
			cm := r.Message.(*v1.ConfigMap)

			proxyImage = getIstioProxyImage(cm)

			return false
		}
		return true
	})

	if proxyImage == "" {
		return
	}

	injectedNamespaces := make(map[string]struct{})

	// Collect the list of namespaces that have istio injection enabled.
	c.ForEach(collections.K8SCoreV1Namespaces.Name(), func(r *resource.Instance) bool {
		if r.Metadata.Labels[InjectionLabelName] == InjectionLabelEnableValue {
			injectedNamespaces[r.Metadata.FullName.String()] = struct{}{}
		}

		return true
	})

	c.ForEach(collections.K8SCoreV1Pods.Name(), func(r *resource.Instance) bool {
		pod := r.Message.(*v1.Pod)

		if _, ok := injectedNamespaces[pod.GetNamespace()]; !ok {
			return true
		}

		// If the pod has been annotated with a custom sidecar, then ignore as
		// it always overrides the injector logic.
		if r.Metadata.Annotations["sidecar.istio.io/proxyImage"] != "" {
			return true
		}

		for _, container := range pod.Spec.Containers {
			if container.Name != istioProxyName {
				continue
			}

			if container.Image != proxyImage {
				c.Report(collections.K8SCoreV1Pods.Name(), msg.NewIstioProxyImageMismatch(r, container.Image, proxyImage))
			}
		}

		return true
	})
}

// getIstioProxyImage retrieves the proxy image name defined in the sidecar injector
// configuration.
func getIstioProxyImage(cm *v1.ConfigMap) string {
	var m injectionConfigMap
	if err := json.Unmarshal([]byte(cm.Data["values"]), &m); err != nil {
		return ""
	}
	return fmt.Sprintf("%s/%s:%s", m.Global.Hub, m.Global.Proxy.Image, m.Global.Tag)
}
