up:
	docker-compose up --build -d
.PHONY: up

api-health:
	curl http://localhost:8282/healthz
.PHONY: api-health

api-task:
	curl http://localhost:8282/createTask
.PHONY: api-task

api-logs:
	docker logs -f nats-docker_api_1
.PHONY: api-logs


worker-health:
	curl http://localhost:8484/healthz
.PHONY: worker-health

worker-logs:
	docker logs -f nats-docker_worker_1
.PHONY: worker-logs

stop:
	docker stop nats-docker_api_1 nats-docker_worker_1 nats-docker_nats_1
.PHONY: stop

start:
	docker start nats-docker_nats_1 nats-docker_worker_1 nats-docker_api_1
.PHONY: start

