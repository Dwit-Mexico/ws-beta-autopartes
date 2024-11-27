package service

import (
	"github.com/RomanshkVolkov/test-api/internal/adapters/repository"
	"github.com/RomanshkVolkov/test-api/internal/core/domain"
	schema "github.com/RomanshkVolkov/test-api/internal/core/domain/schemas"
)

func (server Server) GetDocuments() domain.APIResponse[[]domain.Document] {
	repo := repository.GetDBConnection(server.Host)
	documents, err := repo.GetDocuments()

	if err != nil {
		return domain.APIResponse[[]domain.Document]{
			Message: domain.Message{En: "error on get documents", Es: "error al obtener documentos"},
			Error:   err}
	}

	return domain.APIResponse[[]domain.Document]{
		Success: true,
		Message: domain.Message{En: "Documents retrieved", Es: "Documentos recuperados"},
		Data:    documents,
	}
}

func (server Server) GetDocumentByID(id uint) domain.APIResponse[domain.DocumentWithDetails] {
	repo := repository.GetDBConnection(server.Host)
	document, err := repo.GetDocumentByID(id)

	if err != nil {
		return repository.RecordNotFound[domain.DocumentWithDetails]()
	}

	return domain.APIResponse[domain.DocumentWithDetails]{
		Success: true,
		Message: domain.Message{En: "Document retrieved", Es: "Documento recuperado"},
		Data:    document,
	}
}

func (server Server) CreateDocument(request *domain.DocumentWithDetails) domain.APIResponse[domain.DocumentWithDetails] {
	fields := schema.GenericForm[domain.DocumentWithDetails]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)
	if len(failValidatedFields) > 0 {
		return domain.APIResponse[domain.DocumentWithDetails]{SchemaError: failValidatedFields}
	}

	repo := repository.GetDBConnection(server.Host)
	createdDocument, err := repo.CreateDocument(fields.Data)

	if err != nil {
		if err.Error() == "mssql: RECORD_ALREADY_EXIST" {
			return domain.APIResponse[domain.DocumentWithDetails]{
				Message:     domain.Message{En: "Document already exist", Es: "Documento ya existe"},
				SchemaError: domain.ObjectString{"table": []string{"La tabla ya existe en la base de datos"}},
				Error:       err}
		}
		return domain.APIResponse[domain.DocumentWithDetails]{
			Message: domain.Message{En: "Oops, an error has ocurred in the repository", Es: "Oops, un error ha ocurrido en el repositorio"},
			Error:   err}
	}

	return domain.APIResponse[domain.DocumentWithDetails]{
		Success: true,
		Message: domain.Message{En: "Document created", Es: "Documento creado"},
		Data:    createdDocument,
	}
}

func (server Server) UpdateDocument(request *domain.EditableDocument) domain.APIResponse[domain.DocumentWithDetails] {
	fields := schema.GenericForm[domain.EditableDocument]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)
	if len(failValidatedFields) > 0 {
		return domain.APIResponse[domain.DocumentWithDetails]{SchemaError: failValidatedFields}
	}

	repo := repository.GetDBConnection(server.Host)
	document, err := repo.UpdateDocument(fields.Data)

	if err != nil {
		return repository.HandleDatabaseError[domain.DocumentWithDetails](err, domain.Message{En: "Error on update document", Es: "Error al actualizar documento"})
	}

	if document.ID == 0 {
		return repository.RecordNotFound[domain.DocumentWithDetails]()
	}
	return domain.APIResponse[domain.DocumentWithDetails]{
		Success: true,
		Message: domain.Message{En: "Document updated", Es: "Documento actualizado"},
		Data:    document,
	}
}

func (server Server) DeleteDocument(id uint) domain.APIResponse[any] {
	repo := repository.GetDBConnection(server.Host)
	err := repo.DeleteDocument(id)
	if err != nil {
		return domain.APIResponse[any]{
			Message: domain.Message{En: "Error on delete document", Es: "Error al eliminar documento"},
			Error:   err}
	}

	return domain.APIResponse[any]{
		Success: true,
		Message: domain.Message{En: "Document deleted", Es: "Documento eliminado"},
	}
}

func (server Server) DeleteFieldDocument(id uint) domain.APIResponse[any] {
	repo := repository.GetDBConnection(server.Host)
	err := repo.DeleteFieldDocument(id)
	if err != nil {
		return domain.APIResponse[any]{
			Message: domain.Message{En: "Error on delete document", Es: "Error al eliminar documento"},
			Error:   err}
	}

	return domain.APIResponse[any]{
		Success: true,
		Message: domain.Message{En: "Document deleted", Es: "Documento eliminado"},
	}
}

func (server Server) GetDocumentRowRecord(id uint, documentID uint) domain.APIResponse[[]domain.DocumentRowRecord] {
	repo := repository.GetDBConnection(server.Host)
	records, err := repo.GetDocumentRowRecord(id, documentID)
	if err != nil {
		return domain.APIResponse[[]domain.DocumentRowRecord]{Error: err}
	}

	return domain.APIResponse[[]domain.DocumentRowRecord]{
		Success: true,
		Message: domain.Message{En: "Records retrieved", Es: "Registros recuperados"},
		Data:    records,
	}
}

func (server Server) GetReports() domain.APIResponse[[]domain.Report] {
	repo := repository.GetDBConnection(server.Host)
	reports, err := repo.GetReports()
	if err != nil {
		return domain.APIResponse[[]domain.Report]{Error: err}
	}

	return domain.APIResponse[[]domain.Report]{
		Success: true,
		Message: domain.Message{En: "Reports retrieved", Es: "Reportes recuperados"},
		Data:    reports,
	}
}

func (server Server) GetReportByID(id uint) domain.APIResponse[domain.ReportByID] {
	repo := repository.GetDBConnection(server.Host)
	report, err := repo.GetReportByID(id)
	if err != nil {
		return domain.APIResponse[domain.ReportByID]{Error: err}
	}

	return domain.APIResponse[domain.ReportByID]{
		Success: true,
		Message: domain.Message{En: "Report retrieved", Es: "Reporte recuperado"},
		Data:    report,
	}
}

func (server Server) UpdateDocumentRowRecord(request *domain.EditableDoumentRowRecord) domain.APIResponse[any] {
	fields := schema.GenericForm[domain.EditableDoumentRowRecord]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)
	if len(failValidatedFields) > 0 {
		return domain.APIResponse[any]{SchemaError: failValidatedFields}
	}

	repo := repository.GetDBConnection(server.Host)
	err := repo.UpdateDocumentRowRecord(fields.Data)
	if err != nil {
		return domain.APIResponse[any]{Error: err}
	}

	return domain.APIResponse[any]{
		Success: true,
		Message: domain.Message{En: "Record updated", Es: "Registro actualizado"},
	}
}

func (server Server) DeleteDocumentRowRecord(id uint, documentID uint) domain.APIResponse[any] {
	repo := repository.GetDBConnection(server.Host)
	err := repo.DeleteDocumentRowRecord(id, documentID)
	if err != nil {
		return domain.APIResponse[any]{
			Message: domain.Message{En: "Error on delete record", Es: "Error al eliminar registro"},
			Error:   err}
	}

	return domain.APIResponse[any]{
		Success: true,
		Message: domain.Message{En: "Record deleted", Es: "Registro eliminado"},
	}
}

func (server Server) GetTables() domain.APIResponse[[]domain.Document] {
	repo := repository.GetDBConnection(server.Host)
	tables, err := repo.GetTables()
	if err != nil {
		return domain.APIResponse[[]domain.Document]{Error: err}
	}

	return domain.APIResponse[[]domain.Document]{
		Success: true,
		Message: domain.Message{En: "Reports retrieved", Es: "Reportes recuperados"},
		Data:    tables,
	}
}

func (server Server) GetTableByID(id uint) domain.APIResponse[domain.TableByID] {
	repo := repository.GetDBConnection(server.Host)
	report, err := repo.GetTableByID(id)
	if err != nil {
		return domain.APIResponse[domain.TableByID]{Error: err}
	}

	return domain.APIResponse[domain.TableByID]{
		Success: true,
		Message: domain.Message{En: "Table retrieved", Es: "Tabla recuperada"},
		Data:    report,
	}
}

func (server Server) UploadDocument(request *domain.UploadDocument) domain.APIResponse[any] {
	fields := schema.GenericForm[domain.UploadDocument]{Data: *request}
	failValidatedFields := schema.FormValidator(fields)
	if len(failValidatedFields) > 0 {
		return domain.APIResponse[any]{SchemaError: failValidatedFields}
	}

	repo := repository.GetDBConnection(server.Host)
	err := repo.UploadDocument(fields.Data.DocumentID, fields.Data.File)
	if err != nil {
		return domain.APIResponse[any]{
			Message: domain.Message{
				En: "error on upload document",
				Es: "error al subir documento",
			},
			Error: err}
	}

	return domain.APIResponse[any]{
		Success: true,
		Message: domain.Message{En: "Document uploaded", Es: "Documento subido"},
		Data:    fields.Data,
	}
}
