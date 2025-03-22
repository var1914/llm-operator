package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	llmv1alpha1 "github.com/var1914/llm-operator/api/v1alpha1"
)

// LLMModelReconciler reconciles a LLMModel object
type LLMModelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=llm.example.com,resources=llmmodels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=llm.example.com,resources=llmmodels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=llm.example.com,resources=llmmodels/finalizers,verbs=update

// Reconcile handles LLMModel resources
func (r *LLMModelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the LLMModel instance
	model := &llmv1alpha1.LLMModel{}
	if err := r.Get(ctx, req.NamespacedName, model); err != nil {
		// Handle deletion
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// For now, just set the status to Ready
	if model.Status.Phase != "Ready" {
		model.Status.Phase = "Ready"
		model.Status.Message = "Model is ready to be deployed"
		if err := r.Status().Update(ctx, model); err != nil {
			logger.Error(err, "Failed to update model status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LLMModelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&llmv1alpha1.LLMModel{}).
		Complete(r)
}
