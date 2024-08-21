/*
Copyright 2018 The Kubernetes Authors.

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

package rules

import (
	"reflect"
	"strings"

	"k8s.io/gengo/v2/types"
)

/*
NamingConvention implements APIRule interface. It checks cases that are ignored
by the NamesMatch APIRule.

The embedded metav1.ListMeta and metav1.ObjectMeta Go fields must have the "metadata" JSON field name.
*/
type NamingConvention struct{}

// Name returns the name of APIRule
func (n *NamingConvention) Name() string {
	return "naming_convention"
}

// Validate evaluates API rule on type t and returns a list of field names in
// the type that violate the rule. Empty field name [""] implies the entire
// type violates the rule.
func (n *NamingConvention) Validate(t *types.Type) ([]string, error) {
	fields := make([]string, 0)

	// Only validate struct type and ignore the rest
	switch t.Kind {
	case types.Struct:
		for _, m := range t.Members {
			if !matchesNamingConvention(m) {
				fields = append(fields, m.Name)
			}
		}
	}
	return fields, nil
}

func matchesNamingConvention(m types.Member) bool {
	if !m.Embedded {
		// Check only applies to embedded types.
		return true
	}
	typeName := m.Type.String()
	switch typeName {
	case "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta",
		"k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta":
		jsonTag, ok := reflect.StructTag(m.Tags).Lookup("json")
		if !ok {
			return false
		}
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName != "metadata" {
			return false
		}
	}

	return true
}
