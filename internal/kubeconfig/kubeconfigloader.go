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
	"os"
	"path/filepath"

	"github.com/ahmetb/kubectx/internal/cmdutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"github.com/pkg/errors"
)

var (
	DefaultLoader Loader = new(StandardKubeconfigLoader)
)

type StandardKubeconfigLoader struct{}

type kubeconfigFile struct {
	node *yaml.RNode
}

// TODO :: Bhargav :: Replace []kubeconfigFile with slice of some interface finally
func (*StandardKubeconfigLoader) Load(cfgPath string) ([]kubeconfigFile, error) {
	node, err := yaml.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrap(err, "kubeconfig file not found")
		}
		return nil, errors.Wrap(err, "failed to open file")
	}

	// TODO we'll return all kubeconfig files when we start implementing multiple kubeconfig support
	return []kubeconfigFile{{node}}, nil
}

func kubeconfigPath() (string, error) {
	// KUBECONFIG env var
	if v := os.Getenv("KUBECONFIG"); v != "" {
		list := filepath.SplitList(v)
		if len(list) > 1 {
			// TODO KUBECONFIG=file1:file2 currently not supported
			return "", errors.New("multiple files in KUBECONFIG are currently not supported")
		}
		return v, nil
	}

	// default path
	home := cmdutil.HomeDir()
	if home == "" {
		return "", errors.New("HOME or USERPROFILE environment variable not set")
	}
	return filepath.Join(home, ".kube", "config"), nil
}
