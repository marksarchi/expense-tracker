build: 
	docker build -f zarf/docker/Dockerfile -t expensetracker-1.0 .
run:
	docker run --name expense-tracker -p 8000:8000 expensetracker
compose-up:
	docker-compose -p expense-tracker -f zarf/compose/docker-compose.yaml -f zarf/compose/compose-config.yaml up --detach --remove-orphans