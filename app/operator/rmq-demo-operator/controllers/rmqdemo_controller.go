/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/edge-orchestration/app-fleet-demo/api/v1alpha1"
	edgeoperatorsv1alpha1 "github.com/edge-orchestration/app-fleet-demo/api/v1alpha1"
	"github.com/edge-orchestration/app-fleet-demo/assets"
)

// RmqDemoReconciler reconciles a RmqDemo object
type RmqDemoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=edge-operators.sw.bds.atos.net,resources=rmqdemoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=edge-operators.sw.bds.atos.net,resources=rmqdemoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=edge-operators.sw.bds.atos.net,resources=rmqdemoes/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RmqDemo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *RmqDemoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	logger := log.FromContext(ctx)
	// CR
	operatorCR := &v1alpha1.RmqDemo{}
	err := r.Get(ctx, req.NamespacedName, operatorCR)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("RmqDemo Operator resource object not found.")
		return ctrl.Result{}, nil
	} else if err != nil {
		logger.Error(err, "Error getting RmqDemo resource object")
		return ctrl.Result{}, err
	}
	// Publisher
	if err := r.ReconcilePublisher(ctx, req, operatorCR); err != nil {
		return ctrl.Result{}, err
	}
	// Consumer
	if err := r.ReconcileConsumer(ctx, req, operatorCR); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, err
}

func (r *RmqDemoReconciler) ReconcilePublisher(ctx context.Context, req ctrl.Request, operatorCR *v1alpha1.RmqDemo) error {
	_ = log.FromContext(ctx)
	logger := log.FromContext(ctx)

	// Deployments
	publisherDeploymentName := req.Name + "-publisher"
	namespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      publisherDeploymentName}
	publisherDeployment := &appsv1.Deployment{}
	create := false
	err := r.Get(ctx, namespacedName, publisherDeployment)
	if err != nil && errors.IsNotFound(err) {
		create = true
		publisherDeployment = assets.GetDeploymentFromFile("manifests/rmq-publisher-deployment.yaml")
	} else if err != nil {
		logger.Error(err, "Error getting existing Publisher Deployment.")
		return err
	}
	// update publisher deployment fields
	publisherDeployment.Namespace = req.Namespace
	publisherDeployment.Name = publisherDeploymentName
	if operatorCR.Spec.Publisher.Replicas != nil {
		publisherDeployment.Spec.Replicas = operatorCR.Spec.Publisher.Replicas
	}
	if operatorCR.Spec.Publisher.Port != nil {
		publisherDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *operatorCR.Spec.Publisher.Port
	}
	// ref
	ctrl.SetControllerReference(operatorCR, publisherDeployment, r.Scheme)
	// create or update deployment
	if create {
		err = r.Create(ctx, publisherDeployment)
	} else {
		err = r.Update(ctx, publisherDeployment)
	}
	if err != nil {
		return err
	}

	// Service
	publisherServiceName := req.Name + "-publisher"
	namespacedName = types.NamespacedName{
		Namespace: req.Namespace,
		Name:      publisherServiceName}
	publisherService := &corev1.Service{}
	create = false
	err = r.Get(ctx, namespacedName, publisherService)
	if err != nil && errors.IsNotFound(err) {
		create = true
		publisherService = assets.GetServiceFromFile("manifests/rmq-publisher-service.yaml")
	} else if err != nil {
		logger.Error(err, "Error getting existing Publisher Service.")
		return err
	}
	publisherService.Namespace = req.Namespace
	publisherService.Name = publisherServiceName
	// ref
	ctrl.SetControllerReference(operatorCR, publisherService, r.Scheme)
	// create or update service
	if create {
		err = r.Create(ctx, publisherService)
	} else {
		err = r.Update(ctx, publisherService)
	}
	return err
}

func (r *RmqDemoReconciler) ReconcileConsumer(ctx context.Context, req ctrl.Request, operatorCR *v1alpha1.RmqDemo) error {
	_ = log.FromContext(ctx)
	logger := log.FromContext(ctx)

	// Deployments
	consumerDeploymentName := req.Name + "-consumer"
	namespacedName := types.NamespacedName{
		Namespace: req.Namespace,
		Name:      consumerDeploymentName}
	consumerDeployment := &appsv1.Deployment{}
	create := false
	err := r.Get(ctx, namespacedName, consumerDeployment)
	if err != nil && errors.IsNotFound(err) {
		create = true
		consumerDeployment = assets.GetDeploymentFromFile("manifests/rmq-consumer-deployment.yaml")
	} else if err != nil {
		logger.Error(err, "Error getting existing Consumer Deployment.")
		return err
	}
	// update consumer deployment fields
	consumerDeployment.Namespace = req.Namespace
	consumerDeployment.Name = consumerDeploymentName
	if operatorCR.Spec.Consumer.Replicas != nil {
		consumerDeployment.Spec.Replicas = operatorCR.Spec.Consumer.Replicas
	}
	if operatorCR.Spec.Consumer.Port != nil {
		consumerDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = *operatorCR.Spec.Consumer.Port
	}
	// ref
	ctrl.SetControllerReference(operatorCR, consumerDeployment, r.Scheme)
	// create or update deployment
	if create {
		err = r.Create(ctx, consumerDeployment)
	} else {
		err = r.Update(ctx, consumerDeployment)
	}
	if err != nil {
		return err
	}

	// Service
	consumerServiceName := req.Name + "-consumer"
	namespacedName = types.NamespacedName{
		Namespace: req.Namespace,
		Name:      consumerServiceName}
	consumerService := &corev1.Service{}
	create = false
	err = r.Get(ctx, namespacedName, consumerService)
	if err != nil && errors.IsNotFound(err) {
		create = true
		consumerService = assets.GetServiceFromFile("manifests/rmq-consumer-service.yaml")
	} else if err != nil {
		logger.Error(err, "Error getting existing Consumer Service.")
		return err
	}
	consumerService.Namespace = req.Namespace
	consumerService.Name = consumerServiceName
	// ref
	ctrl.SetControllerReference(operatorCR, consumerService, r.Scheme)
	// create or update service
	if create {
		err = r.Create(ctx, consumerService)
	} else {
		err = r.Update(ctx, consumerService)
	}
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *RmqDemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&edgeoperatorsv1alpha1.RmqDemo{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		//Owns(&networkingv1.Ingress{}).
		Complete(r)
}
