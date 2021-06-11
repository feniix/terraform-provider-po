package po

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	po_types "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	pkgApi "k8s.io/apimachinery/pkg/types"
)

func resourcePoServiceMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoServiceMonitorCreate,
		ReadContext:   resourcePoServiceMonitorRead,
		UpdateContext: resourcePoServiceMonitorUpdate,
		DeleteContext: resourcePoServiceMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("service monitor", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_label": {
							Type:        schema.TypeString,
							Description: "The label to use to retrieve the job name from. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
							Optional:    true,
						},
						"target_labels": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "TargetLabels transfers labels on the Kubernetes Service onto the target. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
						},
						"pod_target_labels": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "PodTargetLabels transfers labels on the Kubernetes Pod onto the target. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
						},
						"endpoints": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of endpoints allowed as part of this ServiceMonitor. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
							Elem: &schema.Resource{
								Schema: endpointSchema(),
							},
						},
						"selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Selector to select Endpoints objects.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: labelSelectorFields(true),
							},
						},
						"namespace_selector": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Selector to select which namespaces the Endpoints objects are discovered from.",
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: namespaceSelectorSchema(),
							},
						},
						"sample_limit": {
							Type:        schema.TypeInt,
							Description: "SampleLimit defines per-scrape limit on number of scraped samples that will be accepted. More info: https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitorspec",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourcePoServiceMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).MonitoringClientset()
	if err != nil {
		return diag.FromErr(err)
	}
	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec, err := expandServiceMonitorSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	monitor := po_types.ServiceMonitor{
		ObjectMeta: metadata,
		Spec:       *spec,
	}
	log.Printf("[INFO] Creating new service monitor: %#v", monitor)
	out, err := conn.MonitoringV1().ServiceMonitors(metadata.Namespace).Create(ctx, &monitor, metav1.CreateOptions{})
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Submitted new service monitor: %#v", out)
	d.SetId(buildId(out.ObjectMeta))
	return resourcePoServiceMonitorRead(ctx, d, meta)
}

func resourcePoServiceMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	exists, err := resourcePoServiceMonitorExists(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !exists {
		d.SetId("")
		return diag.Diagnostics{}
	}
	conn, err := meta.(KubeClientsets).MonitoringClientset()
	if err != nil {
		return diag.FromErr(err)
	}
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Reading service monitor %s", name)
	sm, err := conn.MonitoringV1().ServiceMonitors(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Printf("[DEBUG] Received error: %#v", err)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Received service monitor: %#v", sm)
	err = d.Set("metadata", flattenMetadata(sm.ObjectMeta, d))
	if err != nil {
		return diag.FromErr(err)
	}
	spec, err := flattenServiceMonitorSpec(sm.Spec, d)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("spec", spec)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePoServiceMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).MonitoringClientset()
	if err != nil {
		return diag.FromErr(err)
	}
	namespace, name, err := idParts(d.Id())
	ops := patchMetadata("metaadata.0.", "/metadata/", d)

	if d.HasChange("spec") {
		spec, err := expandServiceMonitorSpec(d.Get("spec").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		ops = append(ops, replace(spec))
	}
	data, err := ops.MarshalJSON()
	if err != nil {
		return diag.FromErr(err)
	}
	out, err := conn.MonitoringV1().ServiceMonitors(namespace).Patch(ctx, name, pkgApi.JSONPatchType, data, metav1.PatchOptions{})
	if err != nil {
		return diag.Errorf("Failed to update Service Monitor: %s", err)
	}
	log.Printf("[INFO] Submitted updated config map: %#v", out)
	d.SetId(buildId(out.ObjectMeta))
	return resourcePoServiceMonitorRead(ctx, d, meta)
}

func resourcePoServiceMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).MonitoringClientset()
	if err != nil {
		return diag.FromErr(err)
	}
	namespace, name, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Deleting Aervice monitor: %#v", name)
	err = conn.MonitoringV1().ServiceMonitors(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Service monitor %s deleted", name)

	d.SetId("")
	return nil
}

func resourcePoServiceMonitorExists(ctx context.Context, d *schema.ResourceData, meta interface{}) (bool, error) {
	conn, err := meta.(KubeClientsets).MonitoringClientset()
	if err != nil {
		return false, err
	}

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] Checking service monitor %s", name)
	_, err = conn.MonitoringV1().ServiceMonitors(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}
	return true, err
}

func expandServiceMonitorSpec(sm []interface{}) (*po_types.ServiceMonitorSpec, error) {
	obj := &po_types.ServiceMonitorSpec{}
	if len(sm) == 0 || sm[0] == nil {
		return obj, nil
	}
	in := sm[0].(map[string]interface{})

	obj.JobLabel = in["job_label"].(string)
	obj.SampleLimit = uint64(in["sample_limit"].(int))
	if tl, ok := in["target_labels"].([]interface{}); ok {
		obj.TargetLabels = expandStringSlice(tl)
	}
	if ptl, ok := in["pod_target_labels"].([]interface{}); ok {
		obj.PodTargetLabels = expandStringSlice(ptl)
	}
	if v, ok := in["endpoints"].([]interface{}); ok && len(v) > 0 {
		endpoints, err := expandEndpoints(v)
		if err != nil {
			return obj, err
		}
		obj.Endpoints = endpoints
	}
	if s, ok := in["selector"].([]interface{}); ok && len(s) > 0 {
		selector := expandLabelSelector(s)
		obj.Selector = *selector
	}
	if ns, ok := in["namespace_selector"].([]interface{}); ok && len(ns) > 0 {
		selector, err := expandNamespaceSelector(ns)
		if err != nil {
			return obj, err
		}
		obj.NamespaceSelector = *selector
	}
	return obj, nil
}

func flattenServiceMonitorSpec(spec po_types.ServiceMonitorSpec, d *schema.ResourceData) ([]interface{}, error) {
	att := make(map[string]interface{})

	if spec.JobLabel != "" {
		att["job_label"] = spec.JobLabel
	}
	if len(spec.TargetLabels) > 0 {
		att["target_labels"] = spec.TargetLabels
	}
	if len(spec.PodTargetLabels) > 0 {
		att["pod_target_labels"] = spec.PodTargetLabels
	}
	att["sample_limit"] = int(spec.SampleLimit)

	endpoints, err := flattenEndpoints(spec.Endpoints)
	if err != nil {
		return nil, err
	}
	att["endpoints"] = endpoints
	att["selector"] = flattenLabelSelector(&spec.Selector)
	att["namespace_selector"] = flattenNamespaceSelector(&spec.NamespaceSelector)

	return []interface{}{att}, nil
}
