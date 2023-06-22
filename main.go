package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/gopkg/util/gctuner"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/memnix/memnix-rest/app/http"
	"github.com/memnix/memnix-rest/app/meilisearch"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func main() {
	// Setup the logger
	zapLogger, undo := logger.CreateZapLogger()

	// Setup the environment variables
	setupEnv()

	// Setup the garbage collector
	gcTuning()

	// Setup the infrastructures
	setupInfrastructures()

	if !fiber.IsChild() {
		// Migrate the models
		migrate()

		// Init MeiliSearch
		err := meilisearch.InitMeiliSearch(internal.InitializeMeiliSearch())
		if err != nil {
			zapLogger.Error("error initializing meilisearch", zap.Error(err))
		}
	}

	zapLogger.Info("starting server")

	// Create the app
	app := http.New()

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":1815"); err != nil {
			zapLogger.Panic("error starting server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received

	shutdown(app)

	zapLogger.Info("server stopped")

	if err := zapLogger.Sync(); err != nil {
		return // can't even log, just exit
	}
	undo()
}

func shutdown(app *fiber.App) {
	otelzap.L().Info("🔒 Server shutting down...")
	_ = app.Shutdown()

	otelzap.L().Info("🧹 Running cleanup tasks...")

	err := infrastructures.DisconnectDB()
	if err != nil {
		otelzap.L().Error("❌ Error closing database connection")
	} else {
		otelzap.L().Info("✅ Disconnected from database")
	}

	err = infrastructures.CloseRedis()
	if err != nil {
		otelzap.L().Error("❌ Error closing Redis connection")
	} else {
		otelzap.L().Info("✅ Disconnected from Redis")
	}

	err = infrastructures.ShutdownTracer()
	if err != nil {
		otelzap.L().Error("❌ Error closing Tracer connection")
	} else {
		otelzap.L().Info("✅ Disconnected from Tracer")
	}
}

func migrate() {
	// Models to migrate
	migrates := []domain.Model{
		&domain.User{}, &domain.Card{}, &domain.Deck{}, &domain.Mcq{},
	}

	otelzap.L().Info("⚙️ Starting database migration...")

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		step := i + 1
		err := infrastructures.GetDBConn().AutoMigrate(&migrates[i])
		if err != nil {
			otelzap.L().Error(fmt.Sprintf("❌ Error migrating model %s %d/%d", migrates[i].TableName(), step, len(migrates)))
		} else {
			otelzap.L().Info(fmt.Sprintf("✅ Migration completed for model %s %d/%d", migrates[i].TableName(), step, len(migrates)))
		}
	}

	otelzap.L().Info("✅ Database migration completed!")
}

func setupEnv() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		otelzap.L().Fatal("❌ Error loading .env file")
	}

	// Init oauth
	infrastructures.InitOauth()
}

func setupInfrastructures() {
	err := infrastructures.ConnectDB()
	if err != nil {
		otelzap.L().Fatal("❌ Error connecting to database")
	} else {
		otelzap.L().Info("✅ Connected to database")
	}

	// Redis connection
	err = infrastructures.ConnectRedis()
	if err != nil {
		otelzap.L().Fatal("❌ Error connecting to Redis")
	} else {
		otelzap.L().Info("✅ Connected to Redis")
	}

	// Connect MeiliSearch
	err = infrastructures.ConnectMeiliSearch(config.EnvHelper)
	if err != nil {
		otelzap.L().Fatal("❌ Error connecting to MeiliSearch")
	} else {
		otelzap.L().Info("✅ Connected to MeiliSearch")
	}

	// Connect to the tracer
	err = infrastructures.InitTracer()
	if err != nil {
		otelzap.L().Fatal("❌ Error connecting to the tracer")
	} else {
		otelzap.L().Info("✅ Connected to the tracer")
	}

	if err = infrastructures.CreateRistrettoCache(); err != nil {
		otelzap.L().Fatal("❌ Error creating Ristretto cache")
	} else {
		otelzap.L().Info("✅ Created Ristretto cache")
	}
}

func gcTuning() {
	var limit float64 = 4 * config.GCLimit
	// Set the GC threshold to 70% of the limit
	threshold := uint64(limit * config.GCThresholdPercent)

	gctuner.Tuning(threshold)

	otelzap.L().Info(fmt.Sprintf("🔧 GC Tuning - Limit: %.2f GB, Threshold: %d bytes, GC Percent: %d, Min GC Percent: %d, Max GC Percent: %d",
		limit/(config.GCLimit),
		threshold,
		gctuner.GetGCPercent(),
		gctuner.GetMinGCPercent(),
		gctuner.GetMaxGCPercent()))

	otelzap.L().Info("✅ GC Tuning completed!")
}
