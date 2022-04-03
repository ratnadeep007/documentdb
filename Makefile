SERVER := "server/"
CLIENT := "client/"
DATABASE_FOLDER := "data.doc"

build-server:
	cd $(SERVER) && go build

start-server: build-server
	cd $(SERVER) && ./documentdb

clean:
	cd $(SERVER) && rm documentdb && rm -rf $(DATABASE_FOLDER)