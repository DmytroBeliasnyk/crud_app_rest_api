// Code generated by MockGen. DO NOT EDIT.
// Source: abstract.go

// Package mock_services is a generated GoMock package.
package mock_services

import (
	reflect "reflect"

	dto "github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"
	gomock "github.com/golang/mock/gomock"
)

// MockProjectService is a mock of ProjectService interface.
type MockProjectService struct {
	ctrl     *gomock.Controller
	recorder *MockProjectServiceMockRecorder
}

// MockProjectServiceMockRecorder is the mock recorder for MockProjectService.
type MockProjectServiceMockRecorder struct {
	mock *MockProjectService
}

// NewMockProjectService creates a new mock instance.
func NewMockProjectService(ctrl *gomock.Controller) *MockProjectService {
	mock := &MockProjectService{ctrl: ctrl}
	mock.recorder = &MockProjectServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProjectService) EXPECT() *MockProjectServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProjectService) Create(p dto.ProjectDTO, userId int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", p, userId)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockProjectServiceMockRecorder) Create(p, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProjectService)(nil).Create), p, userId)
}

// DeleteById mocks base method.
func (m *MockProjectService) DeleteById(id, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", id, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockProjectServiceMockRecorder) DeleteById(id, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockProjectService)(nil).DeleteById), id, userId)
}

// GetAll mocks base method.
func (m *MockProjectService) GetAll(userId int64) ([]dto.ProjectDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", userId)
	ret0, _ := ret[0].([]dto.ProjectDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockProjectServiceMockRecorder) GetAll(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockProjectService)(nil).GetAll), userId)
}

// GetById mocks base method.
func (m *MockProjectService) GetById(id, userId int64) (dto.ProjectDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id, userId)
	ret0, _ := ret[0].(dto.ProjectDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockProjectServiceMockRecorder) GetById(id, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockProjectService)(nil).GetById), id, userId)
}

// UpdateById mocks base method.
func (m *MockProjectService) UpdateById(id int64, p dto.UpdateProjectDTO, userId int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateById", id, p, userId)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateById indicates an expected call of UpdateById.
func (mr *MockProjectServiceMockRecorder) UpdateById(id, p, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateById", reflect.TypeOf((*MockProjectService)(nil).UpdateById), id, p, userId)
}

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// GenerateTokens mocks base method.
func (m *MockAuthService) GenerateTokens(id int64) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokens", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GenerateTokens indicates an expected call of GenerateTokens.
func (mr *MockAuthServiceMockRecorder) GenerateTokens(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokens", reflect.TypeOf((*MockAuthService)(nil).GenerateTokens), id)
}

// HashPassword mocks base method.
func (m *MockAuthService) HashPassword(password string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HashPassword", password)
	ret0, _ := ret[0].(string)
	return ret0
}

// HashPassword indicates an expected call of HashPassword.
func (mr *MockAuthServiceMockRecorder) HashPassword(password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HashPassword", reflect.TypeOf((*MockAuthService)(nil).HashPassword), password)
}

// ParseToken mocks base method.
func (m *MockAuthService) ParseToken(input string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", input)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockAuthServiceMockRecorder) ParseToken(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockAuthService)(nil).ParseToken), input)
}

// SignIn mocks base method.
func (m *MockAuthService) SignIn(si dto.SignInDTO) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", si)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthServiceMockRecorder) SignIn(si interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthService)(nil).SignIn), si)
}

// SignUp mocks base method.
func (m *MockAuthService) SignUp(su dto.SignUpDTO) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", su)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthServiceMockRecorder) SignUp(su interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthService)(nil).SignUp), su)
}

// UpdateTokens mocks base method.
func (m *MockAuthService) UpdateTokens(rt string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTokens", rt)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateTokens indicates an expected call of UpdateTokens.
func (mr *MockAuthServiceMockRecorder) UpdateTokens(rt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTokens", reflect.TypeOf((*MockAuthService)(nil).UpdateTokens), rt)
}
