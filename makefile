all:
	cd client && go build;
	cd server && go build;

delete:
	rm client/client;
	rm server/server;
	rm -f server/tls_key.log;