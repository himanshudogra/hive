// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/openshift/hive/pkg/apis/hive/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// SyncSetInstanceLister helps list SyncSetInstances.
type SyncSetInstanceLister interface {
	// List lists all SyncSetInstances in the indexer.
	List(selector labels.Selector) (ret []*v1.SyncSetInstance, err error)
	// SyncSetInstances returns an object that can list and get SyncSetInstances.
	SyncSetInstances(namespace string) SyncSetInstanceNamespaceLister
	SyncSetInstanceListerExpansion
}

// syncSetInstanceLister implements the SyncSetInstanceLister interface.
type syncSetInstanceLister struct {
	indexer cache.Indexer
}

// NewSyncSetInstanceLister returns a new SyncSetInstanceLister.
func NewSyncSetInstanceLister(indexer cache.Indexer) SyncSetInstanceLister {
	return &syncSetInstanceLister{indexer: indexer}
}

// List lists all SyncSetInstances in the indexer.
func (s *syncSetInstanceLister) List(selector labels.Selector) (ret []*v1.SyncSetInstance, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.SyncSetInstance))
	})
	return ret, err
}

// SyncSetInstances returns an object that can list and get SyncSetInstances.
func (s *syncSetInstanceLister) SyncSetInstances(namespace string) SyncSetInstanceNamespaceLister {
	return syncSetInstanceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SyncSetInstanceNamespaceLister helps list and get SyncSetInstances.
type SyncSetInstanceNamespaceLister interface {
	// List lists all SyncSetInstances in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.SyncSetInstance, err error)
	// Get retrieves the SyncSetInstance from the indexer for a given namespace and name.
	Get(name string) (*v1.SyncSetInstance, error)
	SyncSetInstanceNamespaceListerExpansion
}

// syncSetInstanceNamespaceLister implements the SyncSetInstanceNamespaceLister
// interface.
type syncSetInstanceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all SyncSetInstances in the indexer for a given namespace.
func (s syncSetInstanceNamespaceLister) List(selector labels.Selector) (ret []*v1.SyncSetInstance, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.SyncSetInstance))
	})
	return ret, err
}

// Get retrieves the SyncSetInstance from the indexer for a given namespace and name.
func (s syncSetInstanceNamespaceLister) Get(name string) (*v1.SyncSetInstance, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("syncsetinstance"), name)
	}
	return obj.(*v1.SyncSetInstance), nil
}
