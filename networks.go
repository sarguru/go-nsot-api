package nsot

import (
	"encoding/json"
	"fmt"
)

//Create Site structure to create site
type NetworkOpts struct {
	Cidr       string
	SiteId     int
	Attributes map[string]interface{}
	State      string
}

type Network struct {
	NetAddr    string                 `json:"network_address"`
	State      string                 `json:"state"`
	PrefixLen  int                    `json:"prefix_length"`
	IsIP       bool                   `json:"is_ip"`
	IpVersion  string                 `json:"ip_version"`
	Attributes map[string]interface{} `json:"attributes"`
	SiteId     int                    `json:"site_id"`
	Id         int                    `json:"id"`
}

type NetResponse struct {
	DataResp struct {
		Network Network `json:"network"`
	} `json:"data"`
	Status string `json:"status"`
}

type NetworksResponse struct {
	DataResp struct {
		Networks []Network `json:"networks"`
	} `json:"data"`
	Status string `json:"status"`
}

// get network data
func (r *NetResponse) GetNetwork() (*Network, error) {
	if r.Status != "ok" {
		return nil, fmt.Errorf("API Error: %s", r.Status)
	}

	return &r.DataResp.Network, nil
}

// get network id
func (r *NetworksResponse) GetId() (int, error) {
	if r.Status != "ok" {
		return 0, fmt.Errorf("API Error: %s", r.Status)
	}

	network := &r.DataResp.Networks[0]
	return network.Id, nil
}

//Retrieve Network by ID
func (c *Client) RetrieveNetworkbyID(id int) (*Network, error) {
	action := fmt.Sprintf("networks/%d/", id)
	req, err := c.NewRequest(action, "GET", "")

	if err != nil {
		return nil, fmt.Errorf("Error crafting request: %s", err)
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return nil, fmt.Errorf("Error Retrieving Network: %s", err)
	}

	netResp := new(NetResponse)
	err = decodeBody(resp, &netResp)

	if err != nil {
		return nil, fmt.Errorf("Error parsing network response: %s", err)
	}

	net, err := netResp.GetNetwork()

	if err != nil {
		return nil, fmt.Errorf("API error: %s", err)
	}

	return net, nil
}

//Retrieve Network ID
func (c *Client) RetrieveNetworkIdByCIDR(name string) (int, error) {
	action := fmt.Sprintf("networks/?network_address=%s", name)
	req, err := c.NewRequest(action, "GET", "")
	if err != nil {
		return 0, fmt.Errorf("Error crafting request: %s", err)
	}
	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return 0, fmt.Errorf("Error Retrieving network Id: %s", err)
	}
	netsResp := new(NetworksResponse)
	err = decodeBody(resp, &netsResp)
	if err != nil {
		return 0, fmt.Errorf("Error parsing site response: %s", err)
	}

	netId, err := netsResp.GetId()

	if err != nil {
		return 0, fmt.Errorf("API error: %s", err)
	}
	return netId, nil

}

//Retrieve Network by CIDR
func (c *Client) RetrieveNetbyName(cidr string) (*Network, error) {
	netId, err := c.RetrieveNetworkIdByCIDR(cidr)
	action := fmt.Sprintf("networks/%d/", netId)
	req, err := c.NewRequest(action, "GET", "")

	if err != nil {
		return nil, fmt.Errorf("Error crafting request: %s", err)
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return nil, fmt.Errorf("Error Retrieving network: %s", err)
	}

	netResp := new(NetResponse)
	err = decodeBody(resp, &netResp)

	if err != nil {
		return nil, fmt.Errorf("Error parsing site response: %s", err)
	}

	net, err := netResp.GetNetwork()

	if err != nil {
		return nil, fmt.Errorf("API error: %s", err)
	}

	return net, nil
}

// Function that actually creates network on NSOT
func (c *Client) CreateNetwork(opts *NetworkOpts) (*Network, error) {
	//Create Map of data to send to API
	data := make(map[string]interface{})

	data["cidr"] = opts.Cidr
	data["state"] = opts.State
	data["site_id"] = opts.SiteId
	data["attributes"] = opts.Attributes
	jsonData, err := json.Marshal(data)

	req, err := c.NewRequest("networks/", "POST", string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Error crafting request: %s", err)
	}
	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, fmt.Errorf("Error creating Network: %s", err)
	}
	netResp := new(NetResponse)
	err = decodeBody(resp, &netResp)

	if err != nil {
		return nil, fmt.Errorf("Error parsing Network response: %s", err)
	}
	net, err := netResp.GetNetwork()
	if err != nil {
		return nil, fmt.Errorf("API error: %s", err)
	}
	return net, nil
}

//Destroy Network by CIDR
func (c *Client) DestroyNetworkbyCIDR(cidr string) error {
	netId, err := c.RetrieveNetworkIdByCIDR(cidr)
	if err != nil {
		return fmt.Errorf("Error retrieving ID: %s", err)
	}
	action := fmt.Sprintf("networks/%d/", netId)
	req, err := c.NewRequest(action, "DELETE", "")
	if err != nil {
		return fmt.Errorf("Error crafting request: %s", err)
	}
	_, err = checkResp(c.Http.Do(req))
	if err != nil {
		return fmt.Errorf("Error Deleting Network: %s", err)
	}
	return nil
}

//Destroy Site by ID
func (c *Client) DestroyNetworkbyID(id int) error {
	action := fmt.Sprintf("networks/%d/", id)
	req, err := c.NewRequest(action, "DELETE", "")
	if err != nil {
		return fmt.Errorf("Error crafting request: %s", err)
	}
	_, err = checkResp(c.Http.Do(req))
	if err != nil {
		return fmt.Errorf("Error Deleting Network: %s", err)
	}
	return nil
}
