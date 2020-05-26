up:
	docker-compose up --build -d
.PHONY: up

api-health:
	curl http://localhost:8282/healthz
.PHONY: api-health

api-task:
	curl http://localhost:8282/createTask
.PHONY: api-task

api-recon:
	curl http://localhost:8282/reconnectTask
.PHONY: api-recon


api-logs:
	docker logs -f nats-docker_api_1
.PHONY: api-logs


worker-health:
	curl http://localhost:8484/healthz
.PHONY: worker-health

worker-logs:
	docker logs -f nats-docker_worker_1
.PHONY: worker-logs


worker-health-2:
	curl http://localhost:8686/healthz
.PHONY: worker-health-2

worker-logs-2:
	docker logs -f nats-docker_worker2_1
.PHONY: worker-logs-2


stop:
	docker stop nats-docker_api_1 nats-docker_worker_1 nats-docker_worker2_1  nats-docker_natsx_1
.PHONY: stop

start:
	docker start nats-docker_natsx_1 nats-docker_worker_1 nats-docker_worker2_1  nats-docker_api_1
.PHONY: start



restart-server:
	docker restart nats-docker_natsx_1
.PHONY: restart-server

restart-non-server:
	docker restart nats-docker_worker_1 nats-docker_worker2_1  nats-docker_api_1
.PHONY: restart-non-server


# informmation abou the server
info-server:
	watch -n 1 curl http://localhost:8222/varz --silent

	# message throutput
	# watch -n 1 'curl http://localhost:8222/varz --silent | grep "\(msgs\|byte\)"'
.PHONY: info-server

# statistics and metadata about the clients currently connected to the server.
info-conn:
	watch -n 1 curl http://localhost:8222/connz --silent
.PHONY: info-conn

# cumulative stats about internal state of the sublist data structur that server maintains.
info-sublist:
	watch -n 1 curl http://localhost:8222/subsz --silent
.PHONY: info-sublist