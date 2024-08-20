generate-go-rpc:
	protoc --proto_path=./services/rpc/proto --go_out=./services/rpc/pb --go-grpc_out=./services/rpc/pb ./services/rpc/proto/match_service.proto
	protoc --proto_path=./services/rpc/proto --go_out=./services/rpc/pb --go-grpc_out=./services/rpc/pb ./services/rpc/proto/match_service_model.proto
mount-cassandra: 
	docker cp ./archives/keyspace.cql cassandra:/mnt/keyspace.cql
	docker cp ./archives/player_game_states.cql cassandra:/mnt/player_game_states.cql
