package controller

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// GatewayClassReconciler reconciles a GatewayClass object
type GatewayClassReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses/status,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *GatewayClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the GatewayClass instance
	// The purpose is check if the Custom Resource for the Kind GatewayClass
	// is applied on the cluster if not we return nil to stop the reconciliation
	gatewayClass := &gatewayv1.GatewayClass{}
	if err := r.Get(ctx, req.NamespacedName, gatewayClass); err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("gatewayclass resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get gatewayclass")
		return ctrl.Result{}, err
	}

	// check if GatewayClass controllerName is ours
	if gatewayClass.Spec.ControllerName != controllerName {
		log.Info("Ignoring gatewayclass with non-matching controllerName")
		return ctrl.Result{}, nil
	}

	// validate parameters
	msg := ""
	_, api, err := InitCloudflareApi(ctx, r.Client, gatewayClass.Name)
	if err == nil {
		token, err := api.User.Tokens.Verify(ctx)
		if err == nil {
			if token.Status != "active" {
				msg = fmt.Sprintf("Token status is %s, is not active. Please check the Cloudflare dashboard", token.Status)
			}
		} else {
			msg = err.Error() + " Ensure ACCOUNT_ID and TOKEN are valid"
		}
	} else {
		msg = err.Error() + " Ensure ACCOUNT_ID and TOKEN are set"
	}

	var condition metav1.Condition
	if msg != "" {
		condition = metav1.Condition{
			Type:               string(gatewayv1.GatewayClassConditionStatusAccepted),
			Status:             metav1.ConditionFalse,
			Reason:             string(gatewayv1.GatewayClassReasonInvalidParameters),
			Message:            "Unable to initialize Cloudflare API. " + msg,
			ObservedGeneration: gatewayClass.Generation,
		}
	} else {
		condition = metav1.Condition{
			Type:               string(gatewayv1.GatewayClassConditionStatusAccepted),
			Status:             metav1.ConditionTrue,
			Reason:             string(gatewayv1.GatewayClassReasonAccepted),
			ObservedGeneration: gatewayClass.Generation,
		}
	}

	meta.SetStatusCondition(&gatewayClass.Status.Conditions, condition)
	if err := r.Status().Update(ctx, gatewayClass); err != nil {
		log.Error(err, "Failed to update GatewayClass status. Retrying in 1 minute")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.GenerationChangedPredicate{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&gatewayv1.GatewayClass{
			Spec: gatewayv1.GatewayClassSpec{
				ControllerName: "github.com/alodex/cloudflare-kubernetes-gateway",
			},
		}).
		WithEventFilter(pred).
		Complete(r)
}
