build: 
	docker build -f zarf/docker/Dockerfile -t expensetracker:v1.0.0 .
#run:
#	docker run --name expense-tracker -p 8000:8000 expensetracker-1.1
run: up seed
up:
	docker-compose -p expensetracker-api -f zarf/compose/docker-compose.yaml -f zarf/compose/compose-config.yaml up --detach --remove-orphans
migrate:
		go run app-admin/main.go migrate	
seed: migrate
		go run app-admin/main.go seed

swagger:
	GO111MODULE=off swagger generate spec -o ./app/swagger.yaml --scan-models
go-build: 
	go build .
up-trackerdb:
	docker run --name trackerdb -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres  postgres
test:
		go test ./... -count=1
		staticcheck -checks=all ./...
KIND_CLUSTER := expensetracker-cluster
VERSION := v1.1.0

kind-up:
		kind create cluster \
			--name $(KIND_CLUSTER) \
			--image kindest/node:v1.21.2 \
			--config zarf/k8s/kind/kind-config.yaml
		kubectl config set-context --current --namespace=expensetracker-system

kind-down:
		kind delete cluster --name $(KIND_CLUSTER)
kind-load:
		cd zarf/k8s/kind/expensetracker-pod; kustomize edit set image expensetracker-api-image=36044735/expensetracker:v1.1.0
		kind load docker-image 36044735/expensetracker:$(VERSION) --name $(KIND_CLUSTER)
kind-services: kind-apply
kind-tracker:
		kustomize build zarf/k8s/kind/expensetracker-pod | kubectl apply -f -
kind-apply:
		kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
		kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
		kustomize build zarf/k8s/kind/expensetracker-pod | kubectl apply -f -
kind-update:  kind-load
	kubectl rollout restart deployment e-tracker-pod
kind-shell:
	kubectl exec -it $(shell kubectl get pods | grep e-tracker | cut -c1-30) --container app -- /bin/sh

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=expensetracker				




