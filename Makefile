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
	curl http://localhost:8282/reconnect
.PHONY: api-recon

api-cb:
	curl http://localhost:8282/callbacks
.PHONY: api-cb

tail-api:
	docker logs api -f
tail-worker1: 
	docker logs worker1 -f 
tail-worker2:
	docker logs worker2 -f
.PHONY: tail-api tail-worker1 tail-worker2


worker-health:
	curl http://localhost:8484/healthz
.PHONY: worker-health
 

worker-health-2:
	curl http://localhost:8686/healthz
.PHONY: worker-health-2

 

stop:
	docker stop api worker1 worker2 nats
.PHONY: stop

start:
	docker start nats worker1 worker2 api
.PHONY: start


rebuild-server:
	docker stop nats
	docker rmi nats -f
	docker-compose up --build -d
.PHONY: rebuild-server

stop-server:
	docker stop nats
.PHONY: stop-server

restart-server:
	docker restart nats
.PHONY: restart-server

restart-non-server:
	docker restart worker1 worker2 api
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
