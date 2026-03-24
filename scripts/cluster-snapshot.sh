#!/bin/bash
echo "--- CLUSTER OVERVIEW ---"
kubectl cluster-info

echo -e "\n--- NODE STATUS ---"
kubectl get nodes -o custom-columns=NAME:.metadata.nam,STATUS:.status.conditions[-1].type,VERSION:.status.nodeInfo.kubeletVersion,CPU:.status.capacity.cpu

echo -e "\n--- TOP PODS (RESOURCE USAGE) ---"
kubectl top pods -A --sort-by=cpu | head -n 10

echo -e "\n--- RECENT EVENTS (ERRORS) ---"
kubectl get events -A --field-selector type!=Normal --sort-by='.lastTimestamp' | tail -n 5
