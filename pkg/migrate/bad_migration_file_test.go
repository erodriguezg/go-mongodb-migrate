package migrate

import (
	"testing"

	_ "embed"

	"go.mongodb.org/mongo-driver/mongo"
)

//go:embed bad_migration_file_test.go
var badMigrationFileTest string

func TestBadMigrationFile(t *testing.T) {
	oldMigrate := globalMigrate
	defer func() {
		globalMigrate = oldMigrate
	}()
	globalMigrate = NewMigrate(nil)

	err := Register(&badMigrationFileTest, func(db *mongo.Database) error {
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
	if err == nil {
		t.Errorf("Unexpected nil error")
	}
}

func TestBadMigrationFilePanic(t *testing.T) {
	oldMigrate := globalMigrate
	defer func() {
		globalMigrate = oldMigrate
		if r := recover(); r == nil {
			t.Errorf("Unexpectedly no panic recovered")
		}
	}()
	globalMigrate = NewMigrate(nil)
	MustRegister(&badMigrationFileTest, func(db *mongo.Database) error {
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
}
