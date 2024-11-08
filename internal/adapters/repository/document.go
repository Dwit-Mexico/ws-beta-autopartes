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

func (database *DSNSource) GetTables() ([]domain.Document, error) {
	table := []domain.Document{}
	err := database.DB.Model(&domain.Document{}).Find(&table).Error
	if err != nil {
		return []domain.Document{}, err
	}

	return table, nil
}

func (database *DSNSource) GetTableByID(id uint) (domain.TableByID, error) {
	// exec sp with 2 select results
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

	rows, err := database.DB.Raw("EXEC sp_GetDocumentTableByID @id = ?", id).Rows()
	if err != nil {
		return domain.TableByID{}, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return domain.TableByID{}, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			return domain.TableByID{}, err
		}

		// Create an object to hold the data
		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// tranform data type
			switch v := val.(type) {
			case []byte:
				entry[col] = string(v)
			default:
				entry[col] = v
			}
		}

		response.Table = append(response.Table, entry)
	}

	if rows.NextResultSet() {
		for rows.Next() {
			var uid, name string
			var align domain.Align
			err := rows.Scan(&uid, &name, &align)
			if err != nil {
				return domain.TableByID{}, fmt.Errorf("error al escanear las filas de usuarios: %w", err)
			}
			response.Columns = append(response.Columns, domain.TableViewDefinition{
				UID:   uid,
				Name:  name,
				Align: domain.Align(align),
			})
		}
	}

	if response.Table == nil {
		response.Table = []map[string]interface{}{}
	}
	return response, nil
}

func (database *DSNSource) UploadDocument(documentID uint, file string) error {
	return database.DB.Exec("EXEC sp_AppendDocumentData @document_id = ?, @file = ?", documentID, file).Error
}
