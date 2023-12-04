package tester

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/fbufler/database-monitor/pkg/database"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestSetup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabase(ctrl)
	ctx := context.Background()
	postgresTester := New(Config{
		Databases: []database.Database{
			mockDatabase,
		},
	})
	mockDatabase.EXPECT().SetupTestTable(ctx).Return(nil)
	err := postgresTester.Setup(ctx)
	if err != nil {
		t.Errorf("Setup() error = %v", err)
	}
}

func TestSetupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabase(ctrl)
	ctx := context.Background()
	postgresTester := New(Config{
		Databases: []database.Database{
			mockDatabase,
		},
	})
	setupError := errors.New("SetupTestTable error")
	mockDatabase.EXPECT().SetupTestTable(ctx).Return(setupError)
	err := postgresTester.Setup(ctx)
	assert.EqualError(t, err, fmt.Errorf("setting up databases: %v", []error{setupError}).Error())
}

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabase(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	postgresTester := New(Config{
		TestTimeout:  1,
		TestInterval: 1,
		Databases: []database.Database{
			mockDatabase,
		},
	})
	mockDatabase.EXPECT().Identifier().Return("test")
	mockDatabase.EXPECT().Connect().Return(nil)
	mockDatabase.EXPECT().TestRead(ctx).Return(nil)
	mockDatabase.EXPECT().TestWrite(ctx).Return(nil)
	mockDatabase.EXPECT().Close().Return(nil)
	go postgresTester.Run(ctx)
	result := <-postgresTester.(*TesterImpl).results
	cancel()
	assert.Equal(t, "test", result.Database)
	assert.Equal(t, true, result.Connectable)
	assert.Equal(t, true, result.Readable)
	assert.Equal(t, true, result.Writable)
}

func TestRunDatabaseTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabase(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	postgresTester := New(Config{
		TestTimeout:  1,
		TestInterval: 1,
		Databases: []database.Database{
			mockDatabase,
		},
	})
	mockDatabase.EXPECT().Identifier().Return("test")
	mockDatabase.EXPECT().Connect().Return(nil)
	mockDatabase.EXPECT().TestRead(ctx).Return(nil)
	mockDatabase.EXPECT().TestWrite(ctx).Return(nil)
	mockDatabase.EXPECT().Close().Return(nil)
	go postgresTester.(*TesterImpl).runDatabaseTest(mockDatabase, ctx)
	result := <-postgresTester.(*TesterImpl).results
	cancel()
	assert.Equal(t, "test", result.Database)
	assert.Equal(t, true, result.Connectable)
	assert.Equal(t, true, result.Readable)
	assert.Equal(t, true, result.Writable)
}

func TestRunDatabaseTestConnectError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabase(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	postgresTester := New(Config{
		TestTimeout:  1,
		TestInterval: 1,
		Databases: []database.Database{
			mockDatabase,
		},
	})
	mockDatabase.EXPECT().Identifier().Return("test")
	mockDatabase.EXPECT().Connect().Return(errors.New("Connect error"))
	go postgresTester.(*TesterImpl).runDatabaseTest(mockDatabase, ctx)
	result := <-postgresTester.(*TesterImpl).results
	cancel()
	assert.Equal(t, "test", result.Database)
	assert.Equal(t, false, result.Connectable)
	assert.Equal(t, false, result.Readable)
	assert.Equal(t, false, result.Writable)
}

func TestRunDatabaseTestReadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabase(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	postgresTester := New(Config{
		TestTimeout:  1,
		TestInterval: 1,
		Databases: []database.Database{
			mockDatabase,
		},
	})
	mockDatabase.EXPECT().Identifier().Return("test")
	mockDatabase.EXPECT().Connect().Return(nil)
	mockDatabase.EXPECT().TestRead(ctx).Return(errors.New("Read error"))
	mockDatabase.EXPECT().Close().Return(nil)
	go postgresTester.(*TesterImpl).runDatabaseTest(mockDatabase, ctx)
	result := <-postgresTester.(*TesterImpl).results
	cancel()
	assert.Equal(t, "test", result.Database)
	assert.Equal(t, true, result.Connectable)
	assert.Equal(t, false, result.Readable)
	assert.Equal(t, false, result.Writable)
}

func TestRunDatabaseTestWriteError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatabase := database.NewMockDatabase(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	postgresTester := New(Config{
		TestTimeout:  1,
		TestInterval: 1,
		Databases: []database.Database{
			mockDatabase,
		},
	})
	mockDatabase.EXPECT().Identifier().Return("test")
	mockDatabase.EXPECT().Connect().Return(nil)
	mockDatabase.EXPECT().TestRead(ctx).Return(nil)
	mockDatabase.EXPECT().TestWrite(ctx).Return(errors.New("Write error"))
	mockDatabase.EXPECT().Close().Return(nil)
	go postgresTester.(*TesterImpl).runDatabaseTest(mockDatabase, ctx)
	result := <-postgresTester.(*TesterImpl).results
	cancel()
	assert.Equal(t, "test", result.Database)
	assert.Equal(t, true, result.Connectable)
	assert.Equal(t, true, result.Readable)
	assert.Equal(t, false, result.Writable)
}
