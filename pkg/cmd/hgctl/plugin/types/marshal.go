// Copyright (c) 2022 Alibaba Group Holding Ltd.
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

package types

import (
	"bytes"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/json"
)

func (s JSON) MarshalJSON() ([]byte, error) {
	if len(s.Raw) > 0 {
		var obj interface{}
		err := json.Unmarshal(s.Raw, &obj)
		if err != nil {
			return []byte("null"), err
		}
		return json.Marshal(obj)
	}
	return []byte("null"), nil
}

func (s *JSON) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && !bytes.Equal(data, []byte("null")) {
		s.Raw = data
	}
	return nil
}

func (s JSON) MarshalYAML() (interface{}, error) {
	if len(s.Raw) > 0 {
		var obj interface{}
		err := yaml.Unmarshal(s.Raw, &obj)
		if err != nil {
			return "null", err
		}
		return obj, nil
	}
	return "null", nil
}

func (s JSONSchemaPropsOrArray) MarshalJSON() ([]byte, error) {
	if len(s.JSONSchemas) > 0 {
		return json.Marshal(s.JSONSchemas)
	}
	return json.Marshal(s.Schema)
}

func (s *JSONSchemaPropsOrArray) UnmarshalJSON(data []byte) error {
	var nw JSONSchemaPropsOrArray
	var first byte
	if len(data) > 1 {
		first = data[0]
	}
	if first == '{' {
		var sch JSONSchemaProps
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.JSONSchemas); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}

func (s JSONSchemaPropsOrArray) MarshalYAML() (interface{}, error) {
	if len(s.JSONSchemas) > 0 {
		return s.JSONSchemas, nil
	}
	return s.Schema, nil
}

func (s JSONSchemaPropsOrBool) MarshalJSON() ([]byte, error) {
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}

	if s.Schema == nil && !s.Allows {
		return []byte("false"), nil
	}
	return []byte("true"), nil
}

func (s *JSONSchemaPropsOrBool) UnmarshalJSON(data []byte) error {
	var nw JSONSchemaPropsOrBool
	switch {
	case len(data) == 0:
	case data[0] == '{':
		var sch JSONSchemaProps
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		nw.Allows = true
		nw.Schema = &sch
	case len(data) == 4 && string(data) == "true":
		nw.Allows = true
	case len(data) == 5 && string(data) == "false":
		nw.Allows = false
	default:
		return errors.New("boolean or JSON schema expected")
	}
	*s = nw
	return nil
}

func (s JSONSchemaPropsOrBool) MarshalYAML() (interface{}, error) {
	if s.Schema != nil {
		return yaml.Marshal(s.Schema)
	}

	if s.Schema == nil && !s.Allows {
		return false, nil
	}
	return true, nil
}

func (s JSONSchemaPropsOrStringArray) MarshalJSON() ([]byte, error) {
	if len(s.Property) > 0 {
		return json.Marshal(s.Property)
	}
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}
	return []byte("null"), nil
}

func (s *JSONSchemaPropsOrStringArray) UnmarshalJSON(data []byte) error {
	var first byte
	if len(data) > 1 {
		first = data[0]
	}
	var nw JSONSchemaPropsOrStringArray
	if first == '{' {
		var sch JSONSchemaProps
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.Property); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}

func (s JSONSchemaPropsOrStringArray) MarshalYAML() (interface{}, error) {
	if len(s.Property) > 0 {
		return s.Property, nil
	}
	if s.Schema != nil {
		return s.Schema, nil
	}
	return "null", nil
}
