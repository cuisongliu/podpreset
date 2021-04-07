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
	wk "github.com/cuisongliu/webhook"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func SetupWebhookWithManager(mgr ctrl.Manager) error {
	hookServer := mgr.GetWebhookServer()
	c := mgr.GetClient()
	wkpods := &wk.WebhookObject{
		WK:             hookServer,
		Webhook:        &PodPresetWebhook{},
		Obj:            &v1.Pod{},
		DefaultingPath: "/mutate-apps-v1-sts",
		Client:         c,
	}
	wkpods.Init()
	return nil
}

func SetupWebhookToSelector(obj map[string]*metav1.LabelSelector, namespace map[string]*metav1.LabelSelector) {
	//time zone
}
