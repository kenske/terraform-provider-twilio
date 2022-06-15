/*
 * Twilio - Supersim
 *
 * This is the public Twilio REST API.
 *
 * API version: 1.29.2
 * Contact: support@twilio.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/twilio/terraform-provider-twilio/client"
	. "github.com/twilio/terraform-provider-twilio/core"
	. "github.com/twilio/twilio-go/rest/supersim/v1"
)

func ResourceNetworkAccessProfilesNetworks() *schema.Resource {
	return &schema.Resource{
		CreateContext: createNetworkAccessProfilesNetworks,
		ReadContext:   readNetworkAccessProfilesNetworks,
		DeleteContext: deleteNetworkAccessProfilesNetworks,
		Schema: map[string]*schema.Schema{
			"network_access_profile_sid": AsString(SchemaForceNewRequired),
			"network":                    AsString(SchemaForceNewRequired),
			"sid":                        AsString(SchemaComputed),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				err := parseNetworkAccessProfilesNetworksImportId(d.Id(), d)
				if err != nil {
					return nil, err
				}

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func createNetworkAccessProfilesNetworks(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	params := CreateNetworkAccessProfileNetworkParams{}
	if err := UnmarshalSchema(&params, d); err != nil {
		return diag.FromErr(err)
	}

	networkAccessProfileSid := d.Get("network_access_profile_sid").(string)

	r, err := m.(*client.Config).Client.SupersimV1.CreateNetworkAccessProfileNetwork(networkAccessProfileSid, &params)
	if err != nil {
		return diag.FromErr(err)
	}

	idParts := []string{networkAccessProfileSid}
	idParts = append(idParts, (*r.Sid))
	d.SetId(strings.Join(idParts, "/"))

	err = MarshalSchema(d, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteNetworkAccessProfilesNetworks(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	networkAccessProfileSid := d.Get("network_access_profile_sid").(string)
	sid := d.Get("sid").(string)

	err := m.(*client.Config).Client.SupersimV1.DeleteNetworkAccessProfileNetwork(networkAccessProfileSid, sid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func readNetworkAccessProfilesNetworks(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	networkAccessProfileSid := d.Get("network_access_profile_sid").(string)
	sid := d.Get("sid").(string)

	r, err := m.(*client.Config).Client.SupersimV1.FetchNetworkAccessProfileNetwork(networkAccessProfileSid, sid)
	if err != nil {
		return diag.FromErr(err)
	}

	err = MarshalSchema(d, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func parseNetworkAccessProfilesNetworksImportId(importId string, d *schema.ResourceData) error {
	importParts := strings.Split(importId, "/")
	errStr := "invalid import ID (%q), expected network_access_profile_sid/sid"

	if len(importParts) != 2 {
		return fmt.Errorf(errStr, importId)
	}

	d.Set("network_access_profile_sid", importParts[0])
	d.Set("sid", importParts[1])

	return nil
}
