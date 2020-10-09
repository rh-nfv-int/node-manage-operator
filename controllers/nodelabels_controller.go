/*


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
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	nodemanagev1 "github.com/rh-nfv-int/node-manage-operator/api/v1"
)

// NodeLabelsReconciler reconciles a NodeLabels object
type NodeLabelsReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=nodemanage.openshift.io,resources=nodelabels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=nodemanage.openshift.io,resources=nodelabels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;create;update;patch

func (r *NodeLabelsReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("nodelabels", req.NamespacedName)

	nodeLabels := &nodemanagev1.NodeLabels{}
	err := r.Get(ctx, req.NamespacedName, nodeLabels)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("NodeLabels resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get NodeLabels")
		return ctrl.Result{}, err
	}

	nodeList := &corev1.NodeList{}
	listOpts := []client.ListOption{
		client.MatchingLabels(nodeLabels.Spec.NodeSelectorLabels),
	}
	err = r.List(ctx, nodeList, listOpts...)
	if err != nil {
		log.Error(err, "failed to get nodes list")
		return ctrl.Result{}, err
	}

	// Check if there are enough nodes for applying label group
	count := 0
	for _, group := range nodeLabels.Spec.LabelGroup {
		for k, v := range group.Labels {
			log.Info("Requested labels", k, v)
		}
		count += group.Count
	}

	log.Info("node matching", "available nodes", len(nodeList.Items), "requested count", count)
	if len(nodeList.Items) < count {
		err = fmt.Errorf("Not enough nodes to apply label group")
		log.Error(err, "Node request for applying labels and available nodes mismatch")
		return ctrl.Result{}, err
	}

	nodeIdx := 0
	for _, group := range nodeLabels.Spec.LabelGroup {
		for idx := 0; idx < group.Count; idx++ {
			node := nodeList.Items[nodeIdx]
			nodeIdx++
			log.Info("Patching node with labels", "name", node.Name)
			patch := client.MergeFrom(node.DeepCopy())
			for k, v := range group.Labels {
				node.Labels[k] = v
			}
			r.Patch(ctx, &node, patch)
		}
	}

	return ctrl.Result{}, nil
}

func (r *NodeLabelsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&nodemanagev1.NodeLabels{}).
		Complete(r)
}
