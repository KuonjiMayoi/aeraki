// Copyright Aeraki Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aeraki-framework/aeraki/client-go/pkg/apis/metaprotocol/v1alpha1"
	"istio.io/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var metaRouterLog = log.RegisterScope("meta-router-controller", "meta-routerl-controller debugging", 0)

var (
	metaRouterlPredicates = predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			switch old := e.ObjectOld.(type) {
			case *v1alpha1.MetaRouter:
				new, ok := e.ObjectNew.(*v1alpha1.MetaRouter)
				if !ok {
					return false
				}
				if !reflect.DeepEqual(old.Spec, new.Spec) ||
					old.GetDeletionTimestamp() != new.GetDeletionTimestamp() ||
					old.GetGeneration() != new.GetGeneration() {
					return true
				}
			default:
				return false
			}
			return false
		},
	}
)

// MetaRouterController control ApplicationProtocol
type MetaRouterController struct {
	client.Client
	updateRouteCache func() error
}

// Reconcile will try to trigger once mcp push.
func (r *MetaRouterController) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	metaRouterLog.Infof("reconcile: %s/%s", request.Namespace, request.Name)
	if r.updateRouteCache != nil {
		err := r.updateRouteCache()
		if err != nil {
			return reconcile.Result{Requeue: true}, err
		}
	}
	return reconcile.Result{}, nil
}

// AddMetaRouterController adds MetaRouterController
func AddMetaRouterController(mgr manager.Manager, triggerPush func() error) error {
	metaProtocolCtrl := &MetaRouterController{Client: mgr.GetClient(), updateRouteCache: triggerPush}
	c, err := controller.New("aeraki-meta-protocol-meta-router-controller", mgr,
		controller.Options{Reconciler: metaProtocolCtrl})
	if err != nil {
		return err
	}
	// Watch for changes on MetaRouter CRD
	err = c.Watch(&source.Kind{Type: &v1alpha1.MetaRouter{}}, &handler.EnqueueRequestForObject{},
		metaRouterlPredicates)
	if err != nil {
		return err
	}
	controllerLog.Infof("MetaProtocolMetaRouterController registered")
	return nil
}
