package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/RomanshkVolkov/test-api/internal/core/domain"
)

func (database *DSNSource) GetDocuments() ([]domain.Document, error) {
	documents := []domain.Document{}
	err := database.DB.Model(&domain.Document{}).Find(&documents).Error

	if err != nil {
		return []domain.Document{}, err
	}

	return documents, nil
}

func (database *DSNSource) GetDocumentByID(id uint) (domain.DocumentWithDetails, error) {
	document := domain.Document{}

	database.DB.Model(&domain.Document{}).Where("id = ?", id).First(&document)

	details := []domain.DetailDocument{}

	database.DB.Model(&domain.DetailDocument{}).Where("document_id = ?", id).Find(&details)

	documentWithDetails := domain.DocumentWithDetails{
		Document: document,
		Details:  details,
	}

	if document.ID == 0 {
		return domain.DocumentWithDetails{}, errors.New("record not found")
	}

	return documentWithDetails, nil
}

func (database *DSNSource) CreateDocument(document domain.DocumentWithDetails) (domain.DocumentWithDetails, error) {
	tx := database.DB.Begin()

	document.Name = Capitalize(document.Name)
	document.Table = strings.ToLower(strings.ReplaceAll(document.Table, " ", "_"))
	err := tx.Create(&document.Document).Error
	if err != nil {
		tx.Rollback()
		return domain.DocumentWithDetails{}, err
	}

	for _, detail := range document.Details {
		detail.DocumentID = document.Document.ID
		detail.Field = strings.ToLower(strings.ReplaceAll(detail.Field, " ", "_"))
		err = tx.Create(&detail).Error
		if err != nil {
			tx.Rollback()
			return domain.DocumentWithDetails{}, err
		}
	}

	err = tx.Exec("EXEC sp_CreateTableToDocument @id = ?", document.Document.ID).Error
	if err != nil {
		tx.Rollback()
		return domain.DocumentWithDetails{}, err
	}

	tx.Commit()

	return document, nil
}

func (database *DSNSource) UpdateDocument(document domain.EditableDocument) (domain.DocumentWithDetails, error) {
	tx := database.DB.Begin()

	document.Name = Capitalize(document.Name)
	err := tx.Model(&domain.Document{}).Where("id = ?", document.ID).Updates(domain.Document{Name: document.Name}).Error
	if err != nil {
		tx.Rollback()
		return domain.DocumentWithDetails{}, err
	}

	for _, detail := range document.Details {
		detail.DocumentID = document.ID
		detail.Field = strings.ToLower(strings.ReplaceAll(detail.Field, " ", "_"))
		if detail.ID == 0 {
			// new record
			err = tx.Create(&detail).Error
			if err != nil {
				tx.Rollback()
				return domain.DocumentWithDetails{}, err
			}

			err = tx.Exec("EXEC sp_AddFieldToDocument @id = ?, @documentID = ?", detail.ID, detail.DocumentID).Error
			if err != nil {
				tx.Rollback()
				return domain.DocumentWithDetails{}, err
			}
		} else {
			fmt.Println("detail", detail)
			err = tx.Model(&domain.DetailDocument{}).Where("id = ?", detail.ID).Updates(domain.DetailDocument{DocumentKey: detail.DocumentKey}).Error
			if err != nil {
				tx.Rollback()
				return domain.DocumentWithDetails{}, err
			}
		}
	}

	tx.Commit()

	currentDocument := domain.Document{}

	err = database.DB.Model(&domain.Document{}).Where("id = ?", document.ID).First(&currentDocument).Error
	if err != nil {
		return domain.DocumentWithDetails{}, err
	}

	return domain.DocumentWithDetails{
		Document: currentDocument,
		Details:  document.Details,
	}, nil
}

func (database *DSNSource) DeleteFieldDocument(id uint) error {
	tx := database.DB.Begin()

	document := domain.DetailDocument{}

	err := tx.Model(&domain.DetailDocument{}).Where("id = ?", id).First(&document).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	fieldDocument := domain.DetailDocument{}

	err = tx.Exec("EXEC sp_DropFielToDocument @id = ?, @documentID = ?", id, document.DocumentID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&domain.DetailDocument{}).Where("id = ?", id).Delete(&fieldDocument).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (database *DSNSource) DeleteDocument(id uint) error {
	tx := database.DB.Begin()

	document := domain.Document{}

	err := tx.Model(&domain.Document{}).Where("id = ?", id).Delete(&document).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Exec("EXEC sp_DropTableByDocument @id = ?", id).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (database *DSNSource) GetDocumentRowRecord(id uint, documentID uint) ([]domain.DocumentRowRecord, error) {

	document := domain.Document{}
	err := database.DB.Model(&domain.Document{}).Where("id = ?", documentID).First(&document).Error
	if err != nil {
		return []domain.DocumentRowRecord{}, err
	}

	documentDetails := []domain.DetailDocument{}
	err = database.DB.Model(&domain.DetailDocument{}).Where("document_id = ?", documentID).Find(&documentDetails).Error
	if err != nil {
		return []domain.DocumentRowRecord{}, err
	}

	rows, err := database.DB.Raw("EXEC sp_GetDocumentRowRecordValues @id = ?, @table_name = ?", id, document.Table).Rows()
	if err != nil {
		return []domain.DocumentRowRecord{}, err
	}
	defer rows.Close()

	values, err := SerializedRowsProcedure(rows)
	if err != nil {
		return []domain.DocumentRowRecord{}, err
	}

	records := []domain.DocumentRowRecord{}
	for _, documentDetail := range documentDetails {
		for i := range values {
			if ok := values[i][documentDetail.Field]; ok != nil {
				records = append(records, domain.DocumentRowRecord{
					ID:          documentDetail.ID,
					Field:       documentDetail.Field,
					DocumentKey: documentDetail.DocumentKey,
					TypeField:   documentDetail.TypeField,
					Value:       values[i][documentDetail.Field],
				})
			}
		}
	}

	return records, nil
}

func (database *DSNSource) GetReports() ([]domain.Report, error) {
	reports := []domain.Report{}
	err := database.DB.Model(&domain.Report{}).Find(&reports).Error
	if err != nil {
		return []domain.Report{}, err
	}

	return reports, nil
}

func (database *DSNSource) GetReportByID(id uint) (domain.ReportByID, error) {
	var response domain.ReportByID
	report := domain.Report{}

	err := database.DB.Model(&domain.Report{}).Where("id = ?", id).First(&report).Error
	if err != nil {
		return domain.ReportByID{}, err
	}

	response.Report = domain.ReportData{
		ID:   report.ID,
		Name: report.Name,
	}

	rows, err := database.DB.Raw("EXEC " + report.StoredProcedure).Rows()
	if err != nil {
		return domain.ReportByID{}, err
	}
	defer rows.Close()

	table, columns, err := SerializedTableAndColumns(rows)
	if err != nil {
		return domain.ReportByID{}, err
	}

	response.Table = table
	response.Columns = columns

	chars, err := database.GetChartAndLinesByReport(id)
	if err != nil {
		return domain.ReportByID{}, err
	}

	response.Charts = chars

	return response, nil
}

func (database *DSNSource) GetChartAndLinesByReport(id uint) ([]domain.ChartReport, error) {
	charts := []domain.ChartReport{}

	err := database.DB.Model(&domain.ChartReport{}).Where("report_id = ?", id).Find(&charts).Error
	if err != nil {
		return []domain.ChartReport{}, err
	}

	for i := range charts {
		lines := []domain.ChartLine{}
		err = database.DB.Model(&domain.ChartLine{}).Where("chart_id = ?", charts[i].ID).Find(&lines).Error
		if err != nil {
			return []domain.ChartReport{}, err
		}

		charts[i].Lines = lines
	}

	return charts, nil
}

func (database *DSNSource) UpdateDocumentRowRecord(data domain.EditableDoumentRowRecord) error {
	tx := database.DB.Begin()

	recordsString, err := Stringify(data.Records)
	if err != nil {
		return err
	}

	err = tx.Exec("EXEC sp_UpdateDocumentRowRecord @id = ?, @document_id = ?, @fields = ?", data.ID, data.DocumentID, recordsString).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (database *DSNSource) DeleteDocumentRowRecord(id uint, documentID uint) error {
	return database.DB.Exec("EXEC sp_DeleteDocumentRowRecord @id = ?, @document_id = ?", id, documentID).Error
}

func (database *DSNSource) GetTables() ([]domain.Document, error) {
	table := []domain.Document{}
	err := database.DB.Model(&domain.Document{}).Find(&table).Error
	if err != nil {
		return []domain.Document{}, err
	}

	return table, nil
}

func (database *DSNSource) GetTableByID(id uint) (domain.TableByID, error) {
	var response domain.TableByID
	document := domain.Document{}

	err := database.DB.Model(&domain.Document{}).Where("id = ?", id).First(&document).Error
	if err != nil {
		return domain.TableByID{}, err
	}

	response.Document = domain.ReportData{
		ID:   document.ID,
		Name: document.Name,
	}

	// exec sp with 2 select results
	rows, err := database.DB.Raw("EXEC sp_GetDocumentTableByID @id = ?", id).Rows()
	if err != nil {
		return domain.TableByID{}, err
	}

	// go language data
	// defer on golang is executed when the function ends exacly before the return
	// dato del lenguaje go
	// defer en golang se ejecuta cuando la funcion termina exactamente antes del return
	defer rows.Close() // defer clause + function to exec()

	table, columns, err := SerializedTableAndColumns(rows)
	if err != nil {
		return domain.TableByID{}, err
	}

	response.Table = table
	response.Columns = columns

	return response, nil
}

func (database *DSNSource) UploadDocument(documentID uint, file string) error {
	return database.DB.Exec("EXEC sp_AppendDocumentData @document_id = ?, @file = ?", documentID, file).Error
}
