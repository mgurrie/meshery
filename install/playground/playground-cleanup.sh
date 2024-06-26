#!/bin/bash
#TODO: Make the cleanup smarter than basic namespace deletion
#Script is placed in /root directory in VM.
for ns in $(kubectl get ns -o jsonpath="{.items[*].metadata.name}");
do
    valid_namespaces=("kube-system" "default" "monitoring" "kube-flannel" "kube-node-lease" "kube-public" "meshery" "metallb-system" "projectcontour" "ingress-nginx" "layer5-cloud" "postgres")
    if [[ ! " ${valid_namespaces[*]} " =~ [[:space:]]${ns}[[:space:]] ]]; then
          echo "Deleting namespace $ns"
          kubectl delete ns "$ns"
    fi
done

# remove mutatingwebhookconfigurations
for mwh in $(kubectl get mutatingwebhookconfigurations -o jsonpath="{.items[*].metadata.name}");
do
    if [[ "$mwh" == "consul-consul-connect-injector" ]];then
          echo "Deleting mutatingwebhookconfigurations $mwh"
          kubectl delete mutatingwebhookconfigurations "$mwh"
    fi
done
