package usecase

import (
	"testing"

	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/queue"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/common/queue/event"
	"github.com/kaitsubaka/clubhub_franchises/internal/core/dto"
	"github.com/kaitsubaka/clubhub_franchises/internal/infra/repository/psql/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type sub struct{}

func (s *sub) Subscribe(e event.Event) error {
	return nil
}
func TestPendingFranchiseUseCase_Create(t *testing.T) {
	q, _ := queue.New(1).WithNumOfWorker(1).WithSubscriber(&sub{}).Build()
	defer q.Close()
	mockPendingFranchiseRepository := new(mocks.MockPendingFranchiseRepository)
	nfs := &PendingFranchiseUseCase{
		newFranchizeRepository:     mockPendingFranchiseRepository,
		franchiseCreationPublisher: q,
	}
	t.Run("it should return nil error and the created pending franchise", func(t *testing.T) {
		mockPendingFranchiseRepository.On("Put", mock.Anything).Return(dto.PendingFranchiseDTO{}, nil)
		got, err := nfs.Create(dto.CreateNewFranchizeDTO{})
		require.NoError(t, err)
		require.Equal(t, dto.PendingFranchiseDTO{}, got)
	})

}
