CURRENT_DIR=$(shell pwd)

DBURL='postgres://postgres:feruza1727@localhost:5432/blacklist?sslmode=disable'

proto-gen:
	./internal/scripts/gen-proto.sh ${CURRENT_DIR}
mig-up:
	migrate -path migrations -database $(DBURL) -verbose up

mig-down:
	migrate -path migrations -database $(DBURL) -verbose down

mig-create:
	migrate create -ext sql -dir migrations -seq create_table

mig-insert:
	migrate create -ext sql -dir migrations -seq insert_data
