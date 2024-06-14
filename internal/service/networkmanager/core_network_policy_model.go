// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkmanager

import (
	"encoding/json"
	"sort"

	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

type coreNetworkPolicyDocument struct {
	Version                  string                                     `json:"version,omitempty"`
	CoreNetworkConfiguration *coreNetworkPolicyCoreNetworkConfiguration `json:"core-network-configuration"`
	Segments                 []*coreNetworkPolicySegment                `json:"segments"`
	// TODO NetworkFunctionGroups
	SegmentActions     []*coreNetworkPolicySegmentAction `json:"segment-actions,omitempty"`
	AttachmentPolicies []*coreNetworkAttachmentPolicy    `json:"attachment-policies,omitempty"`
}

type coreNetworkPolicyCoreNetworkConfiguration struct {
	AsnRanges        interface{}                                 `json:"asn-ranges"`
	InsideCidrBlocks interface{}                                 `json:"inside-cidr-blocks,omitempty"`
	VpnEcmpSupport   bool                                        `json:"vpn-ecmp-support"`
	EdgeLocations    []*coreNetworkPolicyCoreNetworkEdgeLocation `json:"edge-locations,omitempty"`
}

type coreNetworkPolicyCoreNetworkEdgeLocation struct {
	Location         string      `json:"location"`
	Asn              int64       `json:"asn,omitempty"`
	InsideCidrBlocks interface{} `json:"inside-cidr-blocks,omitempty"`
}

type coreNetworkPolicySegment struct {
	Name                        string      `json:"name"`
	Description                 string      `json:"description,omitempty"`
	EdgeLocations               interface{} `json:"edge-locations,omitempty"`
	IsolateAttachments          bool        `json:"isolate-attachments"`
	RequireAttachmentAcceptance bool        `json:"require-attachment-acceptance"`
	DenyFilter                  interface{} `json:"deny-filter,omitempty"`
	AllowFilter                 interface{} `json:"allow-filter,omitempty"`
}

type coreNetworkPolicySegmentAction struct {
	Action                string                                    `json:"action"`
	Segment               string                                    `json:"segment,omitempty"`
	Mode                  string                                    `json:"mode,omitempty"`
	ShareWith             interface{}                               `json:"share-with,omitempty"`
	ShareWithExcept       interface{}                               `json:"except,omitempty"`
	DestinationCidrBlocks interface{}                               `json:"destination-cidr-blocks,omitempty"`
	Destinations          interface{}                               `json:"destinations,omitempty"`
	Description           string                                    `json:"description,omitempty"`
	WhenSentTo            *coreNetworkPolicySegmentActionWhenSentTo `json:"when-sent-to,omitempty"`
	Via                   *coreNetworkPolicySegmentActionVia        `json:"via,omitempty"`
}

type coreNetworkPolicySegmentActionWhenSentTo struct {
	Segments interface{} `json:"segments,omitempty"`
}

type coreNetworkPolicySegmentActionVia struct {
	NetworkFunctionGroups interface{}                                      `json:"network-function-groups,omitempty"`
	WithEdgeOverrides     []*coreNetworkPolicySegmentActionViaEdgeOverride `json:"with-edge-overrides,omitempty"`
}
type coreNetworkPolicySegmentActionViaEdgeOverride struct {
	EdgeSets interface{} `json:"edge-sets,omitempty"`
	UseEdge  string      `json:"use-edge,omitempty"`
}

type coreNetworkAttachmentPolicy struct {
	RuleNumber     int                                     `json:"rule-number,omitempty"`
	Description    string                                  `json:"description,omitempty"`
	ConditionLogic string                                  `json:"condition-logic,omitempty"`
	Conditions     []*coreNetworkAttachmentPolicyCondition `json:"conditions"`
	Action         *coreNetworkAttachmentPolicyAction      `json:"action"`
}

type coreNetworkAttachmentPolicyCondition struct {
	Type     string `json:"type,omitempty"`
	Operator string `json:"operator,omitempty"`
	Key      string `json:"key,omitempty"`
	Value    string `json:"value,omitempty"`
}

type coreNetworkAttachmentPolicyAction struct {
	AssociationMethod         string `json:"association-method,omitempty"`
	Segment                   string `json:"segment,omitempty"`
	TagValueOfKey             string `json:"tag-value-of-key,omitempty"`
	RequireAcceptance         bool   `json:"require-acceptance,omitempty"`
	AddToNetworkFunctionGroup string `json:"add-to-network-function-group,omitempty"`
}

func (c coreNetworkPolicySegmentAction) MarshalJSON() ([]byte, error) {
	type Alias coreNetworkPolicySegmentAction
	var share interface{}

	if c.ShareWith != nil {
		sWIntf := c.ShareWith.([]string)

		if sWIntf[0] == "*" {
			share = sWIntf[0]
		} else {
			share = sWIntf
		}
	}

	if c.ShareWithExcept != nil {
		share = c.ShareWithExcept.([]string)
	}

	return json.Marshal(&Alias{
		Action:                c.Action,
		Mode:                  c.Mode,
		Destinations:          c.Destinations,
		DestinationCidrBlocks: c.DestinationCidrBlocks,
		Segment:               c.Segment,
		ShareWith:             share,
	})
}

func coreNetworkPolicyExpandStringList(configured []interface{}) interface{} {
	vs := flex.ExpandStringValueList(configured)
	sort.Sort(sort.Reverse(sort.StringSlice(vs)))

	return vs
}
