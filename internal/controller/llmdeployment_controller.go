package controllers

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	llmv1alpha1 "github.com/var1914/llm-operator/api/v1alpha1"
)

// LLMDeploymentReconciler reconciles a LLMDeployment object
type LLMDeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=llm.example.com,resources=llmdeployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=llm.example.com,resources=llmdeployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=llm.example.com,resources=llmdeployments/finalizers,verbs=update
//+kubebuilder:rbac:groups=llm.example.com,resources=llmmodels,verbs=get;list;watch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile handles LLMDeployment resources
func (r *LLMDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling LLMDeployment", "request", req.NamespacedName)

	// Fetch the LLMDeployment instance
	deployment := &llmv1alpha1.LLMDeployment{}
	if err := r.Get(ctx, req.NamespacedName, deployment); err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return
			logger.Info("LLMDeployment resource not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object
		logger.Error(err, "Failed to get LLMDeployment")
		return ctrl.Result{}, err
	}

	// Fetch the referenced model
	modelName := types.NamespacedName{
		Name:      deployment.Spec.ModelRef,
		Namespace: req.Namespace,
	}
	model := &llmv1alpha1.LLMModel{}
	if err := r.Get(ctx, modelName, model); err != nil {
		logger.Error(err, "Unable to fetch referenced LLMModel")

		// Update status to reflect missing model
		newCondition := llmv1alpha1.DeploymentCondition{
			Type:               "ModelNotFound",
			Status:             "False",
			LastTransitionTime: metav1.Now(),
			Reason:             "ModelNotFound",
			Message:            fmt.Sprintf("Referenced model %s not found", deployment.Spec.ModelRef),
		}

		// Update the status
		deployment.Status.Conditions = []llmv1alpha1.DeploymentCondition{newCondition}
		if err := r.Status().Update(ctx, deployment); err != nil {
			logger.Error(err, "Failed to update deployment status")
		}

		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Define the Kubernetes Deployment
	k8sDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		},
	}

	// Apply the Deployment
	op, err := controllerutil.CreateOrUpdate(ctx, r.Client, k8sDeployment, func() error {
		// Set controller reference
		if err := controllerutil.SetControllerReference(deployment, k8sDeployment, r.Scheme); err != nil {
			return err
		}

		// Configure the deployment spec
		k8sDeployment.Spec.Replicas = pointer.Int32(int32(deployment.Spec.Replicas))
		k8sDeployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: map[string]string{"app": deployment.Name},
		}

		// Configure pod template
		if k8sDeployment.Spec.Template.ObjectMeta.Labels == nil {
			k8sDeployment.Spec.Template.ObjectMeta.Labels = map[string]string{}
		}
		k8sDeployment.Spec.Template.ObjectMeta.Labels["app"] = deployment.Name

		// Configure containers
		k8sDeployment.Spec.Template.Spec.Containers = []corev1.Container{
			{
				Name:  "llm-model",
				Image: model.Spec.Image,
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: int32(deployment.Spec.Port),
						Name:          "http",
					},
				},
			},
		}

		// Add resource requirements if specified
		if model.Spec.Resources.CPU != "" || model.Spec.Resources.Memory != "" {
			resourceRequirements := corev1.ResourceRequirements{
				Requests: corev1.ResourceList{},
				Limits:   corev1.ResourceList{},
			}

			if model.Spec.Resources.CPU != "" {
				cpu := resource.MustParse(model.Spec.Resources.CPU)
				resourceRequirements.Requests[corev1.ResourceCPU] = cpu
				resourceRequirements.Limits[corev1.ResourceCPU] = cpu
			}

			if model.Spec.Resources.Memory != "" {
				memory := resource.MustParse(model.Spec.Resources.Memory)
				resourceRequirements.Requests[corev1.ResourceMemory] = memory
				resourceRequirements.Limits[corev1.ResourceMemory] = memory
			}

			k8sDeployment.Spec.Template.Spec.Containers[0].Resources = resourceRequirements
		}

		return nil
	})

	if err != nil {
		logger.Error(err, "Failed to create or update Deployment")
		return ctrl.Result{}, err
	}

	logger.Info("Deployment reconciliation", "operation", op)

	// Define the Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
		},
	}

	// Apply the Service
	op, err = controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		// Set controller reference
		if err := controllerutil.SetControllerReference(deployment, service, r.Scheme); err != nil {
			return err
		}

		// Configure service spec
		service.Spec.Selector = map[string]string{"app": deployment.Name}
		service.Spec.Ports = []corev1.ServicePort{
			{
				Port:       int32(deployment.Spec.Port),
				TargetPort: intstr.FromInt(deployment.Spec.Port),
				Name:       "http",
			},
		}

		return nil
	})

	if err != nil {
		logger.Error(err, "Failed to create or update Service")
		return ctrl.Result{}, err
	}

	logger.Info("Service reconciliation", "operation", op)

	// Update status
	deployment.Status.AvailableReplicas = k8sDeployment.Status.AvailableReplicas
	deployment.Status.Conditions = []llmv1alpha1.DeploymentCondition{
		{
			Type:               "Available",
			Status:             "True",
			LastTransitionTime: metav1.Now(),
			Reason:             "DeploymentAvailable",
			Message:            "Deployment is available",
		},
	}

	if err := r.Status().Update(ctx, deployment); err != nil {
		logger.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LLMDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&llmv1alpha1.LLMDeployment{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
