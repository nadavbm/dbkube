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

package services

import (
	"context"
	"time"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	servicesv1alpha1 "github.com/nadavbm/dbkube/apis/services/v1alpha1"
	"github.com/nadavbm/dbkube/pkg/kubetz"
	"github.com/nadavbm/zlog"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	Logger *zlog.Logger
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=services.etz.com,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=services.etz.com,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=services.etz.com,resources=services/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := zlog.New()
	r.Logger = logger

	r.Logger.Info("Start reconcile", zap.String("namespace", req.NamespacedName.Namespace))

	var resource servicesv1alpha1.Service
	if err := r.Client.Get(context.Background(), req.NamespacedName, &resource); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("service resource is not found. skipping..", zap.String("namespace", req.Namespace))
			return ctrl.Result{Requeue: false, RequeueAfter: 0}, nil
		}
		r.Logger.Error("could not fetch resource", zap.String("type", resource.Kind))
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, err
	}

	var object v1.Service
	obj := kubetz.BuildService(req.Namespace, &resource)
	if err := r.Get(ctx, req.NamespacedName, &object); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("create service", zap.String("namespace", req.Namespace), zap.Error(err))
			if err := r.Create(ctx, obj); err != nil {
				r.Logger.Error("could not create object", zap.String("object kind", obj.Kind), zap.Error(err))
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		if err := r.Update(ctx, &object); err != nil {
			if errors.IsInvalid(err) {
				r.Logger.Error("invalid update", zap.String("object", object.Name), zap.Error(err))
			} else {
				r.Logger.Error("unable to update", zap.String("object", object.Name), zap.Error(err))
			}

			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servicesv1alpha1.Service{}).
		Owns(&v1.Service{}).
		Complete(r)
}
