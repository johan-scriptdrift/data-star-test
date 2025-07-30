package sql

import (
	"context"
	"embed"
	"fmt"
	"github.com/delaneyj/toolbelt"
	"github.com/jaswdr/faker/v2"
	"github.com/johan-scriptdrift/data-star-test/sql/zz"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"zombiezen.com/go/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func SetupDB(ctx context.Context, dataFolder string, shouldClear bool) (*toolbelt.Database, error) {
	migrationsDir := "migrations"
	migrationsFiles, err := migrationsFS.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	slices.SortFunc(migrationsFiles, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	migrations := make([]string, len(migrationsFiles))
	for i, file := range migrationsFiles {
		fn := filepath.Join(migrationsDir, file.Name())
		f, err := migrationsFS.Open(fn)
		if err != nil {
			return nil, fmt.Errorf("failed to open migration file %s: %w", fn, err)
		}
		defer f.Close()

		content, err := io.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", fn, err)
		}

		migrations[i] = string(content)
	}

	dbFolder := filepath.Join(dataFolder, "database")
	if shouldClear {
		log.Printf("Clearing database folder: %s", dbFolder)
		if err := os.RemoveAll(dbFolder); err != nil {
			return nil, fmt.Errorf("failed to remove database folder: %w", err)
		}
	}
	dbFilename := filepath.Join(dbFolder, "conduit.sqlite")
	db, err := toolbelt.NewDatabase(ctx, dbFilename, migrations)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	if err := SeedDBIfEmpty(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	return db, nil
}

func SeedDBIfEmpty(ctx context.Context, db *toolbelt.Database) error {
	isEmpty := true
	if err := db.ReadTX(ctx, func(tx *sqlite.Conn) error {
		count, err := zz.OnceCountUsers(tx)
		if err != nil {
			return fmt.Errorf("failed to count users: %w", err)
		}
		isEmpty = count == 0
		return nil
	}); err != nil {
		return err
	}

	if !isEmpty {
		return nil
	}

	now := time.Now().UTC()
	fake := faker.NewWithSeedInt64(0)

	return db.WriteTX(ctx, func(tx *sqlite.Conn) error {
		userIds := make([]int64, 64)
		createUserStmt := zz.CreateUser(tx)
		userIds[0] = 1

		passwordHash, err := bcrypt.GenerateFromPassword([]byte("correctHorseBatteryStapler"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		if err := createUserStmt.Run(&zz.UserModel{
			Id:           1,
			FirstName:    "John",
			LastName:     "Doe",
			Email:        "johndoe@example.com",
			PasswordHash: passwordHash,
			CreatedAt:    now,
		}); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		createLocationStmt := zz.CreateLocation(tx)

		if err = createLocationStmt.Run(&zz.LocationModel{
			Id:        1,
			Lat:       37.7749,
			Long:      -122.4194,
			CreatedAt: now,
		}); err != nil {
			return fmt.Errorf("failed to create location: %w", err)
		}

		for i := 1; i < 64; i++ {
			if err := createLocationStmt.Run(&zz.LocationModel{
				Id:        toolbelt.NextID(),
				Lat:       fake.Float64(4, -90, 90),
				Long:      fake.Float64(4, -180, 180),
				CreatedAt: now,
			}); err != nil {
				return fmt.Errorf("failed to create location: %w", err)
			}
		}

		return nil
	})
}
