package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/unitechio/einfra-be/internal/config"
	"github.com/unitechio/einfra-be/internal/domain"
)

func NewPostgresConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.SSLMode,
	)

	var gormLogger logger.Interface
	if cfg.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("✅ Connected to PostgreSQL database: %s", cfg.User)
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	err := db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.RefreshToken{},
		&domain.Session{},
		&domain.LoginAttempt{},
		&domain.AuditLog{},
		&domain.Notification{},
		&domain.NotificationTemplate{},
		&domain.NotificationPreference{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

func SeedDefaultData(db *gorm.DB) error {
	log.Println("🌱 Seeding default data...")

	adminPermissions := []domain.Permission{
		{Name: "user.create", Description: "Create users", Resource: "user", Action: "create"},
		{Name: "user.read", Description: "Read users", Resource: "user", Action: "read"},
		{Name: "user.update", Description: "Update users", Resource: "user", Action: "update"},
		{Name: "user.delete", Description: "Delete users", Resource: "user", Action: "delete"},
		{Name: "role.create", Description: "Create roles", Resource: "role", Action: "create"},
		{Name: "role.read", Description: "Read roles", Resource: "role", Action: "read"},
		{Name: "role.update", Description: "Update roles", Resource: "role", Action: "update"},
		{Name: "role.delete", Description: "Delete roles", Resource: "role", Action: "delete"},
	}

	for _, perm := range adminPermissions {
		var existing domain.Permission
		if err := db.Where("name = ?", perm.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&perm).Error; err != nil {
				return fmt.Errorf("failed to create permission %s: %w", perm.Name, err)
			}
		}
	}

	var adminRole domain.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err == gorm.ErrRecordNotFound {
		adminRole = domain.Role{
			Name:        "admin",
			Description: "Administrator role with full access",
			IsActive:    true,
		}
		if err := db.Create(&adminRole).Error; err != nil {
			return fmt.Errorf("failed to create admin role: %w", err)
		}

		var permissions []domain.Permission
		db.Find(&permissions)
		if err := db.Model(&adminRole).Association("Permissions").Append(permissions); err != nil {
			return fmt.Errorf("failed to assign permissions to admin role: %w", err)
		}
	}

	var userRole domain.Role
	if err := db.Where("name = ?", "user").First(&userRole).Error; err == gorm.ErrRecordNotFound {
		userRole = domain.Role{
			Name:        "user",
			Description: "Standard user role",
			IsActive:    true,
		}
		if err := db.Create(&userRole).Error; err != nil {
			return fmt.Errorf("failed to create user role: %w", err)
		}
	}

	log.Println("✅ Default data seeded successfully")
	return nil
}

func InitDatabases(cfgs map[string]config.DatabaseConfig) (map[string]*gorm.DB, error) {
	connections := make(map[string]*gorm.DB)

	for name, dbCfg := range cfgs {
		conn, err := NewPostgresConnection(dbCfg)
		if err != nil {
			// Close already opened connections
			for _, c := range connections {
				sqlDB, _ := c.DB()
				if sqlDB != nil {
					_ = sqlDB.Close()
				}
			}
			return nil, fmt.Errorf("error connecting to %s: %w", name, err)
		}
		connections[name] = conn
		log.Printf("📊 Database '%s' registered successfully", name)
	}

	return connections, nil
}

func CloseAll(connections map[string]*gorm.DB) {
	for name, conn := range connections {
		sqlDB, err := conn.DB()
		if err != nil {
			log.Printf("⚠️  Error getting sqlDB for %s: %v", name, err)
			continue
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("⚠️  Error closing database %s: %v", name, err)
		} else {
			log.Printf("✅ Closed database connection: %s", name)
		}
	}
}
