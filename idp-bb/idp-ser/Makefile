all:
	protoc \
		--go_out=protos/ --go-grpc_out=protos/ \
		protos/*.proto
		
	protoc \
		--grpc-gateway_out=protos/ \
		protos/login.proto

	protoc \
		--grpc-gateway_out=protos/ \
		protos/user.proto

	protoc \
		--grpc-gateway_out=protos/ \
		protos/token.proto

	protoc \
	--openapiv2_out=gateway/third_party/OpenAPI/ \
    protos/*.proto

clean:
	rm *.pb.go *.pb.gw.go *.swagger.json