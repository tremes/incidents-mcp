THANOS_POD=$(oc get -n openshift-monitoring pod -l app.kubernetes.io/instance=thanos-querier --no-headers -o custom-columns=":metadata.name" | head -n 1)

oc port-forward pod/$THANOS_POD 8080:9090