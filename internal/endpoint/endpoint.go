// Package endpoint provides a consistent hash map over service endpoints.
package endpoint

import (
	"fmt"
	"hash/crc32"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"
	"github.com/inconshreveable/log15"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// Map is a consistent hash map to URLs. It uses the kubernetes API to watch
// the endpoints for a service and update the map when they change. It can
// also fallback to static URLs if not configured for kubernetes.
type Map struct {
	mu      sync.Mutex
	init    func() (*hashMap, error)
	err     error
	urls    *hashMap
	urlspec string
	cli     *kubernetes.Clientset
	ns      string
}

// New creates a new Map for the URL specifier.
//
// If the scheme is prefixed with "k8s+", one URL is expected and the format is
// expected to match e.g. k8s+http://service.namespace:port/path. namespace,
// port and path are optional. URLs of this form will consistently hash among
// the endpoints for the Kubernetes service. The values returned by Get will
// look like http://endpoint:port/path.
//
// If the scheme is not prefixed with "k8s+", a space separated list of URLs is
// expected. The map will consistently hash against these URLs in this case.
// This is useful for specifying non-Kubernetes endpoints.
//
// Examples URL specifiers:
//
// 	"k8s+http://searcher"
// 	"http://searcher-0 http://searcher-1 http://searcher-2"
//
func New(urlspec string) *Map {
	if !strings.HasPrefix(urlspec, "k8s+") {
		return &Map{
			urlspec: urlspec,
			urls:    newConsistentHashMap(strings.Fields(urlspec)),
		}
	}

	m := &Map{urlspec: urlspec}

	// Kick off setting the initial urls or err on first access. We don't rely
	// just on inform since it may not communicate updates.
	m.init = func() (*hashMap, error) {
		u, err := parseURL(urlspec)
		if err != nil {
			return nil, err
		}

		if m.cli == nil {
			m.cli, m.ns, err = loadClient()
			if err != nil {
				return nil, err
			}
		}

		factory := informers.NewSharedInformerFactoryWithOptions(m.cli, 0,
			informers.WithTweakListOptions(func(opts *metav1.ListOptions) {
				opts.FieldSelector = "metadata.name=" + u.Service
			}),
		)

		var informer cache.SharedIndexInformer
		switch u.Kind {
		case "sts", "statefulset":
			informer = factory.Apps().V1().StatefulSets().Informer()
		default:
			informer = factory.Core().V1().Endpoints().Informer()
		}

		handle := func(op string, obj interface{}) {
			var eps []string

			switch o := (obj).(type) {
			case *corev1.Endpoints:
				for _, s := range o.Subsets {
					addrs := append([]corev1.EndpointAddress(nil), s.Addresses...)
					addrs = append(addrs, s.NotReadyAddresses...)

					for _, a := range addrs {
						ep := a.Hostname
						if ep == "" {
							ep = a.IP
						}
						eps = append(eps, ep)
					}
				}
			case *appsv1.StatefulSet:
				replicas := int32(1)
				if o.Spec.Replicas != nil {
					replicas = *o.Spec.Replicas
				}
				for i := int32(0); i < replicas; i++ {
					eps = append(eps, fmt.Sprintf("%s-%d", o.Name, i))
				}
			}
		}

		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
			},
			UpdateFunc: func(_, obj interface{}) {
			},
		})

		stopper := make(chan struct{})
		defer close(stopper)

		go informer.Run(stopper)

		endpoints := []string{}
		return endpointsToMap(u.Service, endpoints)
	}

	return m
}

// Static returns an Endpoint map which consistently hashes over endpoints.
//
// There are no requirements on endpoints, it can be any arbitrary
// string. Unlike static endpoints created via New.
//
// Static Maps are guaranteed to never return an error.
func Static(endpoints ...string) *Map {
	return &Map{
		urlspec: fmt.Sprintf("%v", endpoints),
		urls:    newConsistentHashMap(endpoints),
	}
}

// Empty returns an Endpoint map which always fails with err.
func Empty(err error) *Map {
	return &Map{
		urlspec: "error: " + err.Error(),
		err:     err,
	}
}

func (m *Map) String() string {
	return fmt.Sprintf("endpoint.Map(%s)", m.urlspec)
}

// Get the closest URL in the hash to the provided key that is not in
// exclude. If no URL is found, "" is returned.
//
// Note: For k8s URLs we return URLs based on the registered endpoints. The
// endpoint may not actually be available yet / at the moment. So users of the
// URL should implement a retry strategy.
func (m *Map) Get(key string, exclude map[string]bool) (string, error) {
	urls, err := m.getUrls()
	if err != nil {
		return "", err
	}

	return urls.get(key, exclude), nil
}

// GetMany is the same as calling Get on each item of keys. It will only
// acquire the underlying endpoint map once, so is preferred to calling Get
// for each key which will acquire the endpoint map for each call. The benefit
// is it is faster (O(1) mutex acquires vs O(n)) and consistent (endpoint map
// is immutable vs may change between Get calls).
func (m *Map) GetMany(keys ...string) ([]string, error) {
	urls, err := m.getUrls()
	if err != nil {
		return nil, err
	}

	vals := make([]string, len(keys))
	for i := range keys {
		vals[i] = urls.get(keys[i], nil)
	}
	return vals, nil
}

// Endpoints returns a set of all addresses. Do not modify the returned value.
func (m *Map) Endpoints() (map[string]struct{}, error) {
	urls, err := m.getUrls()
	if err != nil {
		return nil, err
	}

	return urls.values, nil
}

func (m *Map) getUrls() (*hashMap, error) {
	m.mu.Lock()
	if m.init != nil {
		m.urls, m.err = m.init()
		m.init = nil // prevent running again
	}
	urls, err := m.urls, m.err
	m.mu.Unlock()
	return urls, err
}

func endpointsToMap(service string, eps []string) (*hashMap, error) {
	sort.Strings(eps)
	log15.Debug("kubernetes endpoints", "service", service, "endpoints", eps)
	metricEndpointSize.WithLabelValues(service).Set(float64(len(eps)))
	if len(eps) == 0 {
		return nil, errors.Errorf(
			"no %s endpoints could be found (this may indicate more %s replicas are needed, contact support@sourcegraph.com for assistance)",
			service,
			service,
		)
	}
	return newConsistentHashMap(eps), nil
}

type k8sURL struct {
	url.URL

	Service   string
	Namespace string
	Kind      string
}

func (u *k8sURL) endpointURL(endpoint string) string {
	uCopy := u.URL
	if port := u.Port(); port != "" {
		uCopy.Host = endpoint + ":" + port
	} else {
		uCopy.Host = endpoint
	}
	if uCopy.Scheme == "rpc" {
		return uCopy.Host
	}
	return uCopy.String()
}

func parseURL(rawurl string) (*k8sURL, error) {
	u, err := url.Parse(strings.TrimPrefix(rawurl, "k8s+"))
	if err != nil {
		return nil, err
	}

	parts := strings.Split(u.Hostname(), ".")
	var svc, ns string
	switch len(parts) {
	case 1:
		svc = parts[0]
	case 2:
		svc, ns = parts[0], parts[1]
	default:
		return nil, errors.Errorf("invalid k8s url. expected k8s+http://service.namespace:port/path?kind=$kind, got %s", rawurl)
	}

	return &k8sURL{
		URL:       *u,
		Service:   svc,
		Namespace: ns,
		Kind:      strings.ToLower(u.Query().Get("kind")),
	}, nil
}

func newConsistentHashMap(keys []string) *hashMap {
	// 50 replicas and crc32.ChecksumIEEE are the defaults used by
	// groupcache.
	m := hashMapNew(50, crc32.ChecksumIEEE)
	m.add(keys...)
	return m
}

// namespace returns the namespace the pod is currently running in
// this is done because the k8s client we previously used set the namespace
// when the client was created, the official k8s client does not
func namespace() string {
	const filename = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	data, err := os.ReadFile(filename)
	if err != nil {
		log15.Warn("endpoint: falling back to kubernetes default namespace", "error", filename+" is empty")
		return "default"
	}

	ns := strings.TrimSpace(string(data))
	if ns == "" {
		log15.Warn("file: ", filename, " empty using \"default\" ns")
		return "default"
	}
	return ns
}

func loadClient() (client *kubernetes.Clientset, ns string, err error) {
	// Uncomment below to test against a real cluster. This is only important
	// when you are changing how we interact with the k8s API and you want to
	// test against the real thing.
	// Ensure you set your KUBECONFIG env var or your current kubeconfig will be used

	// InClusterConfig only works when running inside of a pod in a k8s
	// cluster.
	// From https://github.com/kubernetes/client-go/tree/master/examples/out-of-cluster-client-configuration
	/*
		c, err := clientcmd.NewDefaultClientConfigLoadingRules().Load()
		if err != nil {
			log15.Error("couldn't load kubeconfig")
			os.Exit(1)
		}
		clientConfig := clientcmd.NewDefaultClientConfig(*c, nil)
		config, err = clientConfig.ClientConfig()
		namespace = "prod"
	*/

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, "", err
	}
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, "", err
	}

	return client, namespace(), err
}

var metricEndpointSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "src_endpoint_k8s_size",
	Help: "The number of urls in a watched kubernetes endpoint",
}, []string{"service"})
