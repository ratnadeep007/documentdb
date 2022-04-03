SERVER := "server/"
CLIENT := "client/"
DATABASE_FOLDER := "data.doc"

build-server:
	cd $(SERVER) && go build

start-server: build-server
	cd $(SERVER) && ./documentdb

js-dev-server:
	cd $(CLIENT) && npm start

clean:
	cd $(SERVER) && rm documentdb && rm -rf $(DATABASE_FOLDER)