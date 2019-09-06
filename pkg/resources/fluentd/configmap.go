/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fluentd

import (
	"bytes"
	"fmt"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"text/template"
)

type fluentdTLSConfig struct {
	Enabled    bool
	SharedKey  string
	CACertFile string
	CertFile   string
	KeyFile    string
}

type fluentdConfig struct {
	TLS          fluentdTLSConfig
	LogJobOutput bool
}

func (r *Reconciler) configMap() runtime.Object {
	tlsConfig := fluentdTLSConfig{
		Enabled:   r.Fluentd.Spec.TLS.Enabled,
		SharedKey: r.Fluentd.Spec.TLS.SharedKey,
	}
	if r.Fluentd.Spec.TLS.SecretType == "tls" {
		tlsConfig.CertFile = "/fluentd/tls/tls.crt"
		tlsConfig.KeyFile = "/fluentd/tls/tls.key"
		tlsConfig.CACertFile = "/fluentd/tls/ca.crt"
	} else {
		tlsConfig.CertFile = "/fluentd/tls/serverCert"
		tlsConfig.KeyFile = "/fluentd/tls/serverKey"
		tlsConfig.CACertFile = "/fluentd/tls/caCert"
	}
	input := fluentdConfig{
		TLS:          tlsConfig,
		LogJobOutput: r.Fluentd.Spec.LogJobOutput,
	}
	return &corev1.ConfigMap{
		ObjectMeta: templates.FluentdObjectMeta(configMapName, util.MergeLabels(r.Fluentd.Labels, labelSelector), r.Fluentd),
		Data: map[string]string{
			"fluent.conf":  fluentdDefaultTemplate,
			"input.conf":   generateConfig(input),
			"devnull.conf": fluentdOutputTemplate,
		},
	}
}

func generateConfig(input fluentdConfig) string {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(fluentdInputTemplate)
	if err != nil {
		return ""
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return ""
	}
	outputString := fmt.Sprint(output.String())
	return outputString
}
