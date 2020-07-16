# Rate Limit Operator

## Deploy the ratelimit-operator on Openshift cluster

* Create the sample CRD
~~~
$ oc create -f deploy/crds/operators.example.com_ratelimiters_crd.yaml
~~~

* Deploy the Operator along with set-up the RBAC
~~~
$ oc create -f deploy/service_account.yaml
$ oc create -f deploy/role.yaml
$ oc create -f deploy/role_binding.yaml
$ oc create -f deploy/operator.yaml
~~~

* Create the RateLimiter Custom Resource(CR)
~~~
$ oc apply -f deploy/crds/operators.example.com_v1_ratelimiter_cr.yaml
~~~

* Verify the application deployment and POD has been created
~~~
$ oc get deployment
NAME                 READY   UP-TO-DATE   AVAILABLE   
ratelimit-service     1/1        1            1        
ratelimit-operator    1/1        1            1        

$ oc get pods
NAME                                  READY     STATUS    RESTARTS  
ratelimit-service-5bf67b784f-np424     3/3     Running       0       
ratelimit-operator-5c988b77c6-cjfcg    1/1     Running       0       

$ oc get ratelimiter
NAME                  AGE
ratelimit-service    5m23s
