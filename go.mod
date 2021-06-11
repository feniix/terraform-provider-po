module github.com/feniix/terraform-provider-po

go 1.16

require (
	cloud.google.com/go v0.79.0 // indirect
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/fatih/color v1.9.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.6.1
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.48.1
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.48.1
	golang.org/x/oauth2 v0.0.0-20210323180902-22b0adad7558 // indirect
	google.golang.org/genproto v0.0.0-20210312152112-fc591d9ea70f // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/api v0.21.1
	k8s.io/apiextensions-apiserver v0.21.0 // indirect
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/kube-aggregator v0.21.1
	k8s.io/utils v0.0.0-20210305010621-2afb4311ab10 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.21.1
