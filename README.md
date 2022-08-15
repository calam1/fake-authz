A fake external auth for envoy proxy
https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter



# create a namespace and enable istio-injection if it doesn't exist
❯❯❯ kubectl create ns mystuff
namespace/mystuff created

❯❯❯ kubectl label namespace mystuff istio-injection=enabled --overwrite=true

namespace/mystuff labeled

❯❯❯ kubectl get namespace -L istio-injection

NAME                 STATUS   AGE    ISTIO-INJECTION
default              Active   5d5h
istio-system         Active   5d5h   disabled
kube-node-lease      Active   5d5h
kube-public          Active   5d5h
kube-system          Active   5d5h
local-path-storage   Active   5d5h
mystuff              Active   29s    enabled


# go to the fake-authz directory and build the fake-authz docker image
❯❯❯ docker build -t fake-authz .


# I use Kind for local k8, so copy the image over to the cluster (nodeport is the name of my cluster)
❯❯❯ kind load docker-image fake-authz:latest --name nodeport
Image: "fake-authz:latest" with ID "sha256:e6cb6d3b8eb99ec342248af5e42a022022299cc4bc8ac1a10147ff2ec1077b1f" not yet present on node "nodeport-control-plane", loading...
Image: "fake-authz:latest" with ID "sha256:e6cb6d3b8eb99ec342248af5e42a022022299cc4bc8ac1a10147ff2ec1077b1f" not yet present on node "nodeport-worker", loading...


# go to the deploymnent directory of the fake-authz project and deploy
❯❯❯ kubectl apply -f fake-authz.yml -n mystuff
service/authz created
deployment.apps/authz created
horizontalpodautoscaler.autoscaling/authz created
poddisruptionbudget.policy/authz created


# test it out, port forward
❯❯❯ kubectl port-forward deployment/authz 50051:50051 -n mystuff
Forwarding from 127.0.0.1:50051 -> 50051
Forwarding from [::1]:50051 -> 50051


# install grpcurl with brew and run the following to make sure it works and is deployed correctly
❯❯❯ grpcurl -d  '{ "attributes": { "request": { "http": { "method": "GET", "headers": {"x-api-key":"123abc"} } } } }' --plaintext  localhost:50051 envoy.service.auth.v3.Authorization/Check
{
  "status": {
    "message": "OK"
  },
  "okResponse": {
    "headers": [
      {
        "header": {
          "key": "x-test-api-header-value",
          "value": "123abc"
        }
      },
      {
        "header": {
          "key": "x-ext-auth-ratelimit",
          "value": "5"
        }
      },
      {
        "header": {
          "key": "x-ext-auth-ratelimit-unit",
          "value": "MINUTE"
        }
      }
    ]
  }
}