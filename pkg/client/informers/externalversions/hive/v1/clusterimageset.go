// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	hivev1 "github.com/openshift/hive/pkg/apis/hive/v1"
	versioned "github.com/openshift/hive/pkg/client/clientset/versioned"
	internalinterfaces "github.com/openshift/hive/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/openshift/hive/pkg/client/listers/hive/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ClusterImageSetInformer provides access to a shared informer and lister for
// ClusterImageSets.
type ClusterImageSetInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ClusterImageSetLister
}

type clusterImageSetInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewClusterImageSetInformer constructs a new informer for ClusterImageSet type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewClusterImageSetInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredClusterImageSetInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredClusterImageSetInformer constructs a new informer for ClusterImageSet type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredClusterImageSetInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HiveV1().ClusterImageSets().List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HiveV1().ClusterImageSets().Watch(options)
			},
		},
		&hivev1.ClusterImageSet{},
		resyncPeriod,
		indexers,
	)
}

func (f *clusterImageSetInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredClusterImageSetInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *clusterImageSetInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&hivev1.ClusterImageSet{}, f.defaultInformer)
}

func (f *clusterImageSetInformer) Lister() v1.ClusterImageSetLister {
	return v1.NewClusterImageSetLister(f.Informer().GetIndexer())
}
