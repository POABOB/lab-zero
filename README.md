# lab-zero

This is a two microserives of Users and Orders.

## how to init a microservice of go-zero
### init a rpc service

1. add a structure of go-zero rpc sample.

```bash
goctl rpc new demo
```

2. edit the protobuffer file of demo, and generate a new structure of real service.

```bash
cd demo
goctl rpc protoc --go_out=. --go-grpc_out=. --zrpc_out=. demo.proto
```

3. edit the bussiness logic in `demo/internal/logic` of your microservice.

4. run it.

```bash
go run demo.go
```

### init a gateway of rpc service

1. generate a proto descriptor of microservice

```bash
protoc --include_imports --proto_path=. --descriptor_set_out=demo.pb demo.proto
```

2. edit `demo/internal/config/config.go` in order to add the config to struct.

```go
package config

import (
	"github.com/zeromicro/go-zero/gateway" // add this line
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Gateway          gateway.GatewayConf // add this line
}
```

3. edit `demo/etc/demo.yaml` in order to add gateway config.

```yaml
# add this
Gateway:
  Name: gateway
  Port: 80
  Upstreams:
    - Grpc:
        Target: 0.0.0.0:8888
      ProtoSet: demo.pb
```

4. edit `demo/demo.go` in order to start two service. One is rpc service, and the other is gateway service.

```go
func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

    // init a gateway server
	gw := gateway.MustNewServer(c.Gateway)

    // init a service group
	group := service.NewServiceGroup()

    // init a demo server
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		demo.RegisterDemoServer(grpcServer, server.NewDemoServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
    // add demo server to service group
	group.Add(s)

    // add gateway server to service group 
	group.Add(gw)
	defer group.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Starting gateway at %s:%d...\n", c.Gateway.Host, c.Gateway.Port)
    // service group start
	group.Start()
}
```


5. run it.

```bash
go run demo.go
```

### add validate to rpc service

1. add rule to the `demo.proto`.

> protoc-gen-validate
 doc: https://github.com/bufbuild/protoc-gen-validate/blob/main/docs.md

```proto
message CreateRequest {
  string name = 1 [(validate.rules).string = {
    max_bytes: 256,
  }];
}
```

2. generate the validate file(`demo.pb.validate.go`).
```bash
protoc --validate_out=lang=go:./ *.proto
```

3.1. add to the logic in `demo/internal/logic` mannually .

```go
func (l *CreateLogic) Create(in *demo.CreateRequest) (*demo.CreateResponse, error) {
	if err := in.ValidateAll(); err != nil {
		return &user.CreateResponse{
			Msg: err.Error(),
		}, nil
	}

	return &demo.CreateResponse{
		Msg: "Ok",
	}, nil
}
```

3.2. use middleware to validate it.

- add `demo/internal/middleware/validate/validate.go`

```go
// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package validate

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// The validate interface starting with protoc-gen-validate v0.6.0.
// See https://github.com/envoyproxy/protoc-gen-validate/pull/455.
type validator interface {
	Validate(all bool) error
}

// The validate interface prior to protoc-gen-validate v0.6.0.
type validatorLegacy interface {
	Validate() error
}

func validate(req interface{}) error {
	switch v := req.(type) {
	case validatorLegacy:
		if err := v.Validate(); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
	case validator:
		if err := v.Validate(false); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return nil
}

// UnaryServerInterceptor returns a new unary server interceptor that validates incoming messages.
//
// Invalid messages will be rejected with `InvalidArgument` before reaching any userspace handlers.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// UnaryClientInterceptor returns a new unary client interceptor that validates outgoing messages.
//
// Invalid messages will be rejected with `InvalidArgument` before sending the request to server.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if err := validate(req); err != nil {
			return err
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that validates incoming messages.
//
// The stage at which invalid messages will be rejected with `InvalidArgument` varies based on the
// type of the RPC. For `ServerStream` (1:m) requests, it will happen before reaching any userspace
// handlers. For `ClientStream` (n:1) or `BidiStream` (n:m) RPCs, the messages will be rejected on
// calls to `stream.Recv()`.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := &recvWrapper{stream}
		return handler(srv, wrapper)
	}
}

type recvWrapper struct {
	grpc.ServerStream
}

func (s *recvWrapper) RecvMsg(m interface{}) error {
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	if err := validate(m); err != nil {
		return err
	}

	return nil
}
```

- edit `demo/demo.go` to use this middleware.

```go
func main () {
    // ...
    s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		demo.RegisterDemoServer(grpcServer, server.NewDemoServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
    s.AddUnaryInterceptors(validate.UnaryServerInterceptor())
    group.Add(s)
    // ...
}
```

4. run it.

```bash
go run demo.go
```

## Start a Kubernetes cluster by Kind

```bash
# CASE A
kind create cluster --name lab --config deployment/kind-cluster.yaml
# CASE B-E
```

## Change HPA sync period

In order to make the HPA controller scale faster, we need change the period of `--horizontal-pod-autoscaler-sync-period` (default is 15s) in `kube-controller-manager`.
We need change this config to 5s.

```bash
# access control-panel
docker exec -it lab-control-plane bash

# edit it
apt update
apt install nano
nano /etc/kubernetes/manifests/kube-controller-manager.yaml
# leave docker contianer
exit

# after edited it, restart control-plane
docker restart lab-control-plane
```

`/etc/kubernetes/manifests/kube-controller-manager.yaml`
```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    component: kube-controller-manager
    tier: control-plane
  name: kube-controller-manager
  namespace: kube-system
spec:
  containers:
  - command:
    - kube-controller-manager
    - --allocate-node-cidrs=true
    - --authentication-kubeconfig=/etc/kubernetes/controller-manager.conf
    - --authorization-kubeconfig=/etc/kubernetes/controller-manager.conf
    - --bind-address=127.0.0.1
    - --client-ca-file=/etc/kubernetes/pki/ca.crt
    - --cluster-cidr=10.244.0.0/16
    - --cluster-name=mac
    - --cluster-signing-cert-file=/etc/kubernetes/pki/ca.crt
    - --cluster-signing-key-file=/etc/kubernetes/pki/ca.key
    - --controllers=*,bootstrapsigner,tokencleaner
    - --enable-hostpath-provisioner=true
    - --kubeconfig=/etc/kubernetes/controller-manager.conf
    - --leader-elect=true
    - --requestheader-client-ca-file=/etc/kubernetes/pki/front-proxy-ca.crt
    - --root-ca-file=/etc/kubernetes/pki/ca.crt
    - --service-account-private-key-file=/etc/kubernetes/pki/sa.key
    - --service-cluster-ip-range=10.96.0.0/16
    - --use-service-account-credentials=true
    - --horizontal-pod-autoscaler-sync-period=5s # add this line
# ...
```
restart it, it will work!

## Install metric-server

In test environment, we need add a args of `--kubelet-insecure-tls`. `deployment/metric-server.yaml` was configured.

```bash
kubectl apply -f deployment/metric-server.yaml
```

Test it.
```bash
kubectl top pod
kubectl top node
```

## Istall Kubernetes Dashboard

```bash
kubectl apply -f https://github.com/kubernetes/dashboard/blob/master/charts/kubernetes-dashboard.yaml

# 建立admin user
kubectl apply -f - <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kube-system-default
  labels:
    k8s-app: kube-system
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: default
    namespace: kube-system

---

apiVersion: v1
kind: Secret
metadata:
  name: default
  namespace: kube-system
  labels:
    k8s-app: kube-system
  annotations:
    kubernetes.io/service-account.name: default
type: kubernetes.io/service-account-token
EOF

TOKEN=$(kubectl -n kube-system describe secret default| awk '$1=="token:"{print $2}')

kubectl config set-credentials docker-desktop --token="${TOKEN}"

echo $TOKEN

kubectl proxy
```

access with this url http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/.

## Set up Case A

1. install microservices and redis.

```bash
kubectl apply -f ./deployment/case-A.yaml
```

2. install nginx ingress.

```bash
# install nginx ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

3. wait it...

```bash
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

4. delete it...

```bash
kubectl delete -f ./deployment/case-A.yaml
kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

## Set up Case B

1. install istio and set the default  configuration profile.

```bash
istioctl install --set profile=default -y
# kubectl patch deployments.apps -n istio-system istio-ingressgateway -p '{"spec":{"template":{"spec":{"containers":[{"name":"istio-proxy","ports":[{"containerPort":8080,"hostPort":80},{"containerPort":8443,"hostPort":443}]}]}}}}'
```

2. add a namespace label to instruct Istio to automatically inject Envoy sidecar proxies.

```bash
kubectl label namespace default istio-injection=enabled
```

3. create a root certificate and private key to sign the certificates for our services.

```bash
mkdir certs
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=local Inc./CN=local.com' -keyout certs/local.com.key -out certs/local.com.crt
```

4. generate a certificate and a private key for lab.local.com.

```bash
openssl req -out certs/lab.local.com.csr -newkey rsa:2048 -nodes -keyout certs/lab.local.com.key -subj "/CN=lab.local.com/O=lab organization"
openssl x509 -req -sha256 -days 365 -CA certs/local.com.crt -CAkey certs/local.com.key -set_serial 0 -in certs/lab.local.com.csr -out certs/lab.local.com.crt
```

5. create a secret for the ingress gateway

```bash
kubectl create -n istio-system secret tls lab-credential \
  --key=certs/lab.local.com.key \
  --cert=certs/lab.local.com.crt
```

6. install microservices and redis.

```bash
kubectl apply -f ./deployment/case-B.yaml
```

7. edit the Service type of istio ingress gateway from `LoadBalancer` to `NodePort` and change to port.

```bash
kubectl patch svc istio-ingressgateway -n istio-system --patch-file ./deployment/gateway-svc-patch.yaml
```

8. delete it...

```bash
kubectl delete -f ./deployment/case-B.yaml
istioctl uninstall --purge
kubectl delete -n istio-system secret tls lab-credential
```