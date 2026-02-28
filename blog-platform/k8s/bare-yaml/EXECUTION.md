
```sh
cd ~/projects/k8s-learning/blog-platform/k8s/bare-yaml
```

```sh
kubectl apply -f database/Secrets.yaml
kubectl apply -f database/Service.yaml
kubectl apply -f database/StatefulSet.yaml
kubectl apply -f cache/Deployment.yaml
kubectl apply -f cache/Service.yaml
kubectl apply -f backend/Deployment.yaml
kubectl apply -f backend/Service.yaml
kubectl apply -f worker/Deployment.yaml
kubectl apply -f worker/Service.yaml
kubectl apply -f frontend/Deployment.yaml
kubectl apply -f frontend/Service.yaml
```

To check if the database is ready call:

```sh
kubectl rollout status statefulset/postgres
```

To access, run:
```sh
minikube service frontend-service --url
```
or access directly via `<node-ip>:30080` (e.g. `http://192.168.49.2:30080`).

**Note:** `minikube service --url` is only required on macOS and Windows, where Docker runs inside a VM (HyperKit/WSL2), making the node IP unreachable from the host. On Linux, Docker uses the host network stack directly via a bridge, so the node IP is routable without any tunnel.
