## Prepare App Image
### Login ACR and Build an app image
```sh
$ az acr login --name azadpk8s2y0anxpaxdq1
Login Succeeded

# Browse to App folder to build an image
$ docker build -t azadpk8s2y0anxpaxdq1.azurecr.io/test/hello-api:1.0.0 .
$ docker images
REPOSITORY                                                  TAG         IMAGE ID       CREATED        SIZE
azadpk8s2y0anxpaxdq1.azurecr.io/test/hello-api              1.0.0       ceeb0db27e0c   13 hours ago   6.24MB
...
```

### Push the image to ACR
```sh
$ docker push azadpk8s2y0anxpaxdq1.azurecr.io/test/hello-api:1.0.0
The push refers to repository [azadpk8s2y0anxpaxdq1.azurecr.io/test/hello-api]
b73751f3aaf1: Pushed
1.0.0: digest: sha256:180716df04d99a331a47149f8d0d87351f087bbb76b86c203ea8e674aec0941e size: 527
```

### Check the ACR to see the pushed image
![ACR](https://user-images.githubusercontent.com/710161/192076537-3d4b04da-d7af-427a-ae6c-fb3e45dcb5b6.png)


## Create K8s Application Deployment

Check current context (K8s cluster) which kubectl command will execute to
```sh
$ kubectl config get-contexts
CURRENT   NAME                                 CLUSTER                              AUTHINFO                                                                                      NAMESPACE
          aks-aks-sif-npr-sea-001-sa1-d-main   aks-aks-sif-npr-sea-001-sa1-d-main   clusterUser_rg-aks-sif-npr-sea-001-sa1-d-main-akscluster_aks-aks-sif-npr-sea-001-sa1-d-main   otp-api-uat
*         aks-aks-sif-prd-sea-001-sa1-p-main   aks-aks-sif-prd-sea-001-sa1-p-main   clusterUser_rg-aks-sif-prd-sea-001-sa1-p-main-akscluster_aks-aks-sif-prd-sea-001-sa1-p-main   voc-maritz-prd
          aks-aks-sif-prd-sea-002-sa1-p-main   aks-aks-sif-prd-sea-002-sa1-p-main   clusterUser_rg-aks-sif-prd-sea-002-sa1-p-main-akscluster_aks-aks-sif-prd-sea-002-sa1-p-main   voc-sugar-prd
```

Switch to desired context
```sh
$ kubectl config use-context aks-aks-sif-npr-sea-001-sa1-d-main
Switched to context "aks-aks-sif-npr-sea-001-sa1-d-main".

# To double check current context
$ kubectl config get-contexts
CURRENT   NAME                                 CLUSTER                              AUTHINFO                                                                                      NAMESPACE
*          aks-aks-sif-npr-sea-001-sa1-d-main   aks-aks-sif-npr-sea-001-sa1-d-main   clusterUser_rg-aks-sif-npr-sea-001-sa1-d-main-akscluster_aks-aks-sif-npr-sea-001-sa1-d-main   otp-api-uat
...
```

Prepare `deployment.yaml` - the deployment script to create an application service
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-api-deployment
  namespace: test-uat
spec:
  selector:
    matchLabels:
      app: hello-api
  template:
    metadata:
      labels:
        app: hello-api
    spec:
      containers:
      - name: hello-api
        image: azadpk8s2y0anxpaxdq1.azurecr.io/test/hello-api:1.0.0
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
```

Apply changes to K8s selected context (cluster)
```sh
$ kubectl apply -f deployment.yaml
Error from server (NotFound): error when creating "pod.yaml": namespaces "test-uat" not found
# this happens because we need to create "test-uat" namespace in the context first

$ kubectl create namespace test-uat
namespace/test-uat created

$ kubectl apply -f deployment.yaml
deployment.apps/hello-api-deployment created

$ kubectl get all --namespace test-uat
NAME                                        READY   STATUS    RESTARTS   AGE
pod/hello-api-deployment-85b75f6cc6-drmn5   1/1     Running   0          10m

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello-api-deployment   1/1     1            1           10m

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/hello-api-deployment-85b75f6cc6   1         1         1       10m

# there is 1 pod (pod/hello-api-deployment-85b75f6cc6-drmn5). this means there is 1 hello-api container running 
```

Test the deployed pod by remoting to the pod and test hello api
```sh
$ kubectl exec -it hello-api-deployment-85b75f6cc6-drmn5 --namespace test-uat -- sh
/ # curl http://localhost:8080
{"message": "Բարեւ"}

# to exit the remote session
/ # exit
```

Scale out 1 more hello-api service to make 2 total hello-api services.

1. Update `deployment.yaml` by adding `replica` setting

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-api-deployment
  namespace: test-uat
spec:
  replicas: 2 # <-- add replica setting
  selector:
    matchLabels:
      app: hello-api
  ...
```
2. Re-apply deployment.yaml and check the pods

```sh
$ kc apply -f deployment.yaml
deployment.apps/hello-api-deployment configured

$ kubectl get all --namespace test-uat
NAME                                        READY   STATUS    RESTARTS   AGE
pod/hello-api-deployment-85b75f6cc6-drmn5   1/1     Running   0          14m
pod/hello-api-deployment-85b75f6cc6-r7k5d   1/1     Running   0          8s
...
```

Create `service` which acts as a load balancer for hello-api pods. This results in other pods in the same K8s cluster be able to request this hello-api service.

1. Update `deployment.yaml` by adding `service` section
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-api-deployment
  namespace: test-uat
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hello-api
  template:
    metadata:
      labels:
        app: hello-api
    spec:
      containers:
      - name: hello-api
        image: azadpk8s2y0anxpaxdq1.azurecr.io/test/hello-api:1.0.0
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"

---

apiVersion: v1
kind: Service
metadata:
  name: hello-api-service
  namespace: test-uat
spec:
  selector:
    app: hello-api
  ports:
    - name: hello-api
      protocol: TCP
      port: 8082
      targetPort: 8080
```

2. Test the service by requesting to it
```sh
$ kubectl apply -f deployment.yaml
deployment.apps/hello-api-deployment unchanged
service/hello-api-service configured

$ kubectl exec -it hello-api-deployment-85b75f6cc6-drmn5 --namespace test-uat -- sh

/ # curl http://hello-api-service.test-uat.svc.cluster.local:8082
{"message": "こんにちは"}

# the service acts as a load balancer by forwarding the request to a pod and responses the message
# notes: K8s service has a service DNS as {service-name}.{namespace}.svc.cluster.local:{service port}
# ref: https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/
```

Expose API to the Internet (ex. https://voc-npr.azay.co.th/test/hello/ekkapob)
1. Update `deployment.yaml` by adding `ingress` setting.
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-api-deployment
  namespace: test-uat
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hello-api
  template:
    metadata:
      labels:
        app: hello-api
    spec:
      containers:
      - name: hello-api
        image: azadpk8s2y0anxpaxdq1.azurecr.io/test/hello-api:1.0.0
        imagePullPolicy: Always
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"

---

apiVersion: v1
kind: Service
metadata:
  name: hello-api-service
  namespace: test-uat
spec:
  selector:
    app: hello-api
  ports:
    - name: hello-api
      protocol: TCP
      port: 8082
      targetPort: 8080

---

apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "75"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "75"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "75"
    nginx.ingress.kubernetes.io/proxy-body-size: 10m
    nginx.ingress.kubernetes.io/configuration-snippet: |
      more_set_headers "Strict-Transport-Security: max-age=31536000; includeSubDomains";
      more_set_headers "Content-Security-Policy: default-src 'self'";
      more_set_headers "X-Frame-Options: DENY";
      more_set_headers "Permissions-Policy: fullscreen=(), geolocation=()";
      more_set_headers "Referrer-Policy: no-referrer";
    nginx.ingress.kubernetes.io/rewrite-target: /$2
  name: hello-api-ingress
  namespace: test-uat
spec:
  rules:
  - host: voc-npr.azay.co.th
    http:
      paths:
      - path: /test/hello/ekkapob(/|$)(.*)
        pathType: Prefix
        backend:
          service:
            name: hello-api-service
            port:
              number: 8082
```
2. Re-apply `deployment.yaml` and test the service

```
$ kubectl apply -f deployment.yaml
deployment.apps/hello-api-deployment unchanged
service/hello-api-service unchanged
ingress.networking.k8s.io/hello-api-ingress configured

$ curl https://voc-npr.azay.co.th/test/hello/ekkapob
{"message": "Ndewo"}
```

### Cleanup

```
# Delete all resources (pods, services, deployments, and etc.) but not the namespace
$ kubectl delete all --all -n test-uat
pod "hello-api-deployment-85b75f6cc6-drmn5" deleted
pod "hello-api-deployment-85b75f6cc6-r7k5d" deleted
service "hello-api-service" deleted
deployment.apps "hello-api-deployment" deleted

# Delete a namespace
$ kubectl delete namespace test-uat
namespace "test-uat" deleted
```

### Tips
```sh
# To set default namespace, therefore --namespace parameter is no longer needed
$ kubectl config set-context --current --namespace=test-uat
```
