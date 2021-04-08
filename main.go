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

package main

import (
	"flag"
	"github.com/cuisongliu/podpreset/webhooks"
	wk "github.com/cuisongliu/webhook"
	level "go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var healthAddr, logLevel string
	var webhookPort int
	flag.IntVar(&webhookPort, "webhook-port", 9443, "The port of the webhook.")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&healthAddr, "health-addr", ":9090", "Health address. Readiness url is  /readyz, Liveness url is /healthz")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&logLevel, "log-level", "info", "log level: debug,info,warn,error,dpanic,panic,fatal")
	flag.Parse()

	l := level.NewAtomicLevel()
	_ = l.UnmarshalText([]byte(logLevel))
	optsl := zap.Level(&l)
	ctrl.SetLogger(zap.New(zap.UseDevMode(true), optsl))

	//webhook
	certDir := os.TempDir() + "/podpreset-webhook/serving-certs"
	obj := make(map[string]*metav1.LabelSelector)
	namespace := make(map[string]*metav1.LabelSelector)
	//webhook
	namespaceName := Env("NAMESPACE_NAME", "kube-system")
	svcName := Env("SVC_NAME", "webhook-service")
	secretName := Env("SECRET_NAME", "webhook-secret")
	csrName := Env("CSR_NAME", "webhook-csr")
	mutateName := Env("MUTATING_NAME", "")
	w := &wk.CertWebHook{
		Subject:     nil,
		CertDir:     certDir,
		Namespace:   namespaceName,
		ServiceName: svcName,
		SecretName:  secretName,
		CsrName:     csrName,
		WebHook: []wk.WebHook{
			{MutatingName: mutateName, ObjectSelect: obj, NamespaceSelect: namespace},
		},
	}
	err := w.Init()
	if err != nil {
		setupLog.Error(err, "unable to create init client")
		os.Exit(1)
	}
	err = w.Generator()
	if err != nil {
		setupLog.Error(err, "unable to generator cert files")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		CertDir:                certDir,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "273cddc9.cuisongliu.com",
		HealthProbeBindAddress: healthAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	if err := webhooks.SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook certs")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder
	//healthz  Liveness
	if err := mgr.AddHealthzCheck("check", func(req *http.Request) error {
		return nil
	}); err != nil {
		setupLog.Error(err, "problem running manager liveness Check")
		os.Exit(1)
	}
	//readyz   Readiness
	if err := mgr.AddReadyzCheck("check", func(req *http.Request) error {
		return nil
	}); err != nil {
		setupLog.Error(err, "problem running manager readiness check")
		os.Exit(1)
	}
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func Env(env, defaultVal string) string {
	var returnVar string
	if len(env) == 0 {
		return ""
	}
	returnVar = os.Getenv(env)
	if returnVar == "" {
		return defaultVal
	}
	return returnVar
}
