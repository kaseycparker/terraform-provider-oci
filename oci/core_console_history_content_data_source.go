// Copyright (c) 2017, 2019, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	oci_core "github.com/oracle/oci-go-sdk/core"
)

func CoreConsoleHistoryContentDataSource() *schema.Resource {
	return &schema.Resource{
		Read: readSingularCoreConsoleHistoryContent,
		Schema: map[string]*schema.Schema{
			"console_history_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
				// GetConsoleHistoryContent returns an error with length < 10240, though this is not documented in the API doc.
				ValidateFunc: func(i interface{}, k string) (s []string, es []error) {
					v, ok := i.(int)
					if !ok {
						es = append(es, fmt.Errorf("expected type of %s to be int", k))
						return
					}

					if v < 10240 {
						es = append(es, fmt.Errorf("expected %s to be less than %d, got %d", k, 10240, v))
						return
					}

					return
				},
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Computed
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func readSingularCoreConsoleHistoryContent(d *schema.ResourceData, m interface{}) error {
	sync := &CoreConsoleHistoryContentDataSourceCrud{}
	sync.D = d
	sync.Client = m.(*OracleClients).computeClient

	return ReadResource(sync)
}

type CoreConsoleHistoryContentDataSourceCrud struct {
	D      *schema.ResourceData
	Client *oci_core.ComputeClient
	Res    *oci_core.GetConsoleHistoryContentResponse
}

func (s *CoreConsoleHistoryContentDataSourceCrud) VoidState() {
	s.D.SetId("")
}

func (s *CoreConsoleHistoryContentDataSourceCrud) Get() error {
	request := oci_core.GetConsoleHistoryContentRequest{}

	if consoleHistoryId, ok := s.D.GetOkExists("console_history_id"); ok {
		tmp := consoleHistoryId.(string)
		request.InstanceConsoleHistoryId = &tmp
	}

	if length, ok := s.D.GetOkExists("length"); ok {
		tmp := length.(int)
		request.Length = &tmp
	}

	if offset, ok := s.D.GetOkExists("offset"); ok {
		tmp := offset.(int)
		request.Offset = &tmp
	}

	request.RequestMetadata.RetryPolicy = getRetryPolicy(false, "core")

	response, err := s.Client.GetConsoleHistoryContent(context.Background(), request)
	if err != nil {
		return err
	}

	s.Res = &response
	return nil
}

func (s *CoreConsoleHistoryContentDataSourceCrud) SetData() error {
	if s.Res == nil {
		return nil
	}

	s.D.SetId(GenerateDataSourceID())

	if s.Res.Value != nil {
		s.D.Set("data", *s.Res.Value)
	}

	return nil
}
