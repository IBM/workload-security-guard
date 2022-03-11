/*
Copyright 2017 The Kubernetes Authors.

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
	"fmt"

	//kubeinformers "k8s.io/client-go/informers"
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	wsecurity "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	guardianclientset "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians"
	//informers "wsecurity.ibm.com/wsecurity/pkg/generated/informers/externalversions"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	//stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	//kubeClient, err := kubernetes.NewForConfig(cfg)
	//if err != nil {
	//	klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	//}

	exampleClient, err := guardianclientset.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	client := exampleClient.WsecurityV1()

	var g *wsecurity.Guardian

	g, err = client.Guardians("default").Get(context.TODO(), "mytest", v1.GetOptions{})
	if err == nil {
		fmt.Printf("guardian %v\n", g)
	} else {
		fmt.Printf("get err %v\n", err)
		g = &wsecurity.Guardian{}
		g.Kind = "Guardian"
		g.APIVersion = "wsecurity.ibmresearch.com/v1"
		g.Name = "mytest"

		fmt.Printf("creating guardian %v\n", g)
		g, err = client.Guardians("default").Create(context.TODO(), g, v1.CreateOptions{})
		fmt.Printf("create err %v\n", err)
		fmt.Printf("created %v\n", g)
	}

	//update
	//g = &wsecurity.Guardian{}
	g.Kind = "Guardian"
	g.APIVersion = "wsecurity.ibmresearch.com/v1"
	g.Name = "mytest"
	g.Spec = new(wsecurity.GuardianSpec)
	g.Spec.Control.Consult = true
	g.Spec.Control.RequestsPerMinuete = 60
	g.Spec.Contigured.Req.AddTypicalVal()
	/*g.Spec.ForceAllow = true
	g.Spec.Req.Url.Segments = make(wsecurity.U8MinmaxSlice, 1)
	g.Spec.Req.Url.Segments[0].Max = 4
	g.Spec.Req.Url.Segments[0].Min = 2
	g.Spec.Req.Url.Val.Numbers = make(wsecurity.U8MinmaxSlice, 1)
	g.Spec.Req.Url.Val.Numbers[0].Max = 7
	g.Spec.Req.Url.Val.Flags = 0x75381
	*/
	fmt.Printf("updating %v\n", g)
	g, err = client.Guardians("default").Update(context.TODO(), g, v1.UpdateOptions{})
	fmt.Printf("update err %v\n", err)
	fmt.Printf("updated %v\n", g)

	g, err = client.Guardians("default").Get(context.TODO(), "mytest", v1.GetOptions{})
	fmt.Printf("last get err %v\n", err)
	fmt.Printf("guardian %v\n", g)
	fmt.Printf("guardian rpm %d\n", g.Spec.Control.RequestsPerMinuete)
	//kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	//exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	//controller := NewController(kubeClient, exampleClient,
	//	kubeInformerFactory.Apps().V1().Deployments(),
	//	exampleInformerFactory.Samplecontroller().V1alpha1().Foos())

	// notice that there is no need to run Start methods in a separate goroutine. (i.e. go kubeInformerFactory.Start(stopCh)
	// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
	//kubeInformerFactory.Start(stopCh)
	//exampleInformerFactory.Start(stopCh)

	//if err = controller.Run(2, stopCh); err != nil {
	//	klog.Fatalf("Error running controller: %s", err.Error())
	//}
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
