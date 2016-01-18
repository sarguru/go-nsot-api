package nsot

import (
	"encoding/json"
	"fmt"
)

//Create Site structure to create site
type SiteOpts struct {
	Name string
	Desc string
}

type Site struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"description"`
}

type SiteResponse struct {
	DataResp struct {
		Site Site `json:"site"`
	} `json:"data"`
	Status string `json:"status"`
}

type SitesResponse struct {
	DataResp struct {
		Sites []Site `json:"sites"`
	} `json:"data"`
	Status string `json:"status"`
}

// get site data
func (r *SiteResponse) GetSite() (*Site, error) {
	if r.Status != "ok" {
		return nil, fmt.Errorf("API Error: %s", r.Status)
	}

	return &r.DataResp.Site, nil
}

// get site id
func (r *SitesResponse) GetId() (int, error) {
	if r.Status != "ok" {
		return 0, fmt.Errorf("API Error: %s", r.Status)
	}

	site := &r.DataResp.Sites[0]
	return site.Id, nil
}

// Function that actually creates site on NSOT
func (c *Client) CreateSite(opts *SiteOpts) (*Site, error) {
	//Create Map of data to send to API
	data := make(map[string]string)

	data["name"] = opts.Name

	if opts.Desc != "" {
		data["description"] = opts.Desc
	}

	jsonData, err := json.Marshal(data)

	req, err := c.NewRequest("sites/", "POST", string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Error crafting request: %s", err)
	}
	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, fmt.Errorf("Error creating site: %s", err)
	}
	siteResp := new(SiteResponse)
	err = decodeBody(resp, &siteResp)

	if err != nil {
		return nil, fmt.Errorf("Error parsing site response: %s", err)
	}
	site, err := siteResp.GetSite()
	if err != nil {
		return nil, fmt.Errorf("API error: %s", err)
	}
	return site, nil
}

//Retrieve Site
func (c *Client) RetrieveSiteIdByName(name string) (int, error) {
	action := fmt.Sprintf("sites/?name=%s", name)
	req, err := c.NewRequest(action, "GET", "")
	if err != nil {
		return 0, fmt.Errorf("Error crafting request: %s", err)
	}
	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return 0, fmt.Errorf("Error creating site: %s", err)
	}
	sitesResp := new(SitesResponse)
	err = decodeBody(resp, &sitesResp)
	if err != nil {
		return 0, fmt.Errorf("Error parsing site response: %s", err)
	}

	siteId, err := sitesResp.GetId()

	if err != nil {
		return 0, fmt.Errorf("API error: %s", err)
	}
	return siteId, nil

}

//Destroy Site by Name
func (c *Client) DestroySitebyName(name string) error {
	siteId, err := c.RetrieveSiteIdByName(name)
	if err != nil {
		return fmt.Errorf("Error retrieving ID: %s", err)
	}
	action := fmt.Sprintf("sites/%d/", siteId)
	req, err := c.NewRequest(action, "DELETE", "")
	if err != nil {
		return fmt.Errorf("Error crafting request: %s", err)
	}
	_, err = checkResp(c.Http.Do(req))
	if err != nil {
		return fmt.Errorf("Error Deleting site: %s", err)
	}
	return nil
}

//Destroy Site by ID
func (c *Client) DestroySitebyID(id int) error {
	action := fmt.Sprintf("sites/%d/", id)
	req, err := c.NewRequest(action, "DELETE", "")
	if err != nil {
		return fmt.Errorf("Error crafting request: %s", err)
	}
	_, err = checkResp(c.Http.Do(req))
	if err != nil {
		return fmt.Errorf("Error Deleting site: %s", err)
	}
	return nil
}

//Retrieve Site by ID
func (c *Client) RetrieveSitebyID(id int) (*Site, error) {
	action := fmt.Sprintf("sites/%d/", id)
	req, err := c.NewRequest(action, "GET", "")

	if err != nil {
		return nil, fmt.Errorf("Error crafting request: %s", err)
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return nil, fmt.Errorf("Error Retrieving site: %s", err)
	}

	siteResp := new(SiteResponse)
	err = decodeBody(resp, &siteResp)

	if err != nil {
		return nil, fmt.Errorf("Error parsing site response: %s", err)
	}

	site, err := siteResp.GetSite()

	if err != nil {
		return nil, fmt.Errorf("API error: %s", err)
	}

	return site, nil
}

//Retrieve Site by name
func (c *Client) RetrieveSitebyName(name string) (*Site, error) {
	siteId, err := c.RetrieveSiteIdByName(name)
	action := fmt.Sprintf("sites/%d/", siteId)
	req, err := c.NewRequest(action, "GET", "")

	if err != nil {
		return nil, fmt.Errorf("Error crafting request: %s", err)
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return nil, fmt.Errorf("Error Retrieving site: %s", err)
	}

	siteResp := new(SiteResponse)
	err = decodeBody(resp, &siteResp)

	if err != nil {
		return nil, fmt.Errorf("Error parsing site response: %s", err)
	}

	site, err := siteResp.GetSite()

	if err != nil {
		return nil, fmt.Errorf("API error: %s", err)
	}

	return site, nil
}

//Update Site by ID
func (c *Client) UpdateSitebyID(id int, opts *SiteOpts) (*Site, error) {
	action := fmt.Sprintf("sites/%d/", id)

	//Create Map of data to send to API
	data := make(map[string]string)

	if opts.Name != "" {
		data["name"] = opts.Name
	}

	if opts.Desc != "" {
		data["description"] = opts.Desc
	}

	jsonData, err := json.Marshal(data)

	req, err := c.NewRequest(action, "PATCH", string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("Error crafting request: %s", err)
	}
	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, fmt.Errorf("Error updating site: %s", err)
	}
	siteResp := new(SiteResponse)
	err = decodeBody(resp, &siteResp)

	if err != nil {
		return nil, fmt.Errorf("Error parsing site response: %s", err)
	}
	site, err := siteResp.GetSite()
	if err != nil {
		return nil, fmt.Errorf("API error: %s", err)
	}
	return site, nil

}
