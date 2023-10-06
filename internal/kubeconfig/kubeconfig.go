// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubeconfig

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// TODO :: Bhargav :: Change name as per final usage
// type ReadWriteResetCloser interface{}

type Loader interface {
	Load(cfgPath string) ([]kubeconfigFile, error)
}

type Kubeconfig struct {
	loader Loader

	f           *kubeconfigFile
	rootNode    *yaml.RNode
	kubeCfgPath string
}

func (k *Kubeconfig) WithLoader(l Loader) *Kubeconfig {
	k.loader = l
	return k
}

func (k *Kubeconfig) Parse() error {
	cfgPath, err := kubeconfigPath()
	if err != nil {
		return errors.Wrap(err, "cannot determine kubeconfig path")
	}
	k.kubeCfgPath = cfgPath

	files, err := k.loader.Load(cfgPath)
	if err != nil {
		return errors.Wrap(err, "failed to load kubeconfig")
	}

	// TODO since we don't support multiple kubeconfig files at the moment, there's just 1 file
	k.f = &files[0]
	k.rootNode = k.f.node

	// Check if kubeconfig document is a map document
	_, err = k.rootNode.FieldRNodes()
	if err != nil {
		return errors.Wrap(err, "kubeconfig file is not a map document")
	}

	return nil
}

func (k *Kubeconfig) Bytes() ([]byte, error) {
	return yaml.Marshal(k.rootNode)
}

func (k *Kubeconfig) Save() error {
	return yaml.WriteFile(k.rootNode, k.kubeCfgPath)
}
