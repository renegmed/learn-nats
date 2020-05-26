# Lesson learned with NATS using Golang #

## Chapter 3 The NATS Client ##

Using Subscribe 

Note about the behaior from a subscription is that,
for a single subscription, only a single message will be handled
as a time sequentially, not in parallel.

If we have multiple subscriptions and one of them is processing 
messages slower than the rest, this will not affect other
subscriptions. 

In this example where there are a couple of subscriptions, one on 
a bare subject and another one on a wildcard.


Flush()

Internally, what is does is send everything that has accumulated in 
pending buffer from client. Then, it sends a PING to the server, and 
then waits for the PONG.

As soon as the client receives the PONG reply, the Flush call will 
unblock and let the client assume that the messages that were fired
have been processed by the server.

Using Request

Request API enables the client to publish a message and then wait for 
someone to reply.

Request/Response 

Client would receive only single response among the workers. Request is 
blocked until a response is received or timeout occurs.

States of NATS Connection

Skipping sending message during reconnecting, do the following steps below:

This is for NATS version 1.1
Steps:

1. Start the NATS service, api and 
	> make up

2. Set the worker log monitor (new terminal)
	> make worker-logs

3. Set the worker log monitor 2 (new terminal)
	> make worker-logs2

4. Set the api log monitor (new terminal)
	> make api-logs

5. Send api reconnect task (new terminal)
	> make api-recon

6. Stop the server
	> make stop-server

7. Restart the server
	> make restart-server

