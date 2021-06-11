# This is a terraform provider for Prometheus operator deployments on Kubernetes

Most of the code is "stolen" from terraform-provider-kubernetes and https://github.com/greg-gajda/terraform-provider-po but upgrade the terraform-plugin-sdk to v2

The main reason for this provider creation is that I was not able to build https://github.com/greg-gajda/terraform-provider-po, I wanted to upgrade the provider to the plugin-sdk v2 to see if that would enable the terraform-ls to do autocompletion of the types, but after cloning the repo from Greg and trying multiple different `go mod` invocations I gave up and this is the result.

For now only support service monitors.

All of this is hardly tested, but generally speaking it works, you can deploy service monitors with it.

To see it in action, i.e. do a local test:

Install kind & terraform

```shell
brew install kind
brew install terraform
```

Install the provider

```shell
make install
```

It should do something like:

```shell
GOOS=darwin GOARCH=amd64 go build -o out/terraform-provider-po-darwin-amd64 ./cmd/terraform-provider-po
mkdir -p ~/.terraform.d/plugins/github.com/feniix/po/0.0.1/darwin_amd64
mv ./out/terraform-provider-po-darwin-amd64 ~/.terraform.d/plugins/github.com/feniix/po/0.0.1/darwin_amd64/terraform-provider-po
```
Start the cluster and deploy prometheus operator

```shell
make setup-cluster
```

It should do something similar to this:

```shell
kind create cluster --config ./config/kind-config.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.19.11) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
make install-operator
make[1]: Entering directory '/Users/otaegui/src/terraform-provider-po'
kubectl apply -f config/bundle.yaml
customresourcedefinition.apiextensions.k8s.io/alertmanagerconfigs.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/alertmanagers.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/podmonitors.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/probes.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/prometheuses.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/prometheusrules.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/servicemonitors.monitoring.coreos.com created
customresourcedefinition.apiextensions.k8s.io/thanosrulers.monitoring.coreos.com created
clusterrolebinding.rbac.authorization.k8s.io/prometheus-operator created
clusterrole.rbac.authorization.k8s.io/prometheus-operator created
deployment.apps/prometheus-operator created
serviceaccount/prometheus-operator created
service/prometheus-operator created
make[1]: Leaving directory '/Users/otaegui/src/terraform-provider-po'
```

Head down to ./examples/simple-service-monitor and run

```shell
$ terraform init

Initializing the backend...

Initializing provider plugins...
- Finding latest version of github.com/feniix/po...
- Using github.com/feniix/po v0.0.1 from the shared cache directory

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

Then 

```shell
$ terraform apply

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # po_service_monitor.example will be created
  + resource "po_service_monitor" "example" {
      + id = (known after apply)

      + metadata {
          + generation       = (known after apply)
          + labels           = {
              + "k8s-app" = "label1"
            }
          + name             = "example"
          + namespace        = "default"
          + resource_version = (known after apply)
          + uid              = (known after apply)
        }

      + spec {
          + job_label = "myapp"

          + endpoints {
              + honor_timestamps = true
              + interval         = "30s"
              + port             = "http-metrics"
            }

          + namespace_selector {
              + match_names = [
                  + "default",
                ]
            }

          + selector {
              + match_labels = {
                  + "k8s-app" = "myapplabel"
                }
            }
        }
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

po_service_monitor.example: Creating...
po_service_monitor.example: Creation complete after 0s [id=default/example]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```
