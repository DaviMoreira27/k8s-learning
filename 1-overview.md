## Kubernetes Overview

**What Kubernetes is**
Kubernetes is an open-source, portable, extensible platform for managing containerized workloads and services with declarative configuration and automation. It has a large ecosystem of tools and support.

**Why Kubernetes matters**
In production, you need to manage containers so applications stay running, scale with demand, and recover from failures. Kubernetes provides a framework for resilient distributed systems, handling scaling, failover, deployments, and more.

**Key capabilities**

* Service discovery and load balancing
* Storage orchestration
* Automated rollouts and rollbacks
* Efficient resource scheduling (bin packing)
* Self-healing (restart/replace failing containers)
* Config and secret management
* Support for batch workloads
* Horizontal scaling
* Dual IPv4/IPv6 support
* Extensible architecture

## Kubernetes vs Traditional Deployment Models

### Physical Servers
In a physical server model, applications run directly on dedicated hardware. Each server is usually configured manually and tied to specific workloads.

**Key differences**
- Low resource utilization: hardware is often underused.
- Scaling is slow and manual (buying and installing new servers).
- Failure recovery is mostly manual.
- Tight coupling between hardware and application.

Kubernetes abstracts the hardware completely. Applications are scheduled dynamically across a cluster, making better use of resources and enabling automatic recovery and scaling.

---

### Hypervisors and Virtual Machines
With hypervisors (VMware, KVM, Hyper-V), applications run inside virtual machines, each with its own operating system.

**Key differences**
- Better isolation than physical servers, but heavier than containers.
- Slower startup times due to full OS per VM.
- Scaling usually requires VM provisioning and configuration.
- Resource overhead is higher.

Kubernetes typically runs containers instead of full VMs. Containers share the host OS, start faster, and use fewer resources. Kubernetes also automates scheduling, scaling, and self-healing at a higher level than most VM-based setups.

---

## Kubernetes vs Docker

### Docker (Standalone)
Docker focuses on building and running containers on a single host.

**Key differences**
- Docker runs containers; it does not manage clusters by itself.
- No native high availability across multiple machines.
- Scaling and recovery are manual or script-based.
- Limited service discovery and load balancing.

Kubernetes uses Docker (or other container runtimes) as a building block but adds orchestration: multi-node scheduling, automatic restarts, scaling, networking, and configuration management.

---

## Kubernetes vs Docker Swarm

### Docker Swarm
Docker Swarm is Docker’s native clustering and orchestration solution.

**Key differences**
- Easier to set up and simpler to understand.
- Limited feature set compared to Kubernetes.
- Smaller ecosystem and community adoption.
- Less flexibility and extensibility.

Kubernetes is more complex but more powerful. It offers richer scheduling, stronger self-healing, advanced networking, extensibility via controllers and CRDs, and is the industry standard for large-scale container orchestration.
