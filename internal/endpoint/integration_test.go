package endpoint

import (
	"flag"
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Tests that even with publishNotReadyAddresses: true in the headless services,
// the endpoints are not always returned during rollouts. Just start a roll-out in
// dogfood and run this test.
//
// k -n dogfood-k8s rollout restart sts/indexed-search
// go test -integration -run=Integration
// === RUN   TestIntegrationK8SNotReadyAddressesBug
// DBUG[08-11|16:49:14] kubernetes endpoints  service=indexed-search urls="[indexed-search-0.indexed-search indexed-search-1.indexed-search]"
// DBUG[08-11|16:49:14] kubernetes endpoints  service=indexed-search urls="[indexed-search-0.indexed-search indexed-search-1.indexed-search]"
// DBUG[08-11|16:49:28] kubernetes endpoints  service=indexed-search urls="[indexed-search-0.indexed-search indexed-search-1.indexed-search]"
// DBUG[08-11|16:49:40] kubernetes endpoints  service=indexed-search urls=[indexed-search-0.indexed-search]
//     endpoint_test.go:163: endpoint set has shrunk from 2 to 1
// 	--- FAIL: TestIntegrationK8SNotReadyAddressesBug (26.94s)

var integration = flag.Bool("integration", false, "Run integration tests")

func TestIntegrationK8SNotReadyAddressesBug(t *testing.T) {
	if !*integration {
		t.Skip("Not running integration tests")
	}

	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatal(err)
	}

	m := New("k8s+rpc://indexed-search?type=statefulset")
	m.cli = clientset
	m.ns = "dogfood-k8s"

	scale, err := m.cli.AppsV1().StatefulSets(m.ns).GetScale("indexed-search", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("scale: %v", scale.Status.Replicas)

	factory := informers.NewSharedInformerFactoryWithOptions(m.cli, 0,
		informers.WithNamespace("dogfood-k8s"),
		informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
			opts.FieldSelector = "metadata.name=searcher"
		}),
	)

	informer := factory.Apps().V1().Deployments().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	endpoints := func(op string, obj interface{}) {
		t.Helper()

		var eps []string

		switch o := (obj).(type) {
		case *corev1.Endpoints:
			for _, s := range o.Subsets {
				for _, a := range s.Addresses {
					eps = append(eps, a.Hostname)
				}
			}
		case *corev1.Pod:
			eps = append(eps, o.Name)
		case *appsv1.Deployment:
			replicas := int32(1)
			if o.Spec.Replicas != nil {
				replicas = *o.Spec.Replicas
			}

			for i := int32(0); i < replicas; i++ {
				eps = append(eps, fmt.Sprintf("%s-%d", o.Name, i))
			}
		}

		sort.Strings(eps)
		t.Logf("%s: %v", op, eps)
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			endpoints("add", obj)
		},
		DeleteFunc: func(obj interface{}) {
			endpoints("delete", obj)
		},
		UpdateFunc: func(_, obj interface{}) {
			endpoints("updated", obj)
		},
	})

	informer.Run(stopper)
	// count := 0
	// began := time.Now()

	//for time.Since(began) <= time.Minute {
	//eps, err := m.Endpoints()
	//if err != nil {
	//t.Fatal(err)
	//}

	//if count == 0 {
	//count = len(eps)
	//} else if len(eps) < count {
	//t.Fatalf("endpoint set has shrunk from %d to %d", count, len(eps))
	//}

	//time.Sleep(500 * time.Millisecond)
	//}
}
