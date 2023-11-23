package uow

type FakeUnitOfWork struct {
	products  ProductRepository
	Committed bool
}

func NewFakeUnitOfWork(products ProductRepository) *FakeUnitOfWork {
	return &FakeUnitOfWork{products: products}
}

func (f *FakeUnitOfWork) Commit() error {
	f.Committed = true
	return nil
}

func (f *FakeUnitOfWork) Rollback() {

}

func (f *FakeUnitOfWork) Products() ProductRepository {
	return f.products
}

func (f *FakeUnitOfWork) CommitOnSuccess(queryFunction QueryFunc) error {
	if err := queryFunction(); err != nil {
		return err
	}
	f.Commit()
	return nil
}
