// Code generated by MockGen. DO NOT EDIT.
// Source: ./user.go
//
// Generated by this command:
//
//	mockgen -source=./user.go -destination=./mock/user.go
//

// Package mock_user is a generated GoMock package.
package mock_user

import (
	reflect "reflect"

	user "github.com/GianOrtiz/bean/pkg/user"
	sessions "github.com/gorilla/sessions"
	gomock "go.uber.org/mock/gomock"
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

// Create mocks base method.
func (m *MockRepository) Create(u *user.User, passwordHash string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", u, passwordHash)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(u, passwordHash any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), u, passwordHash)
}

// GetPasswordHash mocks base method.
func (m *MockRepository) GetPasswordHash(id int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPasswordHash", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPasswordHash indicates an expected call of GetPasswordHash.
func (mr *MockRepositoryMockRecorder) GetPasswordHash(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPasswordHash", reflect.TypeOf((*MockRepository)(nil).GetPasswordHash), id)
}

// GetUser mocks base method.
func (m *MockRepository) GetUser(id int) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockRepositoryMockRecorder) GetUser(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockRepository)(nil).GetUser), id)
}

// GetUserByEmail mocks base method.
func (m *MockRepository) GetUserByEmail(email string) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockRepositoryMockRecorder) GetUserByEmail(email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockRepository)(nil).GetUserByEmail), email)
}

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// GetUser mocks base method.
func (m *MockUseCase) GetUser(id int) (*user.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*user.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUseCaseMockRecorder) GetUser(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUseCase)(nil).GetUser), id)
}

// Login mocks base method.
func (m *MockUseCase) Login(email, password string) (*sessions.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", email, password)
	ret0, _ := ret[0].(*sessions.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUseCaseMockRecorder) Login(email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUseCase)(nil).Login), email, password)
}

// Register mocks base method.
func (m *MockUseCase) Register(email, name, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", email, name, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUseCaseMockRecorder) Register(email, name, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUseCase)(nil).Register), email, name, password)
}
