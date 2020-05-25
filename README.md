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

5. Send api request (new terminal)
	> make api-task


