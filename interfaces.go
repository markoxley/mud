package dtorm

type Updater interface {
	Update(mgr Manager)
}

type Restorer interface {
	Restore(mgr Manager)
}

type Remover interface{}
