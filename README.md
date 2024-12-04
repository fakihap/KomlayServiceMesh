# KomlayServiceMesh
Service mesh orchestration for adapting service endpoints of Grand Oak Hospital and Pine Valley Hospital into a single endpoint using [istio](https://istio.io)🚢 on [minikube](https://minikube.sigs.k8s.io/)⚓

## Service Mesh Orchestration - Tutorial 🚀
adapted from
> https://kubernetes.io/docs/tutorials/hello-minikube/

> https://istio.io/latest/docs/setup/getting-started/

### 1. Preparation
on this tutorial, we are going to use a local image we manually build using local docker daemon (to simplify container registry and build simplicity)

#### Build the corresponding image for each service
1. Booking Service
```sh
cd booking-service/
docker build -t booking-service .
```
2. Grand Oak Hospital Service
```sh
cd grandoak-service/
docker build -t grand-oak .
```
3. Pine Valley Hospital Service
```sh
cd pinevalley-service/
docker build -t pine-valley .
```

#### Load the image into minikube registry
this step is required since `services.yaml` crd refer to this image to be able to create Deployment(s)
```sh
minikube image load booking-service
minikube image load grand-oak
minikube image load pine-valley
```

for other ways to do this, check here https://minikube.sigs.k8s.io/docs/handbook/pushing/

>tips: when you are trying to overwrite an already-loaded image, you can check if `minikube image load` succeed by comparing the `IMAGE ID` from `docker image ls` and `minikube image ls --format='table'`

> tips: if you are unable to do overwriting load for the correct IMAGE, you can check for running k8s deployments which may currently use the image. 

### 2. Minikube
refer to https://minikube.sigs.k8s.io/docs/start/?arch=%2Fwindows%2Fx86-64%2Fstable%2F.exe+download
1. Install Minikube on your platform
2. Start a new cluster
3. (optional) create an alias for kubectl (if you will be using minikube-builtin kubectl)

### 3. Download Istio
download istio and export istioctl to path
refer to https://istio.io/latest/docs/setup/getting-started/#download

### 4. Install Istio
refer to https://istio.io/latest/docs/setup/getting-started/#install
1. install istio using the no-gateway profile
    
    > alternatively, you can use the crd we have in this path "kube/istio-no-gateway.yaml"

2. add namespace label to enable automatic envoy sidecar injection
```sh
kubectl label namespace default istio-injection=enabled
```

### 5. Install the Kubernetes Gateway API CRDs
refer to https://istio.io/latest/docs/setup/getting-started/#gateway-api

### 6. Deploy the app
```sh
kubectl apply -f kube/services.yaml    
```

### 7. Open the app to outside traffic
1. Create Kubernetes Gateway
```sh
kubectl apply -f kube/gateways.yaml    
```

2. Change the service type to ClusterIP by annotating the gateway
```sh
kubectl annotate gateway komlay-gateway networking.istio.io/service-type=ClusterIP --namespace=default
 ```

### 8. Forward the Port
```sh
kubectl port-forward svc/komlay-gateway-istio 8080:80   
```
now, you can access the exposed API from `localhost:8080`