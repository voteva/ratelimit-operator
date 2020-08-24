# Rate Limit Operator

## Deploy the ratelimit-operator on Openshift cluster

* Create the sample CRD
~~~
$ oc create -f deploy/crds/operators.example.com_ratelimiters_crd.yaml
$ oc create -f deploy/crds/operators.example.com_ratelimiterconfigs_crd.yaml
~~~
* Select the namespace
~~~
$ oc project <your-namespace>
~~~
* Deploy the Operator along with set-up the RBAC
~~~
$ oc create -f deploy/service_account.yaml
$ oc create -f deploy/role.yaml
$ oc create -f deploy/role_binding.yaml
$ oc create -f deploy/operator.yaml
~~~

* Create the RateLimiter and RateLimiterConfig Custom Resources (CR)
~~~
$ oc apply -f deploy/crds/operators.example.com_v1_ratelimiter_cr.yaml
$ oc apply -f deploy/crds/operators.example.com_v1_ratelimiterconfig_cr.yaml
~~~