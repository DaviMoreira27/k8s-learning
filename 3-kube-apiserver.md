# kube-apiserver

The kube-apiserver is a **stateless HTTP API process** that sits in front of etcd.
It does **not coordinate execution**, **does not run logic**, and **does not push commands**

Its job is to **control access to cluster state** and **serialize all state changes**.

---

## What “only component that reads or writes cluster state” means

### Writes
- The kube-apiserver is the **only process allowed to write to etcd**.
- Controllers, scheduler, kubelet:
  - DO NOT write to etcd
  - DO NOT share memory
  - DO NOT communicate directly
- They all perform HTTP requests to the API server.

This guarantees:
- Ordered writes
- Consistent state
- No race conditions between components

---

### Reads
- All components read state **through the API server**
- The API server reads from etcd and returns objects
- Watch streams are also served by the API server

No component:
- Reads etcd directly
- Has privileged storage access
- Owns authoritative state locally

---

## Request lifecycle (exact steps)

Example: ReplicaSet controller creating a Pod

1. Controller sends HTTP POST to kube-apiserver
2. API server performs:
   - Authentication (certs / tokens)
   - Authorization (RBAC)
   - Admission control (mutating + validating webhooks)
   - Schema validation (OpenAPI)
3. API server writes object to etcd
4. API server returns success to controller
5. API server emits watch events to all subscribers

At no point does the controller talk to:
- Scheduler
- kubelet
- any node

---

## Watch mechanism

The kube-apiserver maintains **long-lived HTTP connections** (watches).

Components:
- Register a watch on resource types
- Receive change events in order
- Update their local cache

Example:
- Scheduler watches Pods without `spec.nodeName`
- kubelet watches Pods with `spec.nodeName == node`
- Controllers watch the resources they own

The API server does **not decide who should react**.
It only emits state changes.

---

## Informers and caches

Controllers and scheduler use **local caches** populated by watches.

Important:
- Cache is an optimization
- Cache is not authoritative
- Any write must still go through the API server

If cache is stale:
- API server rejects invalid writes
- ResourceVersion prevents conflicts

---

## Why this design exists

This design guarantees:
- Loose coupling between components
- Horizontal scalability of control plane
- Deterministic behavior under failure
- Restart-safe components (no in-memory state loss)

If a controller crashes:
- State is still in etcd
- Watch reconnects
- Reconciliation resumes
