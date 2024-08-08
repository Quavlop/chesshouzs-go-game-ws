generate-go-rpc:
	protoc --proto_path=./services/rpc/proto --go_out=./services/rpc/pb --go-grpc_out=./services/rpc/pb ./services/rpc/proto/match_service.proto
	protoc --proto_path=./services/rpc/proto --go_out=./services/rpc/pb --go-grpc_out=./services/rpc/pb ./services/rpc/proto/match_service_model.proto