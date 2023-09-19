package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	secretv1 "wetee.app/worker/api/v1"
)

// OracleReconciler reconciles a Oracle object
type OracleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secret.wetee.app,resources=oracles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secret.wetee.app,resources=oracles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=secret.wetee.app,resources=oracles/finalizers,verbs=update

func (r *OracleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	fmt.Println("Reconciling Oracle")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OracleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretv1.Oracle{}).
		Complete(r)
}
