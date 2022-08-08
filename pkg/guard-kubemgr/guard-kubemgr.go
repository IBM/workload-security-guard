package guardkubemgr

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	spec "github.com/IBM/workload-security-guard/pkg/apis/wsecurity/v1"
	guardianclientset "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians"
	v1 "github.com/IBM/workload-security-guard/pkg/generated/clientset/guardians/typed/wsecurity/v1"

	pi "github.com/IBM/go-security-plugs/pluginterfaces"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Kubemgr struct {
	clientset *kubernetes.Clientset
	gClient   v1.WsecurityV1Interface
}

func (k *Kubemgr) InitConfigs() {
	var kubeconfig *string
	var cfg *rest.Config
	var errInCluster error
	var errOutOfCluster error

	// creates the in-cluster config
	cfg, errInCluster = rest.InClusterConfig()
	if errInCluster != nil {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		// use the current context in kubeconfig
		cfg, errOutOfCluster = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if errOutOfCluster != nil {
			panic(fmt.Sprintf("No Config found! errInCluster %s", errInCluster.Error()))
		}
	}

	guardianClient, err := guardianclientset.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}
	k.gClient = guardianClient.WsecurityV1()

	var errClientset error
	k.clientset, errClientset = kubernetes.NewForConfig(cfg)
	if errClientset != nil {
		panic(errClientset.Error())
	}
}

//
func (k *Kubemgr) ReadCrd(namespace string, serviceId string) *spec.GuardianSpec {
	g, err := k.gClient.Guardians(namespace).Get(context.TODO(), serviceId, metav1.GetOptions{})
	if err != nil {
		pi.Log.Infof("Missing Guardian CRD %s.%s (Err %s)", serviceId, namespace, err.Error())
		//panic(fmt.Sprintf("No Guardian! for %s.%s", serviceId, namespace))
		return nil
	}
	pi.Log.Debugf("Found Guardian CRD %s.%s", serviceId, namespace)
	//fmt.Print((*spec.WsGate)(g.Spec).Marshal(0))
	(*spec.WsGate)(g.Spec).Reconcile()
	return g.Spec
}

func (k *Kubemgr) ReadConfigMap(namespace string, sid string) *spec.GuardianSpec {
	cmname := "guardian." + sid
	cm, err := k.clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), cmname, metav1.GetOptions{})
	if err != nil {
		pi.Log.Infof("ConfigMap Error: %v\n", err)
		return nil
	}

	g := new(spec.GuardianSpec)
	gdata, ok := cm.Data["Guardian"]
	if ok && len(gdata) > 0 {
		gbytes := []byte(gdata)
		jsonErr := json.Unmarshal(gbytes, g)
		if jsonErr != nil {
			pi.Log.Infof("wsgate getConfig: unmarshel error %v", jsonErr)
			return nil
		}
		(*spec.WsGate)(g).Reconcile()
	}
	pi.Log.Debugf("Get %s ConfigMap succesful", cmname)
	return g
	//	err = json.Unmarshal([]byte(cm.Data["guardian"]), &p.wsGate)
	//	if err != nil {
	//		fmt.Printf("ConfigMap Unmarshal Error: %v\n", err)
	//		panic(err.Error())
	//	}
}

// Set a Guardian Config Map (Update if exists, create if not)
// Not using Kube's apply tp avoid kube's server side merge,
// ...so we need to do this ugly thingy...
// This ugly thingy has a race condition!
// Sometimes it will fail and we will just return false
func (k *Kubemgr) SetCm(ns string, sid string, guardianSpec *spec.GuardianSpec) string {
	var gbytes []byte
	cmname := "guardian." + sid
	//fmt.Printf("SetCm %s: data: %v\n", cmname, guardianSpec)
	cm, err := k.clientset.CoreV1().ConfigMaps(ns).Get(context.TODO(), cmname, metav1.GetOptions{})
	if err == nil {
		//fmt.Printf("setCm %s: guardian read succesful \n", cmname)
		g := new(spec.GuardianSpec)
		gdata, ok := cm.Data["Guardian"]
		if ok && len(gdata) > 0 { // Guardian exists
			gbytes = []byte(gdata)
			if err := json.Unmarshal(gbytes, g); err != nil {
				return fmt.Sprintf("SetCm %s: Unmarshel error %s", cmname, err.Error())
			}
			(*spec.WsGate)(g).Reconcile()
		} else {
			if cm.Data == nil {
				cm.Data = make(map[string]string, 1)
			}
		}
		if guardianSpec != nil {
			if guardianSpec.Control != nil {
				g.Control = guardianSpec.Control
			}
			if guardianSpec.Configured != nil {
				g.Configured = guardianSpec.Configured
			}
			if guardianSpec.Control != nil {
				g.Learned = guardianSpec.Learned
			}
		}
		gbytes, err = json.Marshal(g)
		if err != nil {
			//fmt.Printf("SetCm %s: Error marshaling data: %s\n", cmname, err.Error())
			return fmt.Sprintf("SetCm %s: Error marshaling data: %s", cmname, err.Error())
		}
		cm.Data["Guardian"] = string(gbytes)

		//fmt.Printf("SetCm %s: After marshaling data: %s\n", cmname, gbytes)

		//fmt.Printf("SetCm %s: cm.Data[\"Guardian\"]: %s\n", cmname, cm.Data["Guardian"])

		_, err = k.clientset.CoreV1().ConfigMaps(ns).Update(context.TODO(), cm, metav1.UpdateOptions{})
		if err != nil {
			//fmt.Printf("SetCm %s: Error updating resource: %s\n", cmname, err.Error())
			return fmt.Sprintf("SetCm %s: Error updating resource: %s", cmname, err.Error())
		}

		//fmt.Printf("setCm %s: guardian update succesfull\n", cmname)
	} else {
		//fmt.Printf("setCm %s: guardian read err %s\n", cmname, err.Error())
		cm = new(corev1.ConfigMap)
		cm.Name = cmname
		cm.Data = make(map[string]string, 1)
		gbytes, err = json.Marshal(guardianSpec)
		if err != nil {
			//fmt.Printf("SetCm %s: Error marshaling data during create: %s\n", cmname, err.Error())
			return fmt.Sprintf("SetCm %s: Error marshaling data during create: %s", cmname, err.Error())
		}
		cm.Data["Guardian"] = string(gbytes)
		_, err = k.clientset.CoreV1().ConfigMaps(ns).Create(context.TODO(), cm, metav1.CreateOptions{})
		if err != nil {
			//fmt.Printf("SetCm %s: Error creating resource: %s\n", cmname, err.Error())
			return fmt.Sprintf("SetCm %s: Error creating resource: %s", cmname, err.Error())
		}
		//fmt.Printf("setCm %s: guardian create succesfull\n", cmname)
	}

	return ""
}

// Set a Guardian Custom Resource (Update if exists, create if not)
// Kube does not support server side udpate with create if not exist,
// ...so we need to do this ugly thingy...
// This ugly thingy has a race condition!
// Sometimes it will fail and we will just return false
func (k *Kubemgr) SetCrd(ns string, sid string, guardianSpec *spec.GuardianSpec) string {
	var g *spec.Guardian
	var err error
	g, err = k.gClient.Guardians(ns).Get(context.TODO(), sid, metav1.GetOptions{})
	if err == nil {
		//if g.Spec.Learned != nil {
		//	fmt.Printf("setCrd: guardian read Method %d %v\n", len(g.Spec.Learned.Req.Method), g.Spec.Learned.Req.Method)
		//} else {
		//	fmt.Printf("setCrd: guardian read g.Spec.Learned is nil\n")
		//}
		g.Name = sid
		if guardianSpec != nil {
			if guardianSpec.Control != nil {
				g.Spec.Control = guardianSpec.Control
			}
			if guardianSpec.Configured != nil {
				g.Spec.Configured = guardianSpec.Configured
			}
			if guardianSpec.Control != nil {
				g.Spec.Learned = guardianSpec.Learned
			}
		}
		_, err = k.gClient.Guardians(ns).Update(context.TODO(), g, metav1.UpdateOptions{})
		if err != nil {
			//fmt.Printf("setCrd: update err %v\n", err)
			return fmt.Sprintf("SetCrd: Error updating resource: %s", err.Error())
		}
		//fmt.Printf("setCrd: guardian update succesfull %v\n", g.Spec.Control)
	} else {
		g = new(spec.Guardian)
		fmt.Printf("setCrd: guardian read err %v\n", err)
		fmt.Printf("setCrd: creating guardian ns %s sid %s guardianSpec %v\n", ns, sid, guardianSpec)

		g.Name = sid
		g.Spec = guardianSpec

		_, err = k.gClient.Guardians(ns).Create(context.TODO(), g, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("setCrd: guardian create err %v\n", err)
			return fmt.Sprintf("SetCrd: Error creating resource: %s", err.Error())
		}
		fmt.Printf("setCrd: guardian create succesfull %v\n", g)
	}

	//fmt.Printf("setCrd: success!\n")
	return ""
}

/*
// Set a Guardian Config Map (Update if exists, create if not)
// using Apply
func (k *Kubemgr) SetCmApply(ns string, cmname string, guardianSpec *spec.GuardianSpec) string {
	gbytes, err := json.Marshal(guardianSpec)
	if err != nil {
		//fmt.Printf("setCm: marshal err %v\n", err)
		return fmt.Sprintf("SetCm: Error marshaling data: %s", err.Error())
	}

	cm := applyv1.ConfigMap(cmname, ns)

	cm.Data = make(map[string]string, 1)
	cm.Data["Guardian"] = string(gbytes)
	_, err = k.clientset.CoreV1().ConfigMaps(ns).Apply(context.TODO(), cm, metav1.ApplyOptions{
		FieldManager: "guard",
		Force:        true,
	})
	if err != nil {
		//fmt.Printf("setCm: update err %v\n", err)
		return fmt.Sprintf("SetCm: Error updating resource: %s", err.Error())
	}

	return ""
}

// Set a Guardian Custom Resource (Update if exists, create if not)
// using Apply
func (k *Kubemgr) SetCrdApply(ns string, sid string, guardianSpec *spec.GuardianSpec) string {
	g := new(spec.Guardian)
	g.Name = sid
	g.Spec = guardianSpec
	g.APIVersion = "wsecurity.ibmresearch.com/v1"
	g.Kind = "Guardian"
	gbytes, err := json.Marshal(g)
	if err != nil {
		//fmt.Printf("setCrd: marshal err %v\n", err)
		return fmt.Sprintf("SetCrd: Error marshaling data: %s", err.Error())
	}

	forcetrue := new(bool)
	*forcetrue = true
	_, err = k.gClient.Guardians(ns).Patch(context.TODO(), sid, types.ApplyPatchType, gbytes, metav1.PatchOptions{
		FieldManager: "guard",
		Force:        forcetrue,
	})
	if err != nil {
		//fmt.Printf("setCrd: update err %v\n", err)
		return fmt.Sprintf("SetCrd: Error updating resource: %s", err.Error())
	}

	fmt.Printf("setCrd: success!\n")
	return ""
}
*/

func (k *Kubemgr) FetchConfig(ns string, sid string, cm bool) *spec.GuardianSpec {
	var gurdianSpec *spec.GuardianSpec
	if !strings.EqualFold(sid, "ns-"+ns) {
		if cm {
			gurdianSpec = k.ReadConfigMap(ns, sid)
			if gurdianSpec == nil {
				gurdianSpec = k.ReadConfigMap(ns, "ns-"+ns)
			}
		} else {
			gurdianSpec = k.ReadCrd(ns, sid)
			if gurdianSpec == nil {
				gurdianSpec = k.ReadCrd(ns, "ns-"+ns)
				//if gurdianSpec == nil { forbiden to read from knative-serving
				//	gurdianSpec = k.ReadCrd("knative-serving", "guardian")
				//}
			}
		}
	}

	if gurdianSpec == nil {
		gurdianSpec = new(spec.GuardianSpec)
		(*spec.WsGate)(gurdianSpec).Reconcile()
		// default gurdianSpec has:
		// 		gurdianSpec.falseAllow=false
		// 		gurdianSpec.ConsultGuard.Active = false
		(*spec.WsGate)(gurdianSpec).AutoActivate()
		// now gurdianSpec has:
		// 		gurdianSpec.falseAllow=false
		// 		gurdianSpec.ConsultGuard.Active = false

	}
	return gurdianSpec

	/*
		req, err := http.NewRequest(http.MethodGet, p.guardUrl+"/config", nil)
		if err != nil {
			pi.Log.Infof("wsgate getConfig: http.NewRequest error %v", err)
		}
		query := req.URL.Query()
		query.Add("sid", p.id)
		req.URL.RawQuery = query.Encode()
		res, getErr := p.httpc.Do(req)
		if getErr != nil {
			pi.Log.Infof("wsgate getConfig: httpc.Do error %v", getErr)
			return
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			pi.Log.Infof("wsgate getConfig: read body error %v", readErr)
		}

		//pi.Log.Infof("wsgate getConfig: unmarshal %s", string(body))
		jsonErr := json.Unmarshal(body, &p.gateConfig)
		if jsonErr != nil {
			pi.Log.Infof("wsgate getConfig: unmarshel error %v", jsonErr)
		}

		//pi.Log.Infof("wsgate getConfig: ended %v ", p.gateConfig)
	*/
}

func (k *Kubemgr) WatchOnce(ns string, set func(ns string, sid string, g *spec.GuardianSpec)) {
	defer func() {
		if recovered := recover(); recovered != nil {
			fmt.Printf("Recovered from panic during watchCrdOnce! Recover: %v\n", recovered)
		}
	}()
	watcherCrd, err := k.gClient.Guardians(ns).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("watchCrd gClient.Guardians(%s).Watch err %v\n", ns, err)
		return
	}
	chCrd := watcherCrd.ResultChan()
	watcherCm, err := k.clientset.CoreV1().ConfigMaps(ns).Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("watchCrd gClient.Guardians(%s).Watch err %v\n", ns, err)
		return
	}
	chCm := watcherCm.ResultChan()
	for {
		select {
		case event, ok := <-chCrd:
			if !ok {
				// the channel got closed, so we need to restart
				fmt.Printf("Kubernetes hung up on us, restarting event watcher\n")
				return
			}

			// handle the event
			fmt.Printf("crd event \n")
			//fmt.Printf("--------------> Event\n")

			switch event.Type {
			case watch.Deleted:
				fallthrough
			case watch.Modified:
				fallthrough
			case watch.Added:
				g, ok := event.Object.(*spec.Guardian)
				if !ok {
					fmt.Printf("Kubernetes cant convert to type Guardian\n")
					return
				}
				ns := g.ObjectMeta.Namespace
				sid := g.ObjectMeta.Name

				if event.Type == watch.Deleted {
					set(ns, sid, nil)
					continue
				}
				set(ns, sid, g.Spec)
			case watch.Error:
				s := event.Object.(*metav1.Status)
				fmt.Printf("Error during watch CRD: \n\tListMeta %v\n\tTypeMeta %v\n", s.ListMeta, s.TypeMeta)
			}
		case event, ok := <-chCm:
			if !ok {
				// the channel got closed, so we need to restart
				fmt.Printf("Kubernetes hung up on us, restarting event watcher\n")
				return
			}

			// handle the event
			fmt.Printf("cm event\n")
			//fmt.Printf("--------------> Event\n")

			switch event.Type {
			case watch.Deleted:
				fallthrough
			case watch.Modified:
				fallthrough
			case watch.Added:
				cm, ok := event.Object.(*corev1.ConfigMap)
				if !ok {
					fmt.Printf("Kubernetes cant convert to type ConfigMap\n")
					return
				}
				ns := cm.ObjectMeta.Namespace
				sid := cm.ObjectMeta.Name
				if !strings.HasPrefix(sid, "guardian.") || strings.HasPrefix(sid, "guardian.ns.") {
					// skip...
					continue
				}
				if event.Type == watch.Deleted {
					set(ns, sid, nil)
					continue
				}
				if cm.Data["Guardian"] == "" {
					set(ns, sid, nil)
					continue
				}
				g := new(spec.GuardianSpec)
				gdata := []byte(cm.Data["Guardian"])
				jsonErr := json.Unmarshal(gdata, g)
				if jsonErr != nil {
					pi.Log.Infof("wsgate getConfig: unmarshel error %v", jsonErr)
					set(ns, sid, nil)
					continue
				}
				(*spec.WsGate)(g).Reconcile()
				set(ns, sid, g)
			case watch.Error:
				s := event.Object.(*metav1.Status)
				fmt.Printf("Error during watch CM: \n\tListMeta %v\n\tTypeMeta %v\n", s.ListMeta, s.TypeMeta)
			}
		case <-time.After(10 * time.Minute):
			// deal with the issue where we get no events
			fmt.Printf("Timeout, restarting event watcher\n")
			return
		}
	}
}
