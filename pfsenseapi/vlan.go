package pfsenseapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/markphelps/optional"
	"strconv"
)

type vlansClient interface {
	GetVLAN(ctx context.Context, id int) (*VLAN, error)
	GetVLANs(ctx context.Context) ([]*VLAN, error)
	CreateVLAN(ctx context.Context, vlan VLANRequest) (*VLAN, error)
	DeleteVLAN(ctx context.Context, id int) error
	UpdateVLAN(ctx context.Context, id int, vlan VLANRequest) (*VLAN, error)
}

const (
	interfaceVLANEndpoint  = "api/v2/interface/vlan"
	interfaceVLANsEndpoint = "api/v2/interface/vlans"
)

// VLAN represents a single VLAN.
type VLAN struct {
	VLANRequest
	Id int `json:"id"`
}

type vlanListResponse struct {
	apiResponse
	Data []*VLAN `json:"data"`
}

// ListVLANs returns the VLANs
func (c *client) GetVLANs(ctx context.Context) ([]*VLAN, error) {
	response, err := c.get(ctx, interfaceVLANsEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp := new(vlanListResponse)
	if err = json.Unmarshal(response, resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return resp.Data, nil
}

// GetVLAN returns the VLAN with the given ID.
func (c *client) GetVLAN(ctx context.Context, id int) (*VLAN, error) {
	response, err := c.get(
		ctx,
		interfaceVLANEndpoint,
		map[string]string{
			"id": strconv.Itoa(id),
		},
	)
	if err != nil {
		return nil, err
	}

	resp := new(createVLANResponse)
	if err = json.Unmarshal(response, resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return resp.Data, nil
}

// DeleteVLAN deletes a VLAN.
func (c *client) DeleteVLAN(ctx context.Context, idToDelete int) error {
	_, err := c.delete(
		ctx,
		interfaceVLANEndpoint,
		map[string]string{
			"id": strconv.Itoa(idToDelete),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

type VLANRequest struct {
	If     string           `json:"if"`
	Tag    int              `json:"tag"`
	Vlanif *optional.String `json:"vlanif,omitempty"`
	Pcp    *optional.Int    `json:"pcp,omitempty"`
	Descr  *optional.String `json:"descr,omitempty"`
}

type createVLANResponse struct {
	apiResponse
	Data *VLAN `json:"data"`
}

// CreateVLAN creates a new VLAN.
func (c *client) CreateVLAN(
	ctx context.Context,
	newVLAN VLANRequest,
) (*VLAN, error) {
	jsonData, err := json.Marshal(newVLAN)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request payload into json: %w", err)
	}

	response, err := c.post(ctx, interfaceVLANEndpoint, nil, jsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	resp := new(createVLANResponse)
	if err = json.Unmarshal(response, resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return resp.Data, nil
}

// UpdateVLAN modifies an existing VLAN.
func (c *client) UpdateVLAN(
	ctx context.Context,
	idToUpdate int,
	vlanData VLANRequest,
) (*VLAN, error) {
	requestData := VLAN{
		VLANRequest: vlanData,
		Id:          idToUpdate,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request payload into json: %w", err)
	}

	response, err := c.patch(ctx, interfaceVLANEndpoint, nil, jsonData)
	if err != nil {
		return nil, err
	}

	resp := new(createVLANResponse)
	if err = json.Unmarshal(response, resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return resp.Data, nil
}
