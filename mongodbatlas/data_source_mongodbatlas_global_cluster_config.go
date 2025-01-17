package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasGlobalCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasGlobalClusterRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"managed_namespaces": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collection": {
							Type:     schema.TypeString,
							Required: true,
						},
						"custom_shard_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"db": {
							Type:     schema.TypeString,
							Required: true,
						},
						"is_custom_shard_key_hashed": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"is_shard_key_unique": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"custom_zone_mapping": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasGlobalClusterRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	globalCluster, resp, err := conn.GlobalClusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	if err := d.Set("managed_namespaces", flattenManagedNamespaces(globalCluster.ManagedNamespaces)); err != nil {
		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	if err := d.Set("custom_zone_mapping", globalCluster.CustomZoneMapping); err != nil {
		return diag.FromErr(fmt.Errorf(errorGlobalClusterRead, clusterName, err))
	}

	d.SetId(clusterName)

	return nil
}
