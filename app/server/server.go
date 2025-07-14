package server

import (
	"context"
	"medods-auth/persistance/postgres"
	"medods-auth/service/auth"
	"medods-auth/token"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func Start() error {
	server := setupServer()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-osSignal

	return shutdownServer(server)
}

func setupServer() *http.Server {
	db, err := postgres.InitDatabase(
		Config().Postgres,
	)
	if err != nil {
		panic(err)
	}
	hashRepo := postgres.NewHashRepository(db)
	blacklistRepo := postgres.NewBlackListRepository(db)

	accessTTL := time.Minute * 5
	refreshTTL := time.Hour * 48

	authService, err := auth.NewAuthService(auth.AuthServiceOptions{
		RefreshTokenRepo: hashRepo,
		Blacklist:        blacklistRepo,
		Generator:        &token.SHA512Generator{},
		Hasher:           token.BcryptHasher{},

		Secret: Config().HashSecret,

		// TODO: access TTL and refresh TTL from config
		AccessTTL:  &accessTTL,
		RefreshTTL: &refreshTTL,
	})
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/generate", newGenerateHandler(authService))
	router.POST("/refresh", newRefreshHandler(authService))
	router.POST("/me", newMeHandler(authService))
	router.POST("/logout", newLogoutHandler(authService))

	server := http.Server{
		Addr:    ":" + Config().Port,
		Handler: router,
	}
	return &server
}

func shutdownServer(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
