/*
Copyright The L5 Operator Authors

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

package sub_reconcilers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	petsv1 "github.com/opdev/l5-operator-demo/api/v1"
	"github.com/opdev/l5-operator-demo/internal/util"
)

type DeploymentReconciler struct {
	client client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func NewDeploymentReconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme) *DeploymentReconciler {
	return &DeploymentReconciler{
		client: client,
		Log:    log,
		Scheme: scheme,
	}
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, bestie *petsv1.Bestie) (ctrl.Result, error) {
	logInfo := types.NamespacedName{
		Namespace: bestie.Namespace,
		Name:      "DeploymentReconciler",
	}
	log := r.Log.WithValues("deployment", logInfo)

	// TODO wait for db to be seeded
	// reconcile Deployment.
	log.Info("reconcile deployment if it does not exist")
	dp := &appsv1.Deployment{}
	err := r.client.Get(ctx, types.NamespacedName{Name: bestie.Name + "-app", Namespace: bestie.Namespace}, dp)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new app for bestie")
		fileName := "config/resources/bestie-deploy.yaml"
		err := r.applyManifests(ctx, bestie, dp, fileName)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("Error during Manifests apply - %w", err)
		}
	}
	return ctrl.Result{}, err
}

func (r *DeploymentReconciler) applyManifests(ctx context.Context, bestie *petsv1.Bestie, obj client.Object, fileName string) error {
	err := util.ApplyManifests(r.client, r.Scheme, ctx, bestie, obj, fileName)
	if err != nil {
		return err
	}
	return err
}
