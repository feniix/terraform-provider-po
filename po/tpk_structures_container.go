package po

import (
	v1 "k8s.io/api/core/v1"
)

func expandSecretKeyRef(r []interface{}) (*v1.SecretKeySelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.SecretKeySelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.SecretKeySelector{}

	if v, ok := in["key"].(string); ok {
		obj.Key = v
	}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	if v, ok := in["optional"]; ok {
		obj.Optional = ptrToBool(v.(bool))
	}
	return obj, nil
}

func expandConfigMapKeyRef(r []interface{}) (*v1.ConfigMapKeySelector, error) {
	if len(r) == 0 || r[0] == nil {
		return &v1.ConfigMapKeySelector{}, nil
	}
	in := r[0].(map[string]interface{})
	obj := &v1.ConfigMapKeySelector{}

	if v, ok := in["key"].(string); ok {
		obj.Key = v
	}
	if v, ok := in["name"].(string); ok {
		obj.Name = v
	}
	if v, ok := in["optional"]; ok {
		obj.Optional = ptrToBool(v.(bool))
	}
	return obj, nil

}

func flattenSecretKeyRef(in *v1.SecretKeySelector) []interface{} {
	att := make(map[string]interface{})

	if in.Key != "" {
		att["key"] = in.Key
	}
	if in.Name != "" {
		att["name"] = in.Name
	}
	if in.Optional != nil {
		att["optional"] = *in.Optional
	}
	return []interface{}{att}
}

func flattenConfigMapKeyRef(in *v1.ConfigMapKeySelector) []interface{} {
	att := make(map[string]interface{})

	if in.Key != "" {
		att["key"] = in.Key
	}
	if in.Name != "" {
		att["name"] = in.Name
	}
	if in.Optional != nil {
		att["optional"] = *in.Optional
	}
	return []interface{}{att}
}
