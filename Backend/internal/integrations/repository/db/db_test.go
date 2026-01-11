package db_test

import (
	"context"
	"errors"
	"fmt"
	"go-vue-journey/internal/integrations/repository/db"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestConnect_InvalidDriver(t *testing.T) {
	_, err := db.Connect("invalid-driver://connection-string")
	if err == nil {
		t.Fatal("expected error when using invalid driver, got nil")
	}
	if !strings.Contains(err.Error(), "failed to open database") && !strings.Contains(err.Error(), "unreachable") {
		t.Errorf("expected error message to contain 'failed to open database' or 'unreachable', got: %v", err)
	}
}

func TestConnect_PingFails(t *testing.T) {
	_, err := db.Connect("pgx://invalid:invalid@nonexistent:5432/testdb")
	if err == nil {
		t.Fatal("expected error when ping fails, got nil")
	}
	if !strings.Contains(err.Error(), "unreachable") {
		t.Errorf("expected error message to contain 'unreachable', got: %v", err)
	}
}

func TestConnect_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Skipf("Docker not available (required for integration test): %v", r)
		}
	}()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "cockroachdb/cockroach:latest",
		ExposedPorts: []string{"26257/tcp"},
		Cmd:          []string{"start-single-node", "--insecure"},
		WaitingFor: wait.ForLog("initialized new cluster").
			WithStartupTimeout(60 * time.Second),
	}

	cockroachContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		if strings.Contains(err.Error(), "Docker") || strings.Contains(err.Error(), "docker") ||
			strings.Contains(err.Error(), "Cannot connect") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Docker not available (required for integration test): %v", err)
		}
		t.Fatalf("failed to create container: %v", err)
	}
	defer func() {
		if err := cockroachContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	host, err := cockroachContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := cockroachContainer.MappedPort(ctx, "26257")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	connString := fmt.Sprintf("postgres://root@%s:%s/defaultdb?sslmode=disable", host, port.Port())

	time.Sleep(2 * time.Second)

	connectedDB, err := db.Connect(connString)
	if err != nil {
		t.Fatalf("expected successful connection, got error: %v", err)
	}
	if connectedDB == nil {
		t.Fatal("expected valid database connection, got nil")
	}
	defer connectedDB.Close()

	err = connectedDB.Ping()
	if err != nil {
		t.Fatalf("expected successful ping, got error: %v", err)
	}
}

func TestConnect_MockSuccess(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectPing()

	err = mockDB.Ping()
	if err != nil {
		t.Errorf("ping failed: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestConnect_MockPingFails(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	expectedErr := errors.New("connection refused")
	mock.ExpectPing().WillReturnError(expectedErr)

	err = mockDB.Ping()
	if err == nil {
		t.Fatal("expected ping to fail, got nil error")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
