package service

import (
	"fmt"

	"github.com/RomanshkVolkov/test-api/internal/adapters/repository"
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	schema "github.com/RomanshkVolkov/test-api/internal/core/domain/schemas"
)

func (server Server) GetDocuments() domain.APIResponse[[]domain.Document, any] {
	repo := repository.GetDBConnection(server.Host)
	documents, err := repo.GetDocuments()

	if err != nil {
		return domain.APIResponse[[]domain.Document, any]{
			Message: domain.Message{En: "error on get documents", Es: "error al obtener documentos"},
			Error:   err}
	}

	return domain.APIResponse[[]domain.Document, any]{
		Success: true,
		Message: domain.Message{En: "Documents retrieved", Es: "Documentos recuperados"},
		Data:    documents,
	}
}

func (server Server) GetDocumentByID(id uint) domain.APIResponse[domain.DocumentWithDetails, any] {
	repo := repository.GetDBConnection(server.Host)
	document, err := repo.GetDocumentByID(id)

	if err != nil {
		return repository.RecordNotFound[domain.DocumentWithDetails]()
	}

	return domain.APIResponse[domain.DocumentWithDetails, any]{
		Success: true,
		Message: domain.Message{En: "Document retrieved", Es: "Documento recuperado"},
		Data:    document,
	}
}

func (server Server) CreateDocument(request *domain.DocumentWithDetails) domain.APIResponse[domain.DocumentWithDetails, any] {
	fields := schema.GenericForm[domain.DocumentWithDetails]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)
	if len(failValidatedFields) > 0 {
		return domain.APIResponse[domain.DocumentWithDetails, any]{SchemaError: failValidatedFields}
	}

	repo := repository.GetDBConnection(server.Host)
	createdDocument, err := repo.CreateDocument(fields.Data)

	fmt.Println("print error: ", err.Error())

	if err.Error() == "mssql: RECORD_ALREADY_EXIST" {
		return domain.APIResponse[domain.DocumentWithDetails, any]{
			Message:     domain.Message{En: "Document already exist", Es: "Documento ya existe"},
			SchemaError: domain.ObjectString{"table": "La tabla ya existe en la base de datos"},
			Error:       err}
	}

	if err != nil {
		return domain.APIResponse[domain.DocumentWithDetails, any]{
			Message: domain.Message{En: "Oops, an error has ocurred in the repository", Es: "Oops, un error ha ocurrido en el repositorio"},
			Error:   err}
	}

	return domain.APIResponse[domain.DocumentWithDetails, any]{
		Success: true,
		Message: domain.Message{En: "Document created", Es: "Documento creado"},
		Data:    createdDocument,
	}
}

func (server Server) GetTables() domain.APIResponse[[]domain.Document, any] {
	repo := repository.GetDBConnection(server.Host)
	tables, err := repo.GetTables()
	if err != nil {
		return domain.APIResponse[[]domain.Document, any]{Error: err}
	}

	return domain.APIResponse[[]domain.Document, any]{
		Success: true,
		Message: domain.Message{En: "Reports retrieved", Es: "Reportes recuperados"},
		Data:    tables,
	}
}

func (server Server) GetTableByID(id uint) domain.APIResponse[domain.TableByID, any] {
	repo := repository.GetDBConnection(server.Host)
	report, err := repo.GetTableByID(id)
	if err != nil {
		return domain.APIResponse[domain.TableByID, any]{Error: err}
	}

	return domain.APIResponse[domain.TableByID, any]{
		Success: true,
		Message: domain.Message{En: "Table retrieved", Es: "Tabla recuperada"},
		Data:    report,
	}
}

func (server Server) UploadDocument(request *domain.UploadDocument) domain.APIResponse[any, any] {
	fields := schema.GenericForm[domain.UploadDocument]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)
	if len(failValidatedFields) > 0 {
		return domain.APIResponse[any, any]{SchemaError: failValidatedFields}
	}

	repo := repository.GetDBConnection(server.Host)
	err := repo.UploadDocument(fields.Data.DocumentID, fields.Data.File)
	if err != nil {
		return domain.APIResponse[any, any]{
			Message: domain.Message{
				En: "error on upload document",
				Es: "error al subir documento",
			},
			Error: err}
	}

	return domain.APIResponse[any, any]{
		Success: true,
		Message: domain.Message{En: "Document uploaded", Es: "Documento subido"},
		Data:    fields.Data,
	}
}
