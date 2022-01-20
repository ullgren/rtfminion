package anypointclient

import (
	"fmt"
	"net/http"
)

type Fabric struct {
	ID                        string    `json:"id"`
	Name                      string    `json:"name"`
	Region                    string    `json:"region"`
	Vendor                    string    `json:"vendor"`
	OrganizationID            string    `json:"organizationId"`
	Version                   string    `json:"version"`
	Status                    string    `json:"status"`
	LastUpgradeTimestamp      int64     `json:"lastUpgradeTimestamp"`
	CreatedAt                 int64     `json:"createdAt"`
	Nodes                     []Nodes   `json:"nodes"`
	SecondsSinceHeartbeat     int       `json:"secondsSinceHeartbeat"`
	ClusterVersion            string    `json:"clusterVersion"`
	KubernetesVersion         string    `json:"kubernetesVersion"`
	IsManaged                 bool      `json:"isManaged"`
	Appliance                 Appliance `json:"appliance"`
	ClusterConfigurationLevel string    `json:"clusterConfigurationLevel"`
	Features                  Features  `json:"features"`
}
type Status struct {
	IsHealthy     bool `json:"isHealthy"`
	IsReady       bool `json:"isReady"`
	IsSchedulable bool `json:"isSchedulable"`
}
type Capacity struct {
	CPU       int    `json:"cpu"`
	CPUMillis int    `json:"cpuMillis"`
	Memory    string `json:"memory"`
	MemoryMi  int    `json:"memoryMi"`
	Pods      int    `json:"pods"`
}
type AllocatedRequestCapacity struct {
	CPU       int    `json:"cpu"`
	CPUMillis int    `json:"cpuMillis"`
	Memory    string `json:"memory"`
	MemoryMi  int    `json:"memoryMi"`
	Pods      int    `json:"pods"`
}
type AllocatedLimitCapacity struct {
	CPU       int    `json:"cpu"`
	CPUMillis int    `json:"cpuMillis"`
	Memory    string `json:"memory"`
	MemoryMi  int    `json:"memoryMi"`
	Pods      int    `json:"pods"`
}
type Nodes struct {
	UID                      string                   `json:"uid"`
	Name                     string                   `json:"name"`
	KubeletVersion           string                   `json:"kubeletVersion"`
	DockerVersion            string                   `json:"dockerVersion"`
	Role                     string                   `json:"role"`
	Status                   Status                   `json:"status"`
	Capacity                 Capacity                 `json:"capacity"`
	AllocatedRequestCapacity AllocatedRequestCapacity `json:"allocatedRequestCapacity"`
	AllocatedLimitCapacity   AllocatedLimitCapacity   `json:"allocatedLimitCapacity"`
}
type Appliance struct {
	Version string `json:"version"`
}
type Features struct {
	EnhancedSecurity bool `json:"enhancedSecurity"`
}

type FabricsListResponse struct {
	Items []Fabric `json:"items"`
}

func (client *AnypointClient) ListFabrics(organization Organization) (*FabricsListResponse, error) {

	payload := new(FabricsListResponse)

	req, _ := client.newRequest("GET", fmt.Sprintf("/runtimefabric/api/organizations/%s/fabrics", organization.ID), nil)
	res, err := client.HTTPClient.Do(req)
	err = checkResponse(res, err)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		err = readResponseAs(res, payload)
		if err != nil {
			return nil, err
		}
	}
	return payload, nil
}

func (c *AnypointClient) ResolveFabricByName(organization Organization, fabricName string) (Fabric, error) {
	fabrics, err := c.ListFabrics(organization)
	if err != nil {
		return Fabric{}, err
	}
	for _, fabric := range fabrics.Items {
		if fabric.Name == fabricName {
			return fabric, nil
		}
	}
	return Fabric{}, fmt.Errorf("failed to resolve fabric named %s in organtisation %s", fabricName, organization.Name)
}
