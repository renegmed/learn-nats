## Lesson learned with NATS using Golang ##


To install and run the application:

1. Start the NATS service, api and 
	$ make up

2. Set the worker log monitor (new terminal)
	$ make tail-worker1

3. Set the worker log monitor 2 (new terminal)
	$ make tail-worker2

4. Set the api log monitor (new terminal)
	$> make tail-api

5. Send api reconnect task (new terminal)
	$ make api-recon

6. Stop the server
	$ make stop-server

7. Restart the server
	$ make restart-server

