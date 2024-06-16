// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces/repository_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	redis "github.com/redis/go-redis/v9"
	models "ingenhouzs.com/chesshouzs/go-game/models"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// DeletePlayerFromPool mocks base method.
func (m *MockRepository) DeletePlayerFromPool(params models.PlayerPoolParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlayerFromPool", params)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlayerFromPool indicates an expected call of DeletePlayerFromPool.
func (mr *MockRepositoryMockRecorder) DeletePlayerFromPool(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlayerFromPool", reflect.TypeOf((*MockRepository)(nil).DeletePlayerFromPool), params)
}

// DeletePlayerOnPoolDataToRedis mocks base method.
func (m *MockRepository) DeletePlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlayerOnPoolDataToRedis", params, joinTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlayerOnPoolDataToRedis indicates an expected call of DeletePlayerOnPoolDataToRedis.
func (mr *MockRepositoryMockRecorder) DeletePlayerOnPoolDataToRedis(params, joinTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlayerOnPoolDataToRedis", reflect.TypeOf((*MockRepository)(nil).DeletePlayerOnPoolDataToRedis), params, joinTime)
}

// GetGameTypeVariant mocks base method.
func (m *MockRepository) GetGameTypeVariant(params models.GameTypeVariant) ([]models.GameTypeVariant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGameTypeVariant", params)
	ret0, _ := ret[0].([]models.GameTypeVariant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGameTypeVariant indicates an expected call of GetGameTypeVariant.
func (mr *MockRepositoryMockRecorder) GetGameTypeVariant(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGameTypeVariant", reflect.TypeOf((*MockRepository)(nil).GetGameTypeVariant), params)
}

// GetPlayerPoolData mocks base method.
func (m *MockRepository) GetPlayerPoolData(params models.PlayerPoolParams) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlayerPoolData", params)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlayerPoolData indicates an expected call of GetPlayerPoolData.
func (mr *MockRepositoryMockRecorder) GetPlayerPoolData(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlayerPoolData", reflect.TypeOf((*MockRepository)(nil).GetPlayerPoolData), params)
}

// GetUnderMatchmakingPlayers mocks base method.
func (m *MockRepository) GetUnderMatchmakingPlayers(params models.PoolParams) ([]models.PlayerPool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnderMatchmakingPlayers", params)
	ret0, _ := ret[0].([]models.PlayerPool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnderMatchmakingPlayers indicates an expected call of GetUnderMatchmakingPlayers.
func (mr *MockRepositoryMockRecorder) GetUnderMatchmakingPlayers(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnderMatchmakingPlayers", reflect.TypeOf((*MockRepository)(nil).GetUnderMatchmakingPlayers), params)
}

// GetUserDataByID mocks base method.
func (m *MockRepository) GetUserDataByID(id string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDataByID", id)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDataByID indicates an expected call of GetUserDataByID.
func (mr *MockRepositoryMockRecorder) GetUserDataByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDataByID", reflect.TypeOf((*MockRepository)(nil).GetUserDataByID), id)
}

// InsertGameData mocks base method.
func (m *MockRepository) InsertGameData(params models.InsertGameParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertGameData", params)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertGameData indicates an expected call of InsertGameData.
func (mr *MockRepositoryMockRecorder) InsertGameData(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertGameData", reflect.TypeOf((*MockRepository)(nil).InsertGameData), params)
}

// InsertMoveCacheIdentifier mocks base method.
func (m *MockRepository) InsertMoveCacheIdentifier(params models.MoveCache) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMoveCacheIdentifier", params)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertMoveCacheIdentifier indicates an expected call of InsertMoveCacheIdentifier.
func (mr *MockRepositoryMockRecorder) InsertMoveCacheIdentifier(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMoveCacheIdentifier", reflect.TypeOf((*MockRepository)(nil).InsertMoveCacheIdentifier), params)
}

// InsertPlayerIntoPool mocks base method.
func (m *MockRepository) InsertPlayerIntoPool(params models.PlayerPoolParams, joinTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPlayerIntoPool", params, joinTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertPlayerIntoPool indicates an expected call of InsertPlayerIntoPool.
func (mr *MockRepositoryMockRecorder) InsertPlayerIntoPool(params, joinTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPlayerIntoPool", reflect.TypeOf((*MockRepository)(nil).InsertPlayerIntoPool), params, joinTime)
}

// InsertPlayerOnPoolDataToRedis mocks base method.
func (m *MockRepository) InsertPlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPlayerOnPoolDataToRedis", params, joinTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertPlayerOnPoolDataToRedis indicates an expected call of InsertPlayerOnPoolDataToRedis.
func (mr *MockRepositoryMockRecorder) InsertPlayerOnPoolDataToRedis(params, joinTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPlayerOnPoolDataToRedis", reflect.TypeOf((*MockRepository)(nil).InsertPlayerOnPoolDataToRedis), params, joinTime)
}

// WithRedisTrx mocks base method.
func (m *MockRepository) WithRedisTrx(ctx context.Context, keys []string, fn func(redis.Pipeliner) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithRedisTrx", ctx, keys, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithRedisTrx indicates an expected call of WithRedisTrx.
func (mr *MockRepositoryMockRecorder) WithRedisTrx(ctx, keys, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithRedisTrx", reflect.TypeOf((*MockRepository)(nil).WithRedisTrx), ctx, keys, fn)
}

// MockMatchRepository is a mock of MatchRepository interface.
type MockMatchRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMatchRepositoryMockRecorder
}

// MockMatchRepositoryMockRecorder is the mock recorder for MockMatchRepository.
type MockMatchRepositoryMockRecorder struct {
	mock *MockMatchRepository
}

// NewMockMatchRepository creates a new mock instance.
func NewMockMatchRepository(ctrl *gomock.Controller) *MockMatchRepository {
	mock := &MockMatchRepository{ctrl: ctrl}
	mock.recorder = &MockMatchRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMatchRepository) EXPECT() *MockMatchRepositoryMockRecorder {
	return m.recorder
}

// DeletePlayerFromPool mocks base method.
func (m *MockMatchRepository) DeletePlayerFromPool(params models.PlayerPoolParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlayerFromPool", params)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlayerFromPool indicates an expected call of DeletePlayerFromPool.
func (mr *MockMatchRepositoryMockRecorder) DeletePlayerFromPool(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlayerFromPool", reflect.TypeOf((*MockMatchRepository)(nil).DeletePlayerFromPool), params)
}

// DeletePlayerOnPoolDataToRedis mocks base method.
func (m *MockMatchRepository) DeletePlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlayerOnPoolDataToRedis", params, joinTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlayerOnPoolDataToRedis indicates an expected call of DeletePlayerOnPoolDataToRedis.
func (mr *MockMatchRepositoryMockRecorder) DeletePlayerOnPoolDataToRedis(params, joinTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlayerOnPoolDataToRedis", reflect.TypeOf((*MockMatchRepository)(nil).DeletePlayerOnPoolDataToRedis), params, joinTime)
}

// GetPlayerPoolData mocks base method.
func (m *MockMatchRepository) GetPlayerPoolData(params models.PlayerPoolParams) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlayerPoolData", params)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlayerPoolData indicates an expected call of GetPlayerPoolData.
func (mr *MockMatchRepositoryMockRecorder) GetPlayerPoolData(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlayerPoolData", reflect.TypeOf((*MockMatchRepository)(nil).GetPlayerPoolData), params)
}

// GetUnderMatchmakingPlayers mocks base method.
func (m *MockMatchRepository) GetUnderMatchmakingPlayers(params models.PoolParams) ([]models.PlayerPool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnderMatchmakingPlayers", params)
	ret0, _ := ret[0].([]models.PlayerPool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnderMatchmakingPlayers indicates an expected call of GetUnderMatchmakingPlayers.
func (mr *MockMatchRepositoryMockRecorder) GetUnderMatchmakingPlayers(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnderMatchmakingPlayers", reflect.TypeOf((*MockMatchRepository)(nil).GetUnderMatchmakingPlayers), params)
}

// InsertGameData mocks base method.
func (m *MockMatchRepository) InsertGameData(params models.InsertGameParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertGameData", params)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertGameData indicates an expected call of InsertGameData.
func (mr *MockMatchRepositoryMockRecorder) InsertGameData(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertGameData", reflect.TypeOf((*MockMatchRepository)(nil).InsertGameData), params)
}

// InsertMoveCacheIdentifier mocks base method.
func (m *MockMatchRepository) InsertMoveCacheIdentifier(params models.MoveCache) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMoveCacheIdentifier", params)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertMoveCacheIdentifier indicates an expected call of InsertMoveCacheIdentifier.
func (mr *MockMatchRepositoryMockRecorder) InsertMoveCacheIdentifier(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMoveCacheIdentifier", reflect.TypeOf((*MockMatchRepository)(nil).InsertMoveCacheIdentifier), params)
}

// InsertPlayerIntoPool mocks base method.
func (m *MockMatchRepository) InsertPlayerIntoPool(params models.PlayerPoolParams, joinTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPlayerIntoPool", params, joinTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertPlayerIntoPool indicates an expected call of InsertPlayerIntoPool.
func (mr *MockMatchRepositoryMockRecorder) InsertPlayerIntoPool(params, joinTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPlayerIntoPool", reflect.TypeOf((*MockMatchRepository)(nil).InsertPlayerIntoPool), params, joinTime)
}

// InsertPlayerOnPoolDataToRedis mocks base method.
func (m *MockMatchRepository) InsertPlayerOnPoolDataToRedis(params models.PlayerPoolParams, joinTime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertPlayerOnPoolDataToRedis", params, joinTime)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertPlayerOnPoolDataToRedis indicates an expected call of InsertPlayerOnPoolDataToRedis.
func (mr *MockMatchRepositoryMockRecorder) InsertPlayerOnPoolDataToRedis(params, joinTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertPlayerOnPoolDataToRedis", reflect.TypeOf((*MockMatchRepository)(nil).InsertPlayerOnPoolDataToRedis), params, joinTime)
}

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// GetUserDataByID mocks base method.
func (m *MockUserRepository) GetUserDataByID(id string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserDataByID", id)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserDataByID indicates an expected call of GetUserDataByID.
func (mr *MockUserRepositoryMockRecorder) GetUserDataByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserDataByID", reflect.TypeOf((*MockUserRepository)(nil).GetUserDataByID), id)
}

// MockGameRepository is a mock of GameRepository interface.
type MockGameRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGameRepositoryMockRecorder
}

// MockGameRepositoryMockRecorder is the mock recorder for MockGameRepository.
type MockGameRepositoryMockRecorder struct {
	mock *MockGameRepository
}

// NewMockGameRepository creates a new mock instance.
func NewMockGameRepository(ctrl *gomock.Controller) *MockGameRepository {
	mock := &MockGameRepository{ctrl: ctrl}
	mock.recorder = &MockGameRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGameRepository) EXPECT() *MockGameRepositoryMockRecorder {
	return m.recorder
}

// GetGameTypeVariant mocks base method.
func (m *MockGameRepository) GetGameTypeVariant(params models.GameTypeVariant) ([]models.GameTypeVariant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGameTypeVariant", params)
	ret0, _ := ret[0].([]models.GameTypeVariant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGameTypeVariant indicates an expected call of GetGameTypeVariant.
func (mr *MockGameRepositoryMockRecorder) GetGameTypeVariant(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGameTypeVariant", reflect.TypeOf((*MockGameRepository)(nil).GetGameTypeVariant), params)
}

// MockTransaction is a mock of Transaction interface.
type MockTransaction struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionMockRecorder
}

// MockTransactionMockRecorder is the mock recorder for MockTransaction.
type MockTransactionMockRecorder struct {
	mock *MockTransaction
}

// NewMockTransaction creates a new mock instance.
func NewMockTransaction(ctrl *gomock.Controller) *MockTransaction {
	mock := &MockTransaction{ctrl: ctrl}
	mock.recorder = &MockTransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransaction) EXPECT() *MockTransactionMockRecorder {
	return m.recorder
}

// WithRedisTrx mocks base method.
func (m *MockTransaction) WithRedisTrx(ctx context.Context, keys []string, fn func(redis.Pipeliner) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithRedisTrx", ctx, keys, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithRedisTrx indicates an expected call of WithRedisTrx.
func (mr *MockTransactionMockRecorder) WithRedisTrx(ctx, keys, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithRedisTrx", reflect.TypeOf((*MockTransaction)(nil).WithRedisTrx), ctx, keys, fn)
}
