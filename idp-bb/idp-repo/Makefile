all:
	protoc \
		--go_out=protos/ --go-grpc_out=protos/ \
		protos/*.proto

clean:
	rm *.pb.go *.pb.gw.go *.swagger.json