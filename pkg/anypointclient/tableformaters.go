package anypointclient

import (
	"fmt"

	"github.com/Redpill-Linpro/rtfminion/pkg/rtfprinter"
)

func (resp *FabricsListResponse) GetTableData() rtfprinter.TableData {

	var data rtfprinter.TableData

	data.Headers = []string{"Name", "ID", "Cluster Version", "Version"}

	for _, f := range resp.Items {
		var row []string
		row = append(row, f.Name)
		row = append(row, f.ID)
		row = append(row, f.ClusterVersion)
		row = append(row, f.Version)
		data.Rows = append(data.Rows, row)
	}

	return data
}

func (resp *DeploymentsListResponse) GetTableData() rtfprinter.TableData {
	var data rtfprinter.TableData

	data.Headers = []string{"Name", "Status", "Runtime version", "CPU Reserved", "CPU Limit", "Memory Reserved", "Memory Max", "Replicas"}

	for _, f := range resp.Items {
		var row []string
		row = append(row, f.Name)
		status := fmt.Sprintf("%s / %s (%s)",
			f.Application.Status, f.Application.DesiredState, f.Status)
		row = append(row, status)
		row = append(row, f.Target.DeploymentSettings.RuntimeVersion)
		row = append(row, f.Target.DeploymentSettings.Resources.CPU.Reserved)
		row = append(row, f.Target.DeploymentSettings.Resources.CPU.Limit)
		row = append(row, f.Target.DeploymentSettings.Resources.Memory.Reserved)
		row = append(row, f.Target.DeploymentSettings.Resources.Memory.Limit)
		row = append(row, fmt.Sprintf("%d", f.Target.Replicas))
		data.Rows = append(data.Rows, row)
	}

	return data
}
