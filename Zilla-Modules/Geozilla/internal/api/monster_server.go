// internal/api/monster_server.go
type MonsterServer struct {
    CognitiveService *cognitive_processor.NeuralAnalyzer
    SecurityService  *threat_detection.AnomalyDetector
    ExecutionService *autonomous_actors.JackpotExecutor
}

func (ms *MonsterServer) StartGRPCServer() {
    // High-performance gRPC for AI communication
    pb.RegisterMonsterServiceServer(grpcServer, ms)
}
