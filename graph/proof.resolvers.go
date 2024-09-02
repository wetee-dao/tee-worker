package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.42

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/wetee-dao/go-sdk/pallet/types"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"wetee.app/worker/graph/model"
	"wetee.app/worker/mint"
	"wetee.app/worker/mint/proof"
	"wetee.app/worker/util"

	wtypes "wetee.app/worker/type"
)

// WorkLoglist is the resolver for the work_loglist field.
func (r *queryResolver) WorkLoglist(ctx context.Context, workType string, workID int, page int, size int) (string, error) {
	wid := types.WorkId{
		Id:    uint64(workID),
		Wtype: util.GetWorkType(workType),
	}
	list, err := proof.ListLogsById(wid, page, size, false)
	if err != nil && err.Error() != "the list not found" {
		return "", gqlerror.Errorf("WorkLogList:" + err.Error())
	}

	listCache, err := proof.ListLogsById(wid, page, 200, true)
	if err != nil && err.Error() != "the list not found" {
		return "", gqlerror.Errorf("WorkLogCacheList:" + err.Error())
	}

	listCache = append(listCache, list...)
	bt, err := json.Marshal(listCache)
	if err != nil {
		return "", gqlerror.Errorf("JsonMarshal:" + err.Error())
	}

	return string(bt), nil
}

// WorkWetriclist is the resolver for the work_wetriclist field.
func (r *queryResolver) WorkWetriclist(ctx context.Context, workType string, workID int, page int, size int) (string, error) {
	wid := types.WorkId{
		Id:    uint64(workID),
		Wtype: util.GetWorkType(workType),
	}
	list, err := proof.ListMonitoringsById(wid, page, size, false)
	if err != nil && err.Error() != "the list not found" {
		return "", gqlerror.Errorf("WorkLogList:" + err.Error())
	}

	listCache, err := proof.ListMonitoringsById(wid, page, 200, true)
	if err != nil && err.Error() != "the list not found" {
		return "", gqlerror.Errorf("WorkLogList:" + err.Error())
	}

	listCache = append(listCache, list...)
	bt, err := json.Marshal(listCache)
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

// AttestationReportVerify is the resolver for the attestation_report_verify field.
func (r *queryResolver) AttestationReportVerify(ctx context.Context, report string) (bool, error) {
	reportx := strings.TrimPrefix(report, "0x")
	bt, err := hex.DecodeString(reportx)
	if err != nil {
		return false, gqlerror.Errorf("HexDecodeString:" + err.Error())
	}

	ps := wtypes.TeeParam{}
	json.Unmarshal(bt, &wtypes.TeeParam{})

	_, err = proof.VerifyReportProof(&ps)
	if err != nil {
		return false, gqlerror.Errorf("VerifyLocalReport error:" + err.Error())
	}

	return true, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
