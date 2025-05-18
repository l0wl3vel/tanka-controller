/*
Copyright 2024 The Flux authors

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

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	runtimeClient "github.com/fluxcd/pkg/runtime/client"
	helper "github.com/fluxcd/pkg/runtime/controller"
	tankav2alpha1 "github.com/fluxcd/tanka-controller/api/v2alpha1"
	kuberecorder "k8s.io/client-go/tools/record"
)

// TankaReleaseReconciler reconciles a TankaRelease object
type TankaReleaseReconciler struct {
	client.Client
	kuberecorder.EventRecorder
	helper.Metrics

	GetClusterConfig func() (*rest.Config, error)
	ClientOpts       runtimeClient.Options
	KubeConfigOpts   runtimeClient.KubeConfigOptions
	APIReader        client.Reader

	FieldManager               string
	DefaultServiceAccount      string

	requeueDependency    time.Duration
	artifactFetchRetries int
}

type TankaReleaseReconcilerOptions struct {
	DependencyRequeueInterval time.Duration
	HTTPRetry                 int
}

// +kubebuilder:rbac:groups=tanka.toolkit.fluxcd.io,resources=tankareleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tanka.toolkit.fluxcd.io,resources=tankareleases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=ocirepositories,verbs=get;list;watch
// +kubebuilder:rbac:groups=source.toolkit.fluxcd.io,resources=ocirepositories/status,verbs=get
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *TankaReleaseReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()


	// your logic here

	return ctrl.Result{}, nil
}

func (r *TankaReleaseReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, opts TankaReleaseReconcilerOptions) error {

	if err := mgr.GetFieldIndexer().IndexField(ctx, &tankav2alpha1.TankaRelease{}, tankav2alpha1.SourceIndexKey,
		func(o client.Object) []string {
			obj := o.(*tankav2alpha1.TankaRelease)
			namespacedName, err := getNamespacedName(obj)
			if err != nil {
				return nil
			}
			return []string{
				namespacedName.String(),
			}
		},
	); err != nil {
		return err
	}

	r.requeueDependency = opts.DependencyRequeueInterval
	r.artifactFetchRetries = opts.HTTPRetry

	return ctrl.NewControllerManagedBy(mgr).
		For(&tankav2alpha1.TankaRelease{}).
		Watches(
			&sourcev1beta2.OCIRepository{},
			handler.EnqueueRequestsFromMapFunc(r.requestsForOCIRrepositoryChange),
			builder.WithPredicates(intpredicates.SourceRevisionChangePredicate{}),
		).
		WithOptions(controller.Options{
			RateLimiter: opts.RateLimiter,
		}).
		Complete(r)
}

func getNamespacedName(obj *tankav2alpha1.TankaRelease) (types.NamespacedName, error) {
	namespacedName := types.NamespacedName{}
	
	namespacedName.Namespace = obj.Spec.TankaRef.Namespace
	namespacedName.Name = obj.Spec.TankaRef.Name

	return namespacedName, nil
}