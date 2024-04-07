package services

type DbServiceLinkinterface interface {
}

type DatabaseServiceInterface interface {
	CreateNewLink() DbServiceLinkinterface
}
