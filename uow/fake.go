package uow

type FakeUnitOfWork struct {
	batches   Repository
	Committed bool
}

func NewFakeUnitOfWork(batches Repository) *FakeUnitOfWork {
	return &FakeUnitOfWork{batches: batches}
}

func (f *FakeUnitOfWork) Commit() error {
	f.Committed = true
	return nil
}

func (f *FakeUnitOfWork) Rollback() {

}

func (f *FakeUnitOfWork) Batches() Repository {
	return f.batches
}

func (f *FakeUnitOfWork) DBInstruction(queryFunction QueryFunc) error {
	if err := queryFunction(); err != nil {
		return err
	}
	f.Commit()
	return nil
}
