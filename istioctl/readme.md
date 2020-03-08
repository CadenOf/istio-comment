## ./istioctl --help
```
# ./istioctl --help
Istio configuration command line utility for service operators to
debug and diagnose their Istio mesh.

Usage:
  istioctl [command]

Available Commands:
  authn           Interact with Istio authentication policies
  authz           (authz is experimental.  Use `istioctl experimental authz`)
  convert-ingress Convert Ingress configuration into Istio VirtualService configuration
  dashboard       Access to Istio web UIs
  deregister      De-registers a service instance
  experimental    Experimental commands that may be modified or deprecated
  help            Help about any command
  kube-inject     Inject Envoy sidecar into Kubernetes pod resources
  manifest        Commands related to Istio manifests
  profile         Commands related to Istio configuration profiles
  proxy-config    Retrieve information about proxy configuration from Envoy [kube only]
  proxy-status    Retrieves the synchronization status of each Envoy in the mesh [kube only]
  register        Registers a service instance (e.g. VM) joining the mesh
  validate        Validate Istio policy and rules
  verify-install  Verifies Istio Installation Status or performs pre-check for the cluster before Istio installation
  version         Prints out build version information

Flags:
      --context string            The name of the kubeconfig context to use
  -h, --help                      help for istioctl
  -i, --istioNamespace string     Istio system namespace (default "istio-system")
  -c, --kubeconfig string         Kubernetes configuration file
      --log_output_level string   Comma-separated minimum per-scope logging level of messages to output, in the form of <scope>:<level>,<scope>:<level>,... where scope can be one of [ads, all, analysis, attributes, authn, cacheLog, citadelClientLog, configMapController, conversions, default, googleCAClientLog, grpcAdapter, kube, kube-converter, mcp, meshconfig, model, patch, processing, rbac, resource, runtime, sdsServiceLog, secretFetcherLog, source, stsClientLog, tpath, translator, util, validation, vaultClientLog] and level can be one of [debug, info, warn, error, fatal, none] (default "default:info,validation:error,processing:error,source:error,analysis:warn")
  -n, --namespace string          Config namespace

Use "istioctl [command] --help" for more information about a command.
```

## istioctl kube-inject --help
```
# ./istioctl kube-inject --help

kube-inject manually injects the Envoy sidecar into Kubernetes
workloads. Unsupported resources are left unmodified so it is safe to
run kube-inject over a single file that contains multiple Service,
ConfigMap, Deployment, etc. definitions for a complex application. It's
best to do this when the resource is initially created.

k8s.io/docs/concepts/workloads/pods/pod-overview/#pod-templates is
updated for Job, DaemonSet, ReplicaSet, Pod and Deployment YAML resource
documents. Support for additional pod-based resource types can be
added as necessary.

The Istio project is continually evolving so the Istio sidecar
configuration may change unannounced. When in doubt re-run istioctl
kube-inject on deployments to get the most up-to-date changes.

Usage:
  istioctl kube-inject [flags]

Examples:

# Update resources on the fly before applying.
kubectl apply -f <(istioctl kube-inject -f <resource.yaml>)

# Create a persistent version of the deployment with Envoy sidecar
# injected.
istioctl kube-inject -f deployment.yaml -o deployment-injected.yaml

# Update an existing deployment.
kubectl get deployment -o yaml | istioctl kube-inject -f - | kubectl apply -f -

# Capture cluster configuration for later use with kube-inject
kubectl -n istio-system get cm istio-sidecar-injector  -o jsonpath="{.data.config}" > /tmp/inj-template.tmpl
kubectl -n istio-system get cm istio -o jsonpath="{.data.mesh}" > /tmp/mesh.yaml
kubectl -n istio-system get cm istio-sidecar-injector -o jsonpath="{.data.values}" > /tmp/values.json
# Use kube-inject based on captured configuration
istioctl kube-inject -f samples/bookinfo/platform/kube/bookinfo.yaml \
	--injectConfigFile /tmp/inj-template.tmpl \
	--meshConfigFile /tmp/mesh.yaml \
	--valuesFile /tmp/values.json


Flags:
  -f, --filename string              Input Kubernetes resource filename
  -h, --help                         help for kube-inject
      --injectConfigFile string      injection configuration filename. Cannot be used with --injectConfigMapName
      --injectConfigMapName string   ConfigMap name for Istio sidecar injection, key should be "config". (default "istio-sidecar-injector")
      --meshConfigFile string        mesh configuration filename. Takes precedence over --meshConfigMapName if set
      --meshConfigMapName string     ConfigMap name for Istio mesh configuration, key should be "mesh" (default "istio")
  -o, --output string                Modified output Kubernetes resource filename
      --valuesFile string            injection values configuration filename.

Global Flags:
      --context string            The name of the kubeconfig context to use
  -i, --istioNamespace string     Istio system namespace (default "istio-system")
  -c, --kubeconfig string         Kubernetes configuration file
      --log_output_level string   Comma-separated minimum per-scope logging level of messages to output, in the form of <scope>:<level>,<scope>:<level>,... where scope can be one of [ads, all, analysis, attributes, authn, cacheLog, citadelClientLog, configMapController, conversions, default, googleCAClientLog, grpcAdapter, kube, kube-converter, mcp, meshconfig, model, patch, processing, rbac, resource, runtime, sdsServiceLog, secretFetcherLog, source, stsClientLog, tpath, translator, util, validation, vaultClientLog] and level can be one of [debug, info, warn, error, fatal, none] (default "default:info,validation:error,processing:error,source:error,analysis:warn")
  -n, --namespace string          Config namespace
```