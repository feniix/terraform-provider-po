package po

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func endpointSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"port": {
			Type:        schema.TypeString,
			Description: "Name of the port this endpoint refers to. Mutually exclusive with targetPort. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
			Optional:    true,
		},
		"path": {
			Type: schema.TypeString,
			Description: "HTTP path to scrape for metrics. 	More info https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
			Optional: true,
		},
		"scheme": {
			Type:        schema.TypeString,
			Description: "HTTP scheme to use for scraping.",
			Optional:    true,
		},
		"interval": {
			Type:        schema.TypeString,
			Description: "Interval at which metrics should be scraped.",
			Optional:    true,
		},
		"scrape_timeout": {
			Type:        schema.TypeString,
			Description: "Timeout after which the scrape is ended.",
			Optional:    true,
		},
		"tls_config": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "TLS configuration to use when scraping the endpoint. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#tlsconfig",
			Elem: &schema.Resource{
				Schema: tlsConfigSchema(),
			},
		},
		"bearer_token_file": {
			Type:        schema.TypeString,
			Description: "File to read bearer token for scraping targets.",
			Optional:    true,
		},
		"bearer_token_secret": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Secret to mount to read bearer token for scraping targets. The secret needs to be in the same namespace as the service monitor and accessible by the Prometheus Operator.",
			Elem: &schema.Resource{
				Schema: secretKeySelectorSchema(),
			},
		},
		"honor_labels": {
			Type:        schema.TypeBool,
			Description: "HonorLabels chooses the metric's labels on collisions with target labels.",
			Optional:    true,
		},
		"honor_timestamps": {
			Type:        schema.TypeBool,
			Description: "HonorTimestamps controls whether Prometheus respects the timestamps present in scraped data.",
			Optional:    true,
			Default:     true,
		},
		"basic_auth": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "BasicAuth allow an endpoint to authenticate over basic authentication More info: https://prometheus.io/docs/operating/configuration/#endpoints",
			Elem: &schema.Resource{
				Schema: basicAuthSchema(),
			},
		},
		"metric_relabelings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "MetricRelabelConfigs to apply to samples before ingestion.",
			Elem: &schema.Resource{
				Schema: relabelConfigSchema(),
			},
		},
		"relabelings": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "RelabelConfigs to apply to samples before scraping. More info: https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config",
			Elem: &schema.Resource{
				Schema: relabelConfigSchema(),
			},
		},
		"proxy_url": {
			Type:        schema.TypeString,
			Description: "ProxyURL eg http://proxyserver:2195 Directs scrapes to proxy through this endpoint.",
			Optional:    true,
		},
	}
}

func relabelConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"separator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Separator placed between concatenated source label values. default is ';'.",
		},
		"target_label": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Label to which the resulting value is written in a replace action. It is mandatory for replace actions. Regex capture groups are available.",
		},
		"regex": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Regular expression against which the extracted value is matched. Default is '(.*)'",
		},
		"modulus": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Modulus to take of the hash of the source label values.",
		},
		"replacement": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Replacement value against which a regex replace is performed if the regular expression matches. Regex capture groups are available. Default is '$1'",
		},
		"action": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Action to perform based on regex matching. Default is 'replace'",
		},
		"source_labels": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The source labels select values from existing labels. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#relabelconfig",
		},
	}
}

func basicAuthSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"username": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The secret in the service monitor namespace that contains the username for authentication.",
			Elem: &schema.Resource{
				Schema: secretKeySelectorSchema(),
			},
		},
		"password": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "The secret in the service monitor namespace that contains the password for authentication.",
			Elem: &schema.Resource{
				Schema: secretKeySelectorSchema(),
			},
		},
	}
}

func secretKeySelectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The key of the secret to select from. Must be a valid secret key.",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
		},
	}
}

func tlsConfigSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ca_file": {
			Type:        schema.TypeString,
			Description: "Path to the CA cert in the Prometheus container to use for the targets.",
			Optional:    true,
		},
		"ca": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Struct containing the CA cert to use for the targets.",
			Elem: &schema.Resource{
				Schema: secretOrConfigMapSchema(),
			},
		},
		"cert_file": {
			Type:        schema.TypeString,
			Description: "Path to the client cert file in the Prometheus container for the targets.",
			Optional:    true,
		},
		"cert": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Struct containing the client cert file for the targets.",
			Elem: &schema.Resource{
				Schema: secretOrConfigMapSchema(),
			},
		},
		"key_file": {
			Type:        schema.TypeString,
			Description: "Path to the client key file in the Prometheus container for the targets.",
			Optional:    true,
		},
		"key_secret": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Secret containing data to use for the targets",
			Elem: &schema.Resource{
				Schema: secretKeySelectorSchema(),
			},
		},
		"server_name": {
			Type:        schema.TypeString,
			Description: "Used to verify the hostname for the targets.",
			Optional:    true,
		},
		"insecure_skip_verify": {
			Type:        schema.TypeBool,
			Description: "Disable target certificate validation.",
			Optional:    true,
		},
	}
}

func secretOrConfigMapSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"secret": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Secret containing data to use for the targets",
			Elem: &schema.Resource{
				Schema: secretKeySelectorSchema(),
			},
		},
		"config_map": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Selects a key of a ConfigMap.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The key to select.",
					},
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of the referent. More info: http://kubernetes.io/docs/user-guide/identifiers#names",
					},
				},
			},
		},
	}
}

func namespaceSelectorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"any": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Boolean describing whether all namespaces are selected in contrast to a list restricting them.",
		},
		"match_names": {
			Type:        schema.TypeSet,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "List of namespace names",
		},
	}
}

func labelSelectorFields(updatable bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"match_expressions": {
			Type:        schema.TypeList,
			Description: "A list of label selector requirements. The requirements are ANDed.",
			Optional:    true,
			ForceNew:    !updatable,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Description: "The label key that the selector applies to.",
						Optional:    true,
						ForceNew:    !updatable,
					},
					"operator": {
						Type:        schema.TypeString,
						Description: "A key's relationship to a set of values. Valid operators ard `In`, `NotIn`, `Exists` and `DoesNotExist`.",
						Optional:    true,
						ForceNew:    !updatable,
					},
					"values": {
						Type:        schema.TypeSet,
						Description: "An array of string values. If the operator is `In` or `NotIn`, the values array must be non-empty. If the operator is `Exists` or `DoesNotExist`, the values array must be empty. This array is replaced during a strategic merge patch.",
						Optional:    true,
						ForceNew:    !updatable,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Set:         schema.HashString,
					},
				},
			},
		},
		"match_labels": {
			Type:        schema.TypeMap,
			Description: "A map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of `match_expressions`, whose key field is \"key\", the operator is \"In\", and the values array contains only \"value\". The requirements are ANDed.",
			Optional:    true,
			ForceNew:    !updatable,
		},
	}
}
