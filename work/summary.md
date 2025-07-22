# Summary

## Problem

Within Harvester, the current in-cluster communication between virtual machines
and pods are unencrypted. This exposed user workloads to spoofing and
man-in-the-middle attack.

## Solution

Research into how service mesh solutions like Linkerd and Istio enabling
service-to-service automatic mTLS. The mTLS mechanism should allow us to achieve
a zero-trust networking communication by enabling in-cluster traffic encryption
and service-to-service identification. We need to figure out options to manage
the x509 certificates without too much disruption to running services.
