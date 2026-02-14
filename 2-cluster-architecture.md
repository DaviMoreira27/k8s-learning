# Kubernetes Architecture

## Control Plane — What it actually does and how

The control plane does **three things only**:
1. Stores desired and current state.
2. Runs control loops that compare desired vs actual state.
3. Issues instructions to worker nodes.

It does **not** run containers and does **not** execute workloads.

---

## kube-apiserver — State entry point and coordination hub

The API server is **the only component that reads or writes cluster state**.

### What it does
- Exposes the Kubernetes API (REST).
- Authenticates and authorizes every request.
- Validates object schemas.
- Persists objects to etcd.
- Serves watches to other components.

### How it works
- Receives a request (e.g., create Pod, update Deployment).
- Runs authn/authz (certs, tokens, RBAC).
- Validates object against OpenAPI schema.
- Writes the object to etcd.
- Emits events and allows watchers to observe the change.

All components (scheduler, controllers, kubelet) **watch the API**, not each other.

There is no direct component-to-component communication without the API server.

---

## etcd — Cluster state storage

etcd stores **only data**, no logic.

### What it stores
- Desired state: Deployments, ReplicaSets, PodSpecs.
- Current state: PodStatus, NodeStatus.
- Metadata: labels, annotations.
- Secrets and ConfigMaps.

### How it is used
- API server is the single writer.
- All reads/watches go through the API server.
- Uses Raft for strong consistency.
- Any state change = etcd write.

If etcd is unavailable, the cluster becomes read-only or unusable.

---

## kube-controller-manager — Reconciliation engine

Controllers are **control loops**. They do not execute actions directly on nodes.

### What controllers do
Each controller:
- Watches a specific resource type.
- Compares desired state vs observed state.
- Writes corrective actions back to the API.

Example: ReplicaSet controller  
Desired: replicas = 3  
Observed: 2 Pods running  
Action: create 1 Pod object via API

### How controllers work
Loop pattern:
1. Watch resource changes via API server.
2. Read desired state.
3. Query current state (also via API).
4. If mismatch exists, write new objects or updates.
5. Loop repeats continuously.

Controllers never talk to kubelet or nodes directly.

**Its not responsible for the pods creation directly, it will only send a request to the apiserver for them to do it.**

---

## kube-scheduler — Pod placement decision maker

Scheduler is **purely a decision engine**.

### What it does
- Assigns a node to a Pod by setting `spec.nodeName`.

### How it works
Triggered when:
- A Pod exists with no `nodeName`.

Steps:
1. Watch for unscheduled Pods.
2. Build list of candidate nodes.
3. Run filtering phase:
   - CPU/memory availability
   - nodeSelector / affinity
   - taints and tolerations
4. Run scoring phase:
   - Resource balance
   - Topology preferences
5. Select highest-score node.
6. Write binding decision to API server.

Scheduler does not start Pods.
Scheduler does not monitor runtime health.

---

## cloud-controller-manager — Cloud-side reconciliation

Only exists in cloud environments.

### What it does
- Translates Kubernetes objects into cloud provider resources.

### Examples
- Service of type LoadBalancer → create cloud LB.
- Node object → map to cloud VM lifecycle.
- Route management for Pod networking.

### How it works
Same controller pattern:
- Watch API objects.
- Call cloud provider APIs.
- Write results back to API server.

---

## Worker Nodes — Execution only

Worker nodes **never make decisions**.

---

## kubelet — Node-level executor

kubelet enforces PodSpecs assigned to its node.

### What it does
- Watches Pods where `spec.nodeName == this node`.
- Talks to container runtime via CRI.
- Starts/stops containers.
- Runs health checks.
- Reports status.

### How it works
1. Watch assigned Pods via API server.
2. For each Pod:
   - Pull images if needed.
   - Create containers via CRI.
   - Mount volumes.
3. Periodically report PodStatus and NodeStatus.
4. Restart containers based on restart policy.

kubelet never schedules Pods and never decides desired state.

---

## kube-proxy — Service traffic implementation

kube-proxy is **network rule programming**, nothing else.

### What it does
- Translates Service definitions into dataplane rules.

### How it works
- Watches Services and Endpoints.
- Programs iptables / IPVS / eBPF rules.
- Ensures traffic to a Service IP is forwarded to Pod IPs.

No awareness of applications. No load-balancing logic beyond rule setup.

---

## Container Runtime — Container execution

### What it does
- Pull images.
- Create containers.
- Manage namespaces and cgroups.

### How it is used
- kubelet talks to runtime via CRI.
- Runtime does not know Kubernetes concepts.

---

## End-to-End Control Flow (Concrete)

Example: `replicas = 3`

1. Deployment object written to API server.
2. Stored in etcd.
3. ReplicaSet controller detects mismatch.
4. Controller creates 3 Pod objects.
5. Scheduler assigns each Pod to a node.
6. kubelet on each node sees assigned Pod.
7. kubelet starts containers via runtime.
8. kubelet reports status.
9. Controllers observe status and stop acting.

This loop runs continuously.

---

## Key Architectural Fact

Kubernetes is **not event-driven execution**.
It is **state-driven reconciliation**.

- Control Plane: writes desired state and reconciliation actions.
- Worker Nodes: execute instructions and report status.

Nothing else.
