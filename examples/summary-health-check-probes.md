# Summary

## Problem

Harvester relies on KubeVirt to provision and manage QEMU-based virtual machines.
KubeVirt run these virtual machines inside the `virt-launcher` pods, which are
initialized with the `virt-launcher-monitor` process. This process serves as the
init process (with PID 1) of the container, and is used for launching and
managing the lifecycle of the VMs.

By default, the readiness of the `virt-launcher` pod is determined by the
readiness of the `virt-launcher` container's `virt-launcher-monitor` process.
This is obviously not the same as the readiness of other user processes running
inside the pod.

This can lead to boot time or runtime failure scenarios where the `virt-launcher`
container may mistakenly report itself to be running and ready, but the actual
user processes responsible for serving traffic are not yet ready.

As a result, downstream client services may start sending traffic to the VM
before it is actually ready to handle it. Health check services would also be
misled to believe that the VM is healthy and hence, not initiate any necessary
remediation or graceful termination actions on the VM.

This out-of-sync readiness status can also affect the upgrade process of the
cluster. The upgrade controller may mistakenly believe that all the VMs on a node
are ready and healthy, and hence, proceed to cordoning and draining the next
node. This can cause service disruption and quorum loss to VM groups.

We want to provide a way for service owners to define the readiness of their VMs
in a more accurate way, so that the cluster can make better decisions on when to
send traffic to the VMs and when to initiate remediation actions.

### User Scenario 1

A VM running Nginx web server with an expired TLS certificate continues to appear
as running and healthy. When the traffic is routed to the VM, users would
encounter TLS errors and be unable to access the web server.

### User Scenario 2

In Harvester/Rancher integration setup, Rancher user can provision downstream
guest clusters to have their nodes run on Harvester-managed VMs. If the Kubelet
on the guest cluster node is down, the VM may remain in running and healthy state
on the Harvester side. When the Harvester CSI driver on the guest cluster
initializes an non-graceful termination on the stateful workloads running on the
unhealthy node, it could lead to storage and compute synchronization issues on
the Harvester side, because the Longhorn controllers would not have performed the
necessary remediation actions on the hotplug volumes used by the stateful
workloads.

### User Scenario 3

During cluster upgrade, a PgSQL cluster may have 3 replicas running on 3
different VMs. If the upgrade controller mistakenly believes that a replica on
the upgrading VMs is ready and healthy, it may proceed to cordoning and draining
the next VM. This can cause service disruption and quorum loss to the PgSQL
cluster, leading to downtime and potential data loss.

## Solution

### KubeVirt Liveness and Readiness Probes

Provide service owners with the ability to define custom liveness and readiness
probes for their VMs, which can be used to determine the actual readiness of the
user processes running inside the VMs.

In short, a probe is a diagnostic performed periodically by the kubelet on a
container. To perform a diagnostic, the kubelet either executes code within the
container or makes a network request.

KubeVirt already supports custom readiness probes for VMs, so we can make this
feature available to Harvester users by exposing it through the Harvester UI. No
changes are needed on the Harvester API and the controller backend.

The UI should be organized in such a way to support these 4 probe mechanisms:

Probe Mechanism | Description
--------------- | -----------
`exec`            | Executes a specified command inside the container. The diagnostic is considered successful if the command exits with a status code of 0.
`grpc`            | Performs a remote procedure call using gRPC. The target should implement gRPC health checks. The diagnostic is considered successful if the status of the response is SERVING. For more details, see gRPC probes.
`httpGet`         | Performs an HTTP GET request against the Pod's IP address on a specified port and path. The diagnostic is considered successful if the response has a status code greater than or equal to 200 and less than 400. For more details, see HTTP probes.
`tcpSocket`       | Performs a TCP check against the Pod's IP address on a specified port. The diagnostic is considered successful if the port is open. If the remote system (the container) closes the connection immediately after it opens, this counts as healthy. For more details, see TCP probes.

In addition, we also want to provide a way for service owners to configure the
probes' check behaviour and frequency by exposing the following parameters:

* `initialDelaySeconds`
* `periodSeconds`
* `timeoutSeconds`
* `successThreshold`
* `failureThreshold`
* `terminationGracePeriodSeconds`

#### Limitations on Probes

Note that KubeVirt doesn't support startup probes for VMs, so we won't be able to
provide this feature to Harvester users until KubeVirt implements it.

Incorrect implementation of readiness probes may result in an ever growing number
of processes in the container, and resource starvation if this is left unchecked.

Incorrect implementation of liveness probes can lead to cascading failures. This
results in restarting of container under high load; failed client requests as
your application became less scalable; and increased workload on remaining pods
due to some failed pods.

Unlike the other mechanisms, exec probe's implementation involves the creation/
forking of multiple processes each time when executed. As a result, in case of
the clusters having higher pod densities, lower intervals of initialDelaySeconds,
periodSeconds, configuring any probe with exec mechanism might introduce an
overhead on the cpu usage of the node.

Using the `exec` probe requires the qemu-guest-agent package to be installed in
the VM. A command supplied to an exec probe will be wrapped by virt-probe in the
operator and forwarded to the guest.

#### Dual Stack Considerations

In dual stack clusters, the HTTP and TCP probes only works with IPv4 addresses
because the `host` field of the requests are default to the pod's IPv4 address.

### Guest OS Watchdog

Watchdogs focus on ensuring that an Operating System is still responsive. They
complement the probes which are more workload centric. Watchdogs require kernel
support from the guest and additional tooling like the commonly used `watchdog`
binary.

An example of a KubeVirt `VirtualMachine` spec with a guest OS watchdog configured:

```yaml
apiVersion: kubevirt.io/v1
kind: VirtualMachine
metadata:
  name: example-watchdog
spec:
  runStrategy: Always
  template:
    spec:
      domain:
        devices:
          watchdog:
            name: mywatchdog
            i6300esb:
              action: "poweroff"
          disks:
          - disk:
              bus: virtio
            name: containerdisk
          - disk:
              bus: virtio
            name: cloudinit
          rng: {}
        resources:
          requests:
            memory: 1Gi
      terminationGracePeriodSeconds: 0
      volumes:
      - name: containerdisk
        containerDisk:
          image: quay.io/containerdisks/fedora:latest
      - name: cloudinit
        cloudInitNoCloud:
          userData: |-
            #cloud-config
            password: fedora
            user: fedora
            chpasswd: { expire: False }
            packages:
                busybox
```

The VM in this example will have the device exposed as `/dev/watchdog`. This
device can then be used by the watchdog binary. For example, if root executes
this command inside the VM:

```sh
sudo busybox watchdog -t 2000ms -T 4000ms /dev/watchdog
```

The watchdog will send a heartbeat every two seconds to `/dev/watchdog` and after
four seconds without a heartbeat the defined action will be executed. In this
case a hard `poweroff`.

Possible `action` for watchdog, which defines what will happen if the OS can't
respond anymore, include:

* `poweroff`: Power off the VM.
* `reset`: Reset the VM.
* `shutdown`: Shutdown the VM.

### Pod Disruption Budget

Kubernetes provides a mechanism called Pod Disruption Budget (PDB) to allow
cluster administrators and service owners to plan for _voluntary service
disruptions_, by limiting the number of pods of a replicated application that are
down simultaneously.

Examples of this kind of disruptions include:

* Draining a node for repair or upgrade
* Draining a node from a cluster to scale the cluster down
* Removing a pod from a node to permit something else to fit on that node
* Deleting the deployment or other controller that manages the pod
* Updating a deployment's pod template causing a restart

Kubernetes accounts for the constraints defined in the PDB when controllers or
tools initiate pod eviction via the Eviction API. The `kubectl drain` command is
an example of such API-initiated eviction. If evicting the pod would cause the
number of available replicas to fall below the minimum specified in the PDB, the
eviction will be blocked until it is safe to proceed.

If liveness and readiness probes are configured for the VMs, the PDB can prevent
eviction of other VM replicas are blocked if they are not ready.

When a KubeVirt VM is configured with its `evictionStrategy` set to `LiveMigrate`
or `LiveMigrateIfPossible`, KubeVirt automatically creates and adjusts PDB for the
VM when it is evicted, to allow for live migration.

Hence, we might not need to expose any PDB configuration to Harvester users.
We should investigate into how KubeVirt works with user-defined PDB for the
`virt-launcher` pods.

If we decide to provide a way for Harvester users to define custom PDB, the UI should be organized in such a way to support the following parameters:

* `minAvailable` which defines the minimum number of pods that must be available
after the eviction. This must be an absolute number, not a percentage.
* `unhealthyPodEvictionPolicy` which support only the `IfHealthyBudget` and
`AlwaysAllow` policies

An example of a PDB spec that selects all the `virt-launcher` pods that backs the
VMs of a `demo-green` guest clusters:

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: guest-rke2-pdb
  namespace: demo-green
spec:
  minAvailable: 2
  selector:
    matchLabels:
      guestcluster.harvesterhci.io/name: guest-rke2
```

#### Limitations on PDB

PDB is more commonly used with pods managed by higher level controllers like
Deployments, StatefulSets, etc. While using it with "bare pods" is possible,
there are some limitations and considerations to keep in mind:

* only .spec.minAvailable can be used, not .spec.maxUnavailable.
* only an integer value can be used with .spec.minAvailable, not a percentage.

### Harvester/Rancher Integration

In setup the involves integrating Harvester with Rancher to manage downstream
guest clusters, further investigation is needed to see if the Rancher UI already
supports defining custom readiness probes for generic Kubernetes pods. If yes, the
enhancement to make the probes available to VMs shouldn't be too difficult.
Otherwise, more work is needed to implement this feature in Rancher before we can
expose it to Harvester users.

## References

### Harvester

* GH issue describing original problem - <https://github.com/harvester/harvester/issues/9523>>
* Harvester/Rancher integration documentation - <https://docs.harvesterhci.io/v1.8/rancher/rancher-integration>

### Concepts

* Liveness/readiness probes concept in Kubernetes: <https://kubernetes.io/docs/concepts/workloads/pods/probes/>
* KubeVirt support of liveness/readiness probes and watchdog: <https://kubevirt.io/user-guide/user_workloads/liveness_and_readiness_probes>
* Pod disruption budget concept in Kubernetes: <https://kubernetes.io/docs/concepts/workloads/pods/disruptions/>
