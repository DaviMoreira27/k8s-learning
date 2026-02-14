A k8s manifest is normally a YAML file (can be JSON) that defines the desired state of an a kubernetes object (Pod, Service, Deployment) within a .


```yaml
kind: Pod
apiVersion: v1
metadata:
  name: foo-app
  labels:
    app: foo
spec:
  containers:
    - name: foo-app
      image: 'kicbase/echo-server:1.0'
```

In this example a pod is created named `foo-app` that uses the `kickbase/echo-server:1.0` image, this pod only answers information about the request. Besides that there a label field, were the name used by others objects to see it. Names identify the object individually, labels can be the same for various objects of the same type.

There are multiple types of K8s objects that be selected using the `kind` property, such as:

# 1. Workload Resources

## 1.1 Pod

**Description**  
The smallest deployable unit in Kubernetes. A Pod encapsulates one or more containers and represents a single instance of a running process.

**Lifecycle Characteristics**
- Created directly by the user or another controller
- No lifecycle management beyond container restart policy
- No replica management
- No rollout or rollback support
- No self-healing at the application level
- Ephemeral by design

A Pod does not manage its own lifecycle in a controlled way.  
If it crashes, it may restart depending on its restart policy, but if it is deleted, it is gone permanently unless recreated manually or by a controller.

**Use cases**
- Debugging
- Testing
- Temporary executions
- Short-lived experiments

**Rule**
Pods are disposable units and should not be used directly for production workloads that require controlled lifecycle management.

---

## 1.2 Deployment

**Description**  
A higher-level controller responsible for managing Pods in production environments. It runs inside the `kube-controller-manager` and provides controlled lifecycle management.

**Lifecycle Characteristics**
- Manages replica count
- Ensures desired state is maintained
- Automatically recreates failed Pods
- Supports rolling updates
- Supports rollbacks
- Enables declarative updates

Unlike a Pod, a Deployment actively controls the lifecycle of its Pods.  
It continuously compares the desired state with the current state and takes corrective action when needed.

**Responsibilities**
- Replica management
- Self-healing
- Rolling updates
- Rollbacks

**Use cases**
- REST APIs
- Web applications
- Backend services
- Workers
- Stateless microservices

**Rule**
If the application must run continuously with controlled updates and automatic recovery, use a Deployment.

---

## 1.3 StatefulSet

**Description**  
Manages stateful applications requiring persistent storage and stable network identity.

**Characteristics**
- Stable Pod names (e.g., app-0, app-1)
- One PersistentVolumeClaim per Pod
- Ordered startup and termination
- Stable network identity

**Use cases**
- Databases
- Redis (persistent mode)
- Kafka
- Elasticsearch

**Rule**
If the workload requires identity + persistent storage, use StatefulSet.

---

## 1.4 DaemonSet

**Description**  
Ensures that one Pod runs on each node (or selected nodes) in the cluster.

**Use cases**
- Monitoring agents (Node Exporter)
- Logging agents (Fluent Bit)
- Networking components (CNI)
- Storage drivers (CSI)

**Rule**
If a workload must run on every node, use a DaemonSet.

---

## 1.5 Job

**Description**  
Represents a finite task with a beginning and an end.

**Use cases**
- Database migrations
- Batch processing
- Maintenance tasks
- Data processing jobs

**Rule**
Runs until completion.

---

## 1.6 CronJob

**Description**  
A Job that runs on a schedule using cron syntax.

**Use cases**
- Daily backups
- Scheduled cleanups
- Periodic reports

**Rule**
Use CronJob for recurring scheduled tasks.

---

# 2. Networking Resources

## 2.1 Service

**Description**  
Provides a stable internal endpoint for accessing Pods and performs load balancing across Pods with matching labels.

**Why it exists**
Pods are ephemeral. Services provide a fixed virtual IP and DNS name.

**Types**
- ClusterIP (default)
- NodePort
- LoadBalancer
- Headless (clusterIP: None)

**Rule**
Pods change. Services remain stable.

---

## 2.2 Ingress

**Description**  
Exposes HTTP and HTTPS routes from the cluster to external clients.

**Responsibilities**
- Domain-based routing
- Path-based routing
- TLS termination

**Rule**
Use Ingress to manage external HTTP/HTTPS access.

---

# 3. Storage Resources

## 3.1 PersistentVolume (PV)

**Description**  
Represents actual storage provisioned in the cluster infrastructure.

**Use cases**
- Cloud disks (EBS, GCE PD, Azure Disk)
- NFS
- Local storage

**Provisioning**
- Typically created dynamically via StorageClass
- Can be manually created (static provisioning)

**Rule**
PV is the actual storage resource available in the cluster.

---

## 3.2 PersistentVolumeClaim (PVC)

**Description**  
A request for storage made by a Pod.

**Use cases**
- Database storage
- Redis persistent data
- Elasticsearch data

**Rule**
The application requests storage; the cluster provides it.

---

# 4. Scaling and Control

## 4.1 HorizontalPodAutoscaler (HPA)

**Description**  
Automatically scales the number of Pod replicas in a Deployment, StatefulSet, or ReplicaSet.

**Scaling metrics**
- CPU usage
- Memory usage
- Custom metrics (e.g., request rate)

**Example**
An API that scales automatically based on traffic.

**Rule**
HPA adjusts replica count based on load.

---

# 5. Configuration Management

## 5.1 ConfigMap

**Description**  
Stores non-sensitive configuration data in key-value format.

**Consumed as**
- Environment variables
- Command-line arguments
- Mounted configuration files

**Use cases**
- Application configuration
- Feature flags
- Runtime parameters

---

## 5.2 Secret

**Description**  
Stores sensitive information in key-value format.

**Use cases**
- Passwords
- API keys
- TLS certificates
- Tokens

**Important**
Secrets are base64-encoded by default. For stronger security, enable encryption at rest and configure RBAC properly.
