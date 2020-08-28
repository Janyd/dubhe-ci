// Code generated by MockGen. DO NOT EDIT.
// Source: dubhe-ci/core (interfaces: RepositoryStore,BuildStore,BranchStore,CredentialStore,StepStore,Scheduler,ConfigService,ConvertService,Manager)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	common "dubhe-ci/common"
	core "dubhe-ci/core"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepositoryStore is a mock of RepositoryStore interface
type MockRepositoryStore struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryStoreMockRecorder
}

// MockRepositoryStoreMockRecorder is the mock recorder for MockRepositoryStore
type MockRepositoryStoreMockRecorder struct {
	mock *MockRepositoryStore
}

// NewMockRepositoryStore creates a new mock instance
func NewMockRepositoryStore(ctrl *gomock.Controller) *MockRepositoryStore {
	mock := &MockRepositoryStore{ctrl: ctrl}
	mock.recorder = &MockRepositoryStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepositoryStore) EXPECT() *MockRepositoryStoreMockRecorder {
	return m.recorder
}

// CheckName mocks base method
func (m *MockRepositoryStore) CheckName(arg0 context.Context, arg1 int64, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckName", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckName indicates an expected call of CheckName
func (mr *MockRepositoryStoreMockRecorder) CheckName(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckName", reflect.TypeOf((*MockRepositoryStore)(nil).CheckName), arg0, arg1, arg2)
}

// Create mocks base method
func (m *MockRepositoryStore) Create(arg0 context.Context, arg1 *core.Repository) (*core.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*core.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockRepositoryStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepositoryStore)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockRepositoryStore) Delete(arg0 context.Context, arg1 *core.Repository) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepositoryStore)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockRepositoryStore) Find(arg0 context.Context, arg1 int64) (*core.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockRepositoryStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockRepositoryStore)(nil).Find), arg0, arg1)
}

// Increment mocks base method
func (m *MockRepositoryStore) Increment(arg0 context.Context, arg1 *core.Repository) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Increment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Increment indicates an expected call of Increment
func (mr *MockRepositoryStoreMockRecorder) Increment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increment", reflect.TypeOf((*MockRepositoryStore)(nil).Increment), arg0, arg1)
}

// List mocks base method
func (m *MockRepositoryStore) List(arg0 context.Context, arg1 *common.Page) (*common.Page, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*common.Page)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockRepositoryStoreMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepositoryStore)(nil).List), arg0, arg1)
}

// Update mocks base method
func (m *MockRepositoryStore) Update(arg0 context.Context, arg1 *core.Repository) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockRepositoryStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepositoryStore)(nil).Update), arg0, arg1)
}

// MockBuildStore is a mock of BuildStore interface
type MockBuildStore struct {
	ctrl     *gomock.Controller
	recorder *MockBuildStoreMockRecorder
}

// MockBuildStoreMockRecorder is the mock recorder for MockBuildStore
type MockBuildStoreMockRecorder struct {
	mock *MockBuildStore
}

// NewMockBuildStore creates a new mock instance
func NewMockBuildStore(ctrl *gomock.Controller) *MockBuildStore {
	mock := &MockBuildStore{ctrl: ctrl}
	mock.recorder = &MockBuildStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBuildStore) EXPECT() *MockBuildStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockBuildStore) Create(arg0 context.Context, arg1 *core.Build, arg2 []*core.Step) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockBuildStoreMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBuildStore)(nil).Create), arg0, arg1, arg2)
}

// Delete mocks base method
func (m *MockBuildStore) Delete(arg0 context.Context, arg1 *core.Build) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockBuildStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBuildStore)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockBuildStore) Find(arg0 context.Context, arg1 int64) (*core.Build, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Build)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockBuildStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockBuildStore)(nil).Find), arg0, arg1)
}

// FindNumber mocks base method
func (m *MockBuildStore) FindNumber(arg0 context.Context, arg1 int64, arg2 uint32) (*core.Build, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindNumber", arg0, arg1, arg2)
	ret0, _ := ret[0].(*core.Build)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindNumber indicates an expected call of FindNumber
func (mr *MockBuildStoreMockRecorder) FindNumber(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindNumber", reflect.TypeOf((*MockBuildStore)(nil).FindNumber), arg0, arg1, arg2)
}

// List mocks base method
func (m *MockBuildStore) List(arg0 context.Context, arg1 *common.Page) (*common.Page, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*common.Page)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockBuildStoreMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockBuildStore)(nil).List), arg0, arg1)
}

// ListIncomplete mocks base method
func (m *MockBuildStore) ListIncomplete(arg0 context.Context) ([]*core.Build, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListIncomplete", arg0)
	ret0, _ := ret[0].([]*core.Build)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIncomplete indicates an expected call of ListIncomplete
func (mr *MockBuildStoreMockRecorder) ListIncomplete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIncomplete", reflect.TypeOf((*MockBuildStore)(nil).ListIncomplete), arg0)
}

// Update mocks base method
func (m *MockBuildStore) Update(arg0 context.Context, arg1 *core.Build) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockBuildStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockBuildStore)(nil).Update), arg0, arg1)
}

// MockBranchStore is a mock of BranchStore interface
type MockBranchStore struct {
	ctrl     *gomock.Controller
	recorder *MockBranchStoreMockRecorder
}

// MockBranchStoreMockRecorder is the mock recorder for MockBranchStore
type MockBranchStoreMockRecorder struct {
	mock *MockBranchStore
}

// NewMockBranchStore creates a new mock instance
func NewMockBranchStore(ctrl *gomock.Controller) *MockBranchStore {
	mock := &MockBranchStore{ctrl: ctrl}
	mock.recorder = &MockBranchStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBranchStore) EXPECT() *MockBranchStoreMockRecorder {
	return m.recorder
}

// Activate mocks base method
func (m *MockBranchStore) Activate(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Activate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Activate indicates an expected call of Activate
func (mr *MockBranchStoreMockRecorder) Activate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Activate", reflect.TypeOf((*MockBranchStore)(nil).Activate), arg0, arg1)
}

// Create mocks base method
func (m *MockBranchStore) Create(arg0 context.Context, arg1 int64, arg2 string) (*core.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(*core.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockBranchStoreMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBranchStore)(nil).Create), arg0, arg1, arg2)
}

// Delete mocks base method
func (m *MockBranchStore) Delete(arg0 context.Context, arg1 *core.Branch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockBranchStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBranchStore)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockBranchStore) Find(arg0 context.Context, arg1 int64) (*core.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockBranchStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockBranchStore)(nil).Find), arg0, arg1)
}

// FindByName mocks base method
func (m *MockBranchStore) FindByName(arg0 context.Context, arg1 int64, arg2 string) (*core.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByName", arg0, arg1, arg2)
	ret0, _ := ret[0].(*core.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByName indicates an expected call of FindByName
func (mr *MockBranchStoreMockRecorder) FindByName(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByName", reflect.TypeOf((*MockBranchStore)(nil).FindByName), arg0, arg1, arg2)
}

// Increment mocks base method
func (m *MockBranchStore) Increment(arg0 context.Context, arg1 *core.Branch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Increment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Increment indicates an expected call of Increment
func (mr *MockBranchStoreMockRecorder) Increment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Increment", reflect.TypeOf((*MockBranchStore)(nil).Increment), arg0, arg1)
}

// List mocks base method
func (m *MockBranchStore) List(arg0 context.Context, arg1 int64) ([]*core.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*core.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockBranchStoreMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockBranchStore)(nil).List), arg0, arg1)
}

// Update mocks base method
func (m *MockBranchStore) Update(arg0 context.Context, arg1 *core.Branch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockBranchStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockBranchStore)(nil).Update), arg0, arg1)
}

// MockCredentialStore is a mock of CredentialStore interface
type MockCredentialStore struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialStoreMockRecorder
}

// MockCredentialStoreMockRecorder is the mock recorder for MockCredentialStore
type MockCredentialStoreMockRecorder struct {
	mock *MockCredentialStore
}

// NewMockCredentialStore creates a new mock instance
func NewMockCredentialStore(ctrl *gomock.Controller) *MockCredentialStore {
	mock := &MockCredentialStore{ctrl: ctrl}
	mock.recorder = &MockCredentialStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCredentialStore) EXPECT() *MockCredentialStoreMockRecorder {
	return m.recorder
}

// CheckName mocks base method
func (m *MockCredentialStore) CheckName(arg0 context.Context, arg1 int64, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckName", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckName indicates an expected call of CheckName
func (mr *MockCredentialStoreMockRecorder) CheckName(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckName", reflect.TypeOf((*MockCredentialStore)(nil).CheckName), arg0, arg1, arg2)
}

// Create mocks base method
func (m *MockCredentialStore) Create(arg0 context.Context, arg1 *core.Credential) (*core.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*core.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockCredentialStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCredentialStore)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockCredentialStore) Delete(arg0 context.Context, arg1 *core.Credential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockCredentialStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCredentialStore)(nil).Delete), arg0, arg1)
}

// Find mocks base method
func (m *MockCredentialStore) Find(arg0 context.Context, arg1 int64) (*core.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockCredentialStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockCredentialStore)(nil).Find), arg0, arg1)
}

// List mocks base method
func (m *MockCredentialStore) List(arg0 context.Context) ([]*core.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*core.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockCredentialStoreMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCredentialStore)(nil).List), arg0)
}

// ListRegistryCred mocks base method
func (m *MockCredentialStore) ListRegistryCred(arg0 context.Context) ([]*core.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRegistryCred", arg0)
	ret0, _ := ret[0].([]*core.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRegistryCred indicates an expected call of ListRegistryCred
func (mr *MockCredentialStoreMockRecorder) ListRegistryCred(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRegistryCred", reflect.TypeOf((*MockCredentialStore)(nil).ListRegistryCred), arg0)
}

// Update mocks base method
func (m *MockCredentialStore) Update(arg0 context.Context, arg1 *core.Credential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockCredentialStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCredentialStore)(nil).Update), arg0, arg1)
}

// MockStepStore is a mock of StepStore interface
type MockStepStore struct {
	ctrl     *gomock.Controller
	recorder *MockStepStoreMockRecorder
}

// MockStepStoreMockRecorder is the mock recorder for MockStepStore
type MockStepStoreMockRecorder struct {
	mock *MockStepStore
}

// NewMockStepStore creates a new mock instance
func NewMockStepStore(ctrl *gomock.Controller) *MockStepStore {
	mock := &MockStepStore{ctrl: ctrl}
	mock.recorder = &MockStepStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStepStore) EXPECT() *MockStepStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockStepStore) Create(arg0 context.Context, arg1 []*core.Step) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockStepStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStepStore)(nil).Create), arg0, arg1)
}

// Find mocks base method
func (m *MockStepStore) Find(arg0 context.Context, arg1 int64) (*core.Step, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Step)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockStepStoreMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockStepStore)(nil).Find), arg0, arg1)
}

// FindNumber mocks base method
func (m *MockStepStore) FindNumber(arg0 context.Context, arg1 int64, arg2 uint32) (*core.Step, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindNumber", arg0, arg1, arg2)
	ret0, _ := ret[0].(*core.Step)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindNumber indicates an expected call of FindNumber
func (mr *MockStepStoreMockRecorder) FindNumber(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindNumber", reflect.TypeOf((*MockStepStore)(nil).FindNumber), arg0, arg1, arg2)
}

// List mocks base method
func (m *MockStepStore) List(arg0 context.Context, arg1 int64) ([]*core.Step, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*core.Step)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockStepStoreMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockStepStore)(nil).List), arg0, arg1)
}

// Update mocks base method
func (m *MockStepStore) Update(arg0 context.Context, arg1 *core.Step) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockStepStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStepStore)(nil).Update), arg0, arg1)
}

// MockScheduler is a mock of Scheduler interface
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// Cancel mocks base method
func (m *MockScheduler) Cancel(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Cancel indicates an expected call of Cancel
func (mr *MockSchedulerMockRecorder) Cancel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockScheduler)(nil).Cancel), arg0, arg1)
}

// Cancelled mocks base method
func (m *MockScheduler) Cancelled(arg0 context.Context, arg1 int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancelled", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Cancelled indicates an expected call of Cancelled
func (mr *MockSchedulerMockRecorder) Cancelled(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancelled", reflect.TypeOf((*MockScheduler)(nil).Cancelled), arg0, arg1)
}

// Pause mocks base method
func (m *MockScheduler) Pause(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Pause", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Pause indicates an expected call of Pause
func (mr *MockSchedulerMockRecorder) Pause(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Pause", reflect.TypeOf((*MockScheduler)(nil).Pause), arg0)
}

// Request mocks base method
func (m *MockScheduler) Request(arg0 context.Context) (*core.Build, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Request", arg0)
	ret0, _ := ret[0].(*core.Build)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Request indicates an expected call of Request
func (mr *MockSchedulerMockRecorder) Request(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Request", reflect.TypeOf((*MockScheduler)(nil).Request), arg0)
}

// Resume mocks base method
func (m *MockScheduler) Resume(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resume", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Resume indicates an expected call of Resume
func (mr *MockSchedulerMockRecorder) Resume(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resume", reflect.TypeOf((*MockScheduler)(nil).Resume), arg0)
}

// Schedule mocks base method
func (m *MockScheduler) Schedule(arg0 context.Context, arg1 *core.Build) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Schedule", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Schedule indicates an expected call of Schedule
func (mr *MockSchedulerMockRecorder) Schedule(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Schedule", reflect.TypeOf((*MockScheduler)(nil).Schedule), arg0, arg1)
}

// MockConfigService is a mock of ConfigService interface
type MockConfigService struct {
	ctrl     *gomock.Controller
	recorder *MockConfigServiceMockRecorder
}

// MockConfigServiceMockRecorder is the mock recorder for MockConfigService
type MockConfigServiceMockRecorder struct {
	mock *MockConfigService
}

// NewMockConfigService creates a new mock instance
func NewMockConfigService(ctrl *gomock.Controller) *MockConfigService {
	mock := &MockConfigService{ctrl: ctrl}
	mock.recorder = &MockConfigServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigService) EXPECT() *MockConfigServiceMockRecorder {
	return m.recorder
}

// Find mocks base method
func (m *MockConfigService) Find(arg0 context.Context, arg1 *core.ConfigArgs) (*core.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*core.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find
func (mr *MockConfigServiceMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockConfigService)(nil).Find), arg0, arg1)
}

// MockConvertService is a mock of ConvertService interface
type MockConvertService struct {
	ctrl     *gomock.Controller
	recorder *MockConvertServiceMockRecorder
}

// MockConvertServiceMockRecorder is the mock recorder for MockConvertService
type MockConvertServiceMockRecorder struct {
	mock *MockConvertService
}

// NewMockConvertService creates a new mock instance
func NewMockConvertService(ctrl *gomock.Controller) *MockConvertService {
	mock := &MockConvertService{ctrl: ctrl}
	mock.recorder = &MockConvertServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConvertService) EXPECT() *MockConvertServiceMockRecorder {
	return m.recorder
}

// Convert mocks base method
func (m *MockConvertService) Convert(arg0 context.Context, arg1 *core.ConvertArgs) (*core.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Convert", arg0, arg1)
	ret0, _ := ret[0].(*core.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Convert indicates an expected call of Convert
func (mr *MockConvertServiceMockRecorder) Convert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Convert", reflect.TypeOf((*MockConvertService)(nil).Convert), arg0, arg1)
}

// MockManager is a mock of Manager interface
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Accept mocks base method
func (m *MockManager) Accept(arg0 context.Context, arg1 int64) (*core.Build, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accept", arg0, arg1)
	ret0, _ := ret[0].(*core.Build)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Accept indicates an expected call of Accept
func (mr *MockManagerMockRecorder) Accept(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockManager)(nil).Accept), arg0, arg1)
}

// Details mocks base method
func (m *MockManager) Details(arg0 context.Context, arg1 int64) (*core.Context, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Details", arg0, arg1)
	ret0, _ := ret[0].(*core.Context)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Details indicates an expected call of Details
func (mr *MockManagerMockRecorder) Details(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Details", reflect.TypeOf((*MockManager)(nil).Details), arg0, arg1)
}

// Request mocks base method
func (m *MockManager) Request(arg0 context.Context) (*core.Build, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Request", arg0)
	ret0, _ := ret[0].(*core.Build)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Request indicates an expected call of Request
func (mr *MockManagerMockRecorder) Request(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Request", reflect.TypeOf((*MockManager)(nil).Request), arg0)
}

// Watch mocks base method
func (m *MockManager) Watch(arg0 context.Context, arg1 int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch
func (mr *MockManagerMockRecorder) Watch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockManager)(nil).Watch), arg0, arg1)
}
