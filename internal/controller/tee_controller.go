package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	secretv1 "wetee.app/worker/api/v1"
)

// TeeReconciler reconciles a Tee object
type TeeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=secret.wetee.app,resources=tees,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=secret.wetee.app,resources=tees/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=secret.wetee.app,resources=tees/finalizers,verbs=update
func (r *TeeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TeeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretv1.Tee{}).
		Complete(r)
}
