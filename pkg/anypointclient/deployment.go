package anypointclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type Deployment struct {
	ID                    string      `json:"id"`
	Name                  string      `json:"name"`
	CreationDate          int64       `json:"creationDate"`
	LastModifiedDate      int64       `json:"lastModifiedDate"`
	Target                Target      `json:"target"`
	Status                string      `json:"status"`
	Application           Application `json:"application"`
	DesiredVersion        string      `json:"desiredVersion"`
	Replicas              []Replicas  `json:"replicas"`
	LastSuccessfulVersion interface{} `json:"lastSuccessfulVersion"`
}
type Jvm struct {
	Args string `json:"args"`
}
type Inbound struct {
	PublicURL string `json:"publicUrl"`
}
type HTTP struct {
	Inbound Inbound `json:"inbound"`
}
type CPU struct {
	Limit    string `json:"limit"`
	Reserved string `json:"reserved"`
}
type Memory struct {
	Limit    string `json:"limit"`
	Reserved string `json:"reserved"`
}
type Resources struct {
	CPU    CPU    `json:"cpu"`
	Memory Memory `json:"memory"`
}
type AnypointMonitoring struct {
	Image     string    `json:"image"`
	Resources Resources `json:"resources"`
}
type Sidecars struct {
	AnypointMonitoring AnypointMonitoring `json:"anypoint-monitoring"`
}
type DeploymentSettings struct {
	Jvm                                 Jvm       `json:"jvm"`
	HTTP                                HTTP      `json:"http"`
	Sidecars                            Sidecars  `json:"sidecars"`
	Clustered                           bool      `json:"clustered"`
	Resources                           Resources `json:"resources"`
	RuntimeVersion                      string    `json:"runtimeVersion"`
	UpdateStrategy                      string    `json:"updateStrategy"`
	LastMileSecurity                    bool      `json:"lastMileSecurity"`
	ForwardSslSession                   bool      `json:"forwardSslSession"`
	EnforceDeployingReplicasAcrossNodes bool      `json:"enforceDeployingReplicasAcrossNodes"`
}
type Target struct {
	Provider           string             `json:"provider"`
	TargetID           string             `json:"targetId"`
	DeploymentSettings DeploymentSettings `json:"deploymentSettings"`
	Replicas           int                `json:"replicas"`
	Type               string             `json:"Type"`
}
type Ref struct {
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Version    string `json:"version"`
	Packaging  string `json:"packaging"`
}

type MuleAgentApplicationPropertiesService struct {
	Properties      map[string]interface{} `json:"properties"`
	ApplicationName string                 `json:"applicationName"`
}
type Configuration struct {
	MuleAgentApplicationPropertiesService MuleAgentApplicationPropertiesService `json:"mule.agent.application.properties.service"`
}
type Application struct {
	Status        string        `json:"status"`
	DesiredState  string        `json:"desiredState"`
	Ref           Ref           `json:"ref"`
	Configuration Configuration `json:"configuration"`
}
type Replicas struct {
	State                    string `json:"state"`
	DeploymentLocation       string `json:"deploymentLocation"`
	CurrentDeploymentVersion string `json:"currentDeploymentVersion"`
	Reason                   string `json:"reason"`
}

type DeploymentsListResponse struct {
	Total int64        `json:"total"`
	Items []Deployment `json:"items"`
}

func (client *AnypointClient) ListDeployments(organization Organization, environment Environment, fabric Fabric) (*DeploymentsListResponse, *http.Response, error) {
	listDeploymentsUrl := fmt.Sprintf("/hybrid/api/v2/organizations/%s/environments/%s/deployments", organization.ID, environment.ID)
	req, _ := client.newRequest("GET", listDeploymentsUrl, nil)

	q := req.URL.Query()
	// Make sure we only get deployments with the provider MC (which indirectly indicates RTF)
	q.Add("provider", "MC")
	if fabric.ID != "" {
		q.Add("target", fabric.ID)
	}
	req.URL.RawQuery = q.Encode()

	res, err := client.HTTPClient.Do(req)
	err = checkResponse(res, err)
	if err != nil {
		return nil, res, errors.Wrapf(err, "failed to list deployments for %s in %s", environment.Name, organization.Name)
	}
	if res.StatusCode != 200 {
		return nil, res, fmt.Errorf("listing deployments failed with status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	payload := new(DeploymentsListResponse)
	err = readResponseAs(res, payload)
	if err != nil {
		return nil, res, errors.Wrapf(err, "failed to decode response")
	}
	/*
		for i, d := range payload.Items {
			dd, err := client.GetDeploymentDetails(organization, environment, d.ID)
			if err == nil {
				payload.Items[i] = dd
			}
		}
		log.Printf("Number of deployments %d", len(payload.Items))
	*/

	// Create two channels one for the jobs and one for the results
	jobs := make(chan Deployment, 10)
	results := make(chan Deployment, len(payload.Items))
	var wg sync.WaitGroup

	// Create 4 worker, more than that seems to give 503 errors
	for j := 0; j < 4; j++ {
		wg.Add(1)
		go deploymentListWorker(j, client, organization, environment, jobs, results, &wg)
	}

	// Send jobs to the workers and close it once all
	for _, d := range payload.Items {
		jobs <- d
	}
	close(jobs)
	// Close the results once all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	payload.Items = nil
	for d := range results {
		payload.Items = append(payload.Items, d)
	}
	log.Printf("Number of deployments %d", len(payload.Items))
	return payload, res, nil
}

func deploymentListWorker(id int, client *AnypointClient, organization Organization, environment Environment, deployments <-chan Deployment, results chan<- Deployment, wg *sync.WaitGroup) {
	defer wg.Done()
	for d := range deployments {
		dd, err := client.GetDeploymentDetails(organization, environment, d.ID)
		if err == nil {
			results <- dd
		} else {
			log.Printf("Failed to fetch details for %s : %+v", d.Name, err)
			results <- d
		}
	}
}

func (client *AnypointClient) GetDeploymentDetails(organization Organization, environment Environment, deploymentID string) (Deployment, error) {
	deploymentDetailsUrl := fmt.Sprintf("/hybrid/api/v2/organizations/%s/environments/%s/deployments/%s", organization.ID, environment.ID, deploymentID)
	req, _ := client.newRequest("GET", deploymentDetailsUrl, nil)

	q := req.URL.Query()
	// Make sure we only get deployments with the provider MC (which indirectly indicates RTF)
	q.Add("provider", "MC")
	req.URL.RawQuery = q.Encode()

	res, err := client.HTTPClient.Do(req)
	err = checkResponse(res, err)
	if err != nil {
		return Deployment{}, errors.Wrapf(err, "failed to list deployments for %s in %s", environment.Name, organization.Name)
	}
	if res.StatusCode != 200 {
		return Deployment{}, fmt.Errorf("listing deployments failed with status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	payload := new(Deployment)
	err = readResponseAs(res, payload)
	if err != nil {
		return Deployment{}, err
	}
	return *payload, nil
}

func (client *AnypointClient) UpdateDeployment(organization Organization, environment Environment, deployment Deployment) error {
	deploymentDetailsUrl := fmt.Sprintf("/hybrid/api/v2/organizations/%s/environments/%s/deployments/%s", organization.ID, environment.ID, deployment.ID)
	bytesToSend := new(bytes.Buffer)
	json.NewEncoder(bytesToSend).Encode(deployment)
	req, _ := client.newRequest("PATCH", deploymentDetailsUrl, bytesToSend)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	res, err := client.HTTPClient.Do(req)
	err = checkResponse(res, err)
	if err != nil {
		return errors.Wrapf(err, "failed to update deployment for %s in %s", environment.Name, organization.Name)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to update deployment with status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	return nil
}

func (client *AnypointClient) StopDeployment(organization Organization, environment Environment, deployment Deployment) error {
	deploymentDetailsUrl := fmt.Sprintf("/hybrid/api/v2/organizations/%s/environments/%s/deployments/%s", organization.ID, environment.ID, deployment.ID)
	dataToSend := strings.NewReader(`{"application":{"desiredState":"STOPPED"}}`)
	req, _ := client.newRequest("PATCH", deploymentDetailsUrl, dataToSend)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	res, err := client.HTTPClient.Do(req)
	err = checkResponse(res, err)
	if err != nil {
		return errors.Wrapf(err, "failed to stop deployment %s in environment %s in %s", deployment.Name, environment.Name, organization.Name)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to stop deployment with status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	return nil
}

func (client *AnypointClient) StartDeployment(organization Organization, environment Environment, deployment Deployment) error {
	deploymentDetailsUrl := fmt.Sprintf("/hybrid/api/v2/organizations/%s/environments/%s/deployments/%s", organization.ID, environment.ID, deployment.ID)
	dataToSend := strings.NewReader(`{"application":{"desiredState":"STARTED"}}`)
	req, _ := client.newRequest("PATCH", deploymentDetailsUrl, dataToSend)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	res, err := client.HTTPClient.Do(req)
	err = checkResponse(res, err)
	if err != nil {
		return errors.Wrapf(err, "failed to start deployment %s in environment %s in %s", deployment.Name, environment.Name, organization.Name)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to start deployment with status code %d", res.StatusCode)
	}
	defer res.Body.Close()

	return nil
}
