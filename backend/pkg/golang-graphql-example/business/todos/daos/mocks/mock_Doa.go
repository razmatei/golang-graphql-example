// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/business/todos/daos (interfaces: Dao)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/mock_Doa.go -package=mocks github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/business/todos/daos Dao
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/business/todos/models"
	pagination "github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/database/pagination"
	gomock "go.uber.org/mock/gomock"
)

// MockDao is a mock of Dao interface.
type MockDao struct {
	ctrl     *gomock.Controller
	recorder *MockDaoMockRecorder
	isgomock struct{}
}

// MockDaoMockRecorder is the mock recorder for MockDao.
type MockDaoMockRecorder struct {
	mock *MockDao
}

// NewMockDao creates a new mock instance.
func NewMockDao(ctrl *gomock.Controller) *MockDao {
	mock := &MockDao{ctrl: ctrl}
	mock.recorder = &MockDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDao) EXPECT() *MockDaoMockRecorder {
	return m.recorder
}

// CreateOrUpdate mocks base method.
func (m *MockDao) CreateOrUpdate(ctx context.Context, tt *models.Todo) (*models.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdate", ctx, tt)
	ret0, _ := ret[0].(*models.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOrUpdate indicates an expected call of CreateOrUpdate.
func (mr *MockDaoMockRecorder) CreateOrUpdate(ctx, tt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdate", reflect.TypeOf((*MockDao)(nil).CreateOrUpdate), ctx, tt)
}

// Find mocks base method.
func (m *MockDao) Find(ctx context.Context, sort []*models.SortOrder, filter *models.Filter, projection *models.Projection) ([]*models.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", ctx, sort, filter, projection)
	ret0, _ := ret[0].([]*models.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockDaoMockRecorder) Find(ctx, sort, filter, projection any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockDao)(nil).Find), ctx, sort, filter, projection)
}

// FindByID mocks base method.
func (m *MockDao) FindByID(ctx context.Context, id string, projection *models.Projection) (*models.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id, projection)
	ret0, _ := ret[0].(*models.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockDaoMockRecorder) FindByID(ctx, id, projection any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockDao)(nil).FindByID), ctx, id, projection)
}

// GetAllPaginated mocks base method.
func (m *MockDao) GetAllPaginated(ctx context.Context, page *pagination.PageInput, sort []*models.SortOrder, filter *models.Filter, projection *models.Projection) ([]*models.Todo, *pagination.PageOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllPaginated", ctx, page, sort, filter, projection)
	ret0, _ := ret[0].([]*models.Todo)
	ret1, _ := ret[1].(*pagination.PageOutput)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllPaginated indicates an expected call of GetAllPaginated.
func (mr *MockDaoMockRecorder) GetAllPaginated(ctx, page, sort, filter, projection any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllPaginated", reflect.TypeOf((*MockDao)(nil).GetAllPaginated), ctx, page, sort, filter, projection)
}

// PatchUpdate mocks base method.
func (m *MockDao) PatchUpdate(ctx context.Context, tt *models.Todo, input map[string]any) (*models.Todo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchUpdate", ctx, tt, input)
	ret0, _ := ret[0].(*models.Todo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchUpdate indicates an expected call of PatchUpdate.
func (mr *MockDaoMockRecorder) PatchUpdate(ctx, tt, input any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchUpdate", reflect.TypeOf((*MockDao)(nil).PatchUpdate), ctx, tt, input)
}
