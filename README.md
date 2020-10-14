node-manage-operator
===================

An operator to apply specific labels to the nodes based on the provided CR. The operator has been introduced to support example-cnf application deployment, which needs to allocate the packet generator TRex to only one node and DUT TestPMD can be across multiple modes. With this labelling, we will also ensure that TRex and TestPMD is not deployed on a same node.

Deploy the operator
-------------------
This operator is bundled with catalog index image of [NFV Example CNF Catalog](https://quay.io/repository/krsacme/nfv-example-cnf-catalog). Deploy this catalog index image in the cluster using `CatalogSource`. Alternatively, use this role [`example-cnf-catalog`](https://github.com/rh-nfv-int/nfv-example-cnf-deploy/tree/master/roles/example-cnf-catalog) to deploy the catalog.

Once the cataglog is deployed, then create a subscription of the `node-manage-operator` to deploy this operator.

```
apiVersion: v1
kind: Namespace
metadata:
  name: example-cnf
spec: {}
---
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: example-cnf-operator-group
  namespace: example-cnf
spec:
  targetNamespaces:
    - example-cnf
---
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: node-manage-operator-subscription
  namespace: example-cnf
spec:
  channel: "alpha"
  name: node-manage-operator
  source: nfv-example-cnf-catalog
  sourceNamespace: openshift-marketplace
```

Create Custom Resource to Label nodes
-------------------------------------
The below specified custom resource `NodeLabels` applies the labels for the worker nodes as per the requirement of the example-cnf application deployment.

```
apiVersion: nodemanage.openshift.io/v1
kind: NodeLabels
metadata:
  name: nodelabels-sample
spec:
  nodeSelectorLabels:
    node-role.kubernetes.io/worker: ""
  labelGroup:
    - count: 1
      labels:
        examplecnf.openshift.io/trex-new: ""
    - count: 2
      labels:
        examplecnf.openshift.io/testpmd-new: ""
```
