package mocks

import (
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/port/repository/psql"
	"github.com/stretchr/testify/mock"
)

type MockPendingFranchiseRepository struct {
	mock.Mock
	psql.PendingFranchiseRepository
}

func (m *MockPendingFranchiseRepository) Put(n dto.PendingFranchiseDTO) (dto.PendingFranchiseDTO, error) {
	args := m.Called()
	return args.Get(0).(dto.PendingFranchiseDTO), args.Error(1)
}
