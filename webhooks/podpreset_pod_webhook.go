/*
Copyright 2021 cuisongliu@qq.com.

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

package webhooks

import (
	"context"
	wk "github.com/cuisongliu/webhook"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// log is for logging in this package.
var podLog = logf.Log.WithName("pod-resource")

type PodPresetWebhook struct {
	object *v1.Pod
	client client.Client
}

func (a *PodPresetWebhook) OutRuntimeObject() runtime.Object {
	return a.object
}
func (a *PodPresetWebhook) GetClient() client.Client {
	return a.client
}
func (r *PodPresetWebhook) IntoRuntimeObject(object runtime.Object) {
	obj := &v1.Pod{}
	_ = wk.JsonConvert(object, obj)
	r.object = obj
}

func (a *PodPresetWebhook) InjectClient(c client.Client) error {
	a.client = c
	return nil
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

var _ wk.Defaulter = &PodPresetWebhook{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PodPresetWebhook) Default() {
	podLog.Info("default", "name", r.object.Name)
	label := r.object.Labels
	label["webhook"] = "webhook1"
	_ = r.client.Patch(context.Background(), r.object, client.RawPatch(types.MergePatchType, PatchData([]PatchStruct{
		{
			Op:    OpReplace,
			Path:  "/metadata/labels",
			Value: label,
		},
	})), client.DryRunAll)
}

// +kubebuilder:webhook:path=/mutate-podpreset-core-v1-pod,mutating=true,failurePolicy=ignore,groups=core,resources=pods,verbs=create;update,versions=v1,name=mpodpreset.pod.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
