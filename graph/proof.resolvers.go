package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/wetee-dao/go-sdk/gen/types"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/graph/model"
	"wetee.app/worker/mint"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/util"
)

// WorkLoglist is the resolver for the work_loglist field.
func (r *queryResolver) WorkLoglist(ctx context.Context, workType string, workID int, page int, size int) (string, error) {
	list, err := proof.ListLogsById(types.WorkId{
		Id:    uint64(workID),
		Wtype: util.GetWorkType(workType),
	}, page, size)
	if err != nil {
		return "", gqlerror.Errorf("WorkLogList:" + err.Error())
	}

	bt, err := json.Marshal(list)
	if err != nil {
		return "", gqlerror.Errorf("JsonMarshal:" + err.Error())
	}

	return string(bt), nil
}

// WorkWetriclist is the resolver for the work_wetriclist field.
func (r *queryResolver) WorkWetriclist(ctx context.Context, workType string, workID int, page int, size int) (string, error) {
	list, err := proof.ListMonitoringsById(types.WorkId{
		Id:    uint64(workID),
		Wtype: util.GetWorkType(workType),
	}, page, size)

	if err != nil {
		return "", gqlerror.Errorf("WorkLogList:" + err.Error())
	}

	bt, err := json.Marshal(list)
	if err != nil {
		return "", gqlerror.Errorf("JsonMarshal:" + err.Error())
	}

	return string(bt), nil
}

// WorkServicelist is the resolver for the work_servicelist field.
func (r *queryResolver) WorkServicelist(ctx context.Context, projectID string, workType string, workID int) ([]*model.Service, error) {
	if mint.MinterIns.ChainClient == nil {
		return nil, gqlerror.Errorf("Invalid chain client")
	}

	wid := types.WorkId{
		Id:    uint64(workID),
		Wtype: util.GetWorkType(workType),
	}
	name := util.GetWorkTypeStr(wid) + "-" + fmt.Sprint(wid.Id)

	client := mint.MinterIns.K8sClient
	ServiceSpace := client.CoreV1().Services(mint.HexStringToSpace(projectID))
	list, err := ServiceSpace.List(ctx, v1.ListOptions{
		LabelSelector: "service=" + name,
	})
	if err != nil {
		return nil, gqlerror.Errorf("WorkServiceList:" + err.Error())
	}

	var services []*model.Service = make([]*model.Service, 0, 10)
	if list != nil {
		for i := 0; i < len(list.Items); i++ {
			item := list.Items[i]
			var ports = make([]*model.ServicePort, 0, len(item.Spec.Ports))
			for j := 0; j < len(item.Spec.Ports); j++ {
				ports = append(ports, &model.ServicePort{
					Name:     item.ObjectMeta.Name,
					Port:     int(item.Spec.Ports[j].Port),
					Protocol: fmt.Sprint(item.Spec.Ports[j].Protocol),
					NodePort: int(item.Spec.Ports[j].NodePort),
				})
			}
			services = append(services, &model.Service{
				Type:  fmt.Sprint(item.Spec.Type),
				Ports: ports,
			})
		}
	}

	return services, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }