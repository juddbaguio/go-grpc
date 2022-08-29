generate-pb-hello:
	rm -rf ./grpc/hello
	mkdir ./grpc/hello
	protoc --proto_path=proto --go_out=grpc/hello --go-grpc_out=grpc/hello \
	--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
	hello.proto

generate-pb-auth:
	rm -rf ./grpc/auth
	mkdir ./grpc/auth
	protoc --proto_path=proto --go_out=grpc/auth --go-grpc_out=grpc/auth \
	--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
	auth.proto

gen-pb: generate-pb-hello generate-pb-auth

build-auth:
	docker build -t juddbaguio/auth:latest -f ./infra/docker/auth.Dockerfile .

build-hello:
	docker build -t juddbaguio/hello:latest -f ./infra/docker/hello.Dockerfile .

build-docker: build-auth build-hello

internal-nginx:
	helm install grpc-nginx ingress-nginx/ingress-nginx  \
	--namespace default \
	--set controller.electionID=grpc-internal-nginx \
	--set controller.ingressClassResource.name=nginx-internal \
	--set controller.ingressClassResource.controllerValue="juddbaguio/grpc-internal-nginx" \
	--set controller.ingressClassResource.enabled=true \
	--set controller.ingressClassByName=true \
	--set controller.extraArgs.default-ssl-certificate="default/grpc-tls"

external-nginx:
	helm install http-nginx ingress-nginx/ingress-nginx  \
	--namespace default \
	--set controller.electionID=http-external-nginx \
	--set controller.ingressClassResource.name=nginx-external \
	--set controller.ingressClassResource.controllerValue="juddbaguio/http-external-nginx" \
	--set controller.ingressClassResource.enabled=true \
	--set controller.ingressClassByName=true

nginx: internal-nginx

k8s-deployment:
	kubectl apply -f ./infra/k8s/auth.deployment.yaml -f ./infra/k8s/hello.deployment.yaml

k8s-service:
	kubectl apply -f ./infra/k8s/service.yaml

k8s-ingress:
	kubectl apply -f ./infra/k8s/ingress.yaml

k8s: k8s-deployment k8s-service

k8s-tls:
	kubectl create secret tls grpc-tls --key ./infra/k8s/tls/tls.key --cert ./infra/k8s/tls/tls.crt

delete-istio-tls:
	kubectl -n istio-system delete secret grpc-tls

istio-tls: delete-istio-tls
	kubectl -n istio-system create secret tls grpc-tls --save-config --key ./infra/k8s/tls/tls.key --cert ./infra/k8s/tls/tls.crt

update-k8s-tls:
	kubectl update secret tls grpc-tls --key ./infra/k8s/tls/tls.key --cert ./infra/k8s/tls/tls.crt

gen-tls:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ./infra/k8s/tls/tls.key \
	-out ./infra/k8s/tls/tls.crt -subj "/CN=*.juddbaguio.dev" \
	-addext "subjectAltName = DNS:*.juddbaguio.dev"

run-client:
	go run ./client

istio:
	kubectl apply -f ./infra/k8s/istio