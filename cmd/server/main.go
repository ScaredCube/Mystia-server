package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"mystia-voice-backend/internal/db"
	"mystia-voice-backend/internal/livekit"
	"mystia-voice-backend/internal/service"
	"mystia-voice-backend/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// CLI Flags
	flagPort := flag.String("port", "", "Mystia gRPC server port")
	flagLkPort := flag.String("lk-port", "", "LiveKit TCP/Signaling port")
	flag.Parse()

	// Configuration
	dbPath := getEnv("DB_PATH", "voice.db")
	jwtSecret := []byte(getEnv("JWT_SECRET", "super-secret-key"))
	lkApiKey := getEnv("LIVEKIT_API_KEY", "devkey")
	lkApiSecret := getEnv("LIVEKIT_API_SECRET", "secret")
	lkHost := getEnv("LIVEKIT_HOST", "http://127.0.0.1:7880")
	port := getEnv("PORT", "50051")

	if *flagPort != "" {
		port = *flagPort
	}
	if *flagLkPort != "" {
		if strings.Contains(lkHost, ":") {
			lastColon := strings.LastIndex(lkHost, ":")
			lkHost = lkHost[:lastColon+1] + *flagLkPort
		} else {
			lkHost = lkHost + ":" + *flagLkPort
		}
	}

	// Initialize DB
	database, err := db.NewDB(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	// Initialize LiveKit
	lkProvider := livekit.NewLiveKitProvider(lkApiKey, lkApiSecret, lkHost)

	// Initialize gRPC Server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(service.AuthInterceptor(jwtSecret)),
	)

	// Register Services
	proto.RegisterAuthServiceServer(s, &service.AuthService{
		DB:        database,
		JWTSecret: jwtSecret,
	})
	proto.RegisterChannelServiceServer(s, &service.ChannelService{
		DB:      database,
		LiveKit: lkProvider,
	})
	proto.RegisterAdminServiceServer(s, &service.AdminService{
		DB: database,
	})

	reflection.Register(s)

	lkManager := livekit.NewServerManager()
	lkPort := "7880"
	if strings.Contains(lkHost, ":") {
		parts := strings.Split(lkHost, ":")
		lkPort = parts[len(parts)-1]
	}

	if err := lkManager.Start(lkApiKey, lkApiSecret, lkPort); err != nil {
		log.Printf("Warning: failed to start LiveKit server: %v (Is livekit-server.exe in the current directory?)", err)
	} else {
		defer lkManager.Stop()
	}

	log.Printf("Starting voice backend on port %s...", port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down...")
	s.GracefulStop()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
