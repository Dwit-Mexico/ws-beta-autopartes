package repository

import (
	"fmt"

	"gorm.io/gorm"
)

// stored procedures
const sp_CreateTableToDocument = `CREATE PROCEDURE [dbo].sp_CreateTableToDocument @id BIGINT
AS
BEGIN
    DECLARE @table NVARCHAR(300) = (SELECT [table] FROM documents WHERE id = @id)
    IF (EXISTS (SELECT *
                FROM INFORMATION_SCHEMA.tables
                WHERE TABLE_SCHEMA = 'dbo'
                  AND TABLE_NAME = @table))
        BEGIN
            THROW 50000, 'RECORD_ALREADY_EXIST', 1;
        END
    DECLARE @columns TABLE
                     (
                         name NVARCHAR(300),
                         type NVARCHAR(100)
                     )

    INSERT INTO @columns
    SELECT 'id',
           'BIGINT IDENTITY
                  PRIMARY KEY'

    INSERT INTO @columns
    SELECT field, type_field
    FROM detail_documents
    WHERE document_id = @id

    DECLARE @sql NVARCHAR(MAX) = 'CREATE TABLE ' + @table + ' (';
    SELECT @sql = @sql + name + ' ' + type + ',' FROM @columns;
    SET @sql = LEFT(@sql, LEN(@sql) - 1) + ');';

    EXEC sp_executesql @sql;
END
`
const sp_GetDocumentTableByID = `CREATE PROCEDURE [dbo].sp_GetDocumentTableByID @id INT
AS
BEGIN
    DECLARE @table_name NVARCHAR(100)
    DECLARE @document TABLE
                      (
                          id      BIGINT,
                          name    NVARCHAR(300),
                          [table] NVARCHAR(300)
                      )

    INSERT INTO @document (id, name, [table])
    SELECT id, name, [table]
    FROM documents d
    WHERE d.id = @id

    SET @table_name = (SELECT TOP 1 [table] FROM @document)

    DECLARE @data TABLE
                  (
                      uid         NVARCHAR(300),
                      table_name  NVARCHAR(300),
                      table_align NVARCHAR(20)
                  )
    INSERT INTO @data (uid, table_name, table_align) SELECT 'id', 'none', 'none'

    INSERT INTO @data (uid, table_name, table_align)
    SELECT dd.field, dd.document_key, IIF(dd.type_field = 'DECIMAL(10, 2)' OR dd.type_field = 'INT', 'end', 'start')
    FROM detail_documents dd
    WHERE dd.document_id = @id

    DECLARE @sql NVARCHAR(MAX) = 'SELECT '

    SELECT @sql =
           @sql + CONCAT(uid, '''', uid, '''', ',')

    FROM @data
    SET @sql = LEFT(@sql, LEN(@sql) - 1);
    SET @sql = CONCAT(@sql, ' FROM ', @table_name)

    PRINT @sql

    EXEC sp_executesql @sql

    SELECT uid, table_name AS name, table_align AS align FROM @data WHERE uid <> 'id'

END
`
const sp_AppendDocumentData = `CREATE PROCEDURE [dbo].sp_AppendDocumentData @document_id BIGINT, @file NVARCHAR(MAX)
AS
BEGIN
    DECLARE @table_name NVARCHAR(300);
    DECLARE @details TABLE
                     (
                         field        NVARCHAR(300),
                         document_key NVARCHAR(300),
                         type_field   NVARCHAR(300)
                     );
    DECLARE @sql NVARCHAR(MAX);

    SET @table_name = (SELECT [table] FROM documents WHERE id = @document_id);

    INSERT INTO @details (field, document_key, type_field)
    SELECT field, document_key, type_field
    FROM detail_documents
    WHERE document_id = @document_id;

    SET @sql = 'MERGE INTO ' + @table_name + ' AS trg USING (SELECT ';

    SELECT @sql = @sql + CONCAT(field, ',')
    FROM @details

    SET @sql = dbo.fn_DropEndChar(@sql) + ' FROM OPENJSON(@file) WITH (';

    SELECT @sql = @sql + CONCAT(d.field, ' ', d.type_field, ' ''', '$.', REPLACE(d.document_key, ' ', '_'), '''', ',')
    FROM @details d

    SET @sql = dbo.fn_DropEndChar(@sql) + ')) AS src ON ';

    SELECT @sql = @sql + CONCAT('trg.', d.field, ' = src.', d.field, ' AND ') FROM @details d

    SET @sql = dbo.fn_DropEndChars(@sql, 4) + ' WHEN NOT MATCHED THEN INSERT (';

    SELECT @sql = @sql + field + ',' FROM @details;

    SET @sql = dbo.fn_DropEndChar(@sql) + ') VALUES (';

    SELECT @sql = @sql + 'src.' + field + ', ' FROM @details;

    SET @sql = dbo.fn_DropEndChar(@sql) + ');';

    PRINT @sql

    EXEC sp_executesql @sql, N'@file NVARCHAR(MAX)', @file = @file;
END
`
const sp_DropTableByDocument = `CREATE PROCEDURE [dbo].sp_DropTableByDocument @id BIGINT
AS
BEGIN
    DECLARE @table_name NVARCHAR(100);
    SET @table_name = (SELECT TOP 1 [table] FROM documents WHERE id = @id)

    IF (EXISTS (SELECT *
                FROM INFORMATION_SCHEMA.tables
                WHERE table_schema = 'dbo'
                  AND table_name = @table_name))
        BEGIN
            DECLARE @sql NVARCHAR(200) = CONCAT('DROP TABLE dbo.', @table_name)

            EXEC sp_executesql @sql;
        END
END
`
const sp_AddFieldToDocument = `CREATE PROCEDURE [dbo].sp_AddFieldToDocument @id BIGINT, @documentID BIGINT
AS
BEGIN
    DECLARE
        @field NVARCHAR(100),
        @type_field NVARCHAR(50), @table_name NVARCHAR(100)


    SELECT @field = field, @type_field = type_field
    FROM detail_documents
    WHERE id = @id

    SET @table_name = (SELECT [table] FROM documents WHERE id = @documentID)
    DECLARE @sql NVARCHAR(MAX) = CONCAT('ALTER TABLE ', @table_name, ' ADD ', @field, ' ', @type_field)

    EXEC sp_executesql @sql
END
`
const sp_DeleteDocumentRowRecord = `CREATE PROCEDURE [dbo].sp_DeleteDocumentRowRecord @id BIGINT, @document_id BIGINT
AS
BEGIN
    DECLARE @table_name NVARCHAR(100) = (SELECT TOP 1 [table] FROM documents WHERE id = @document_id)

    DECLARE @sql NVARCHAR(500) = CONCAT('DELETE FROM ', @table_name, ' WHERE id = ', @id)

    EXEC sp_executesql @sql
END
`
const sp_DropFielToDocument = `CREATE PROCEDURE [dbo].sp_DropFielToDocument @id BIGINT, @documentID BIGINT
AS
BEGIN
    BEGIN
        DECLARE
            @field NVARCHAR(100),
            @table_name NVARCHAR(100)

        SET @field = (SELECT field FROM detail_documents WHERE id = @id);
        SET @table_name = (SELECT [table] FROM documents WHERE id = @documentID);

        DECLARE @sql NVARCHAR(MAX) = CONCAT('ALTER TABLE ', @table_name, ' DROP COLUMN ', @field, ';')

        EXEC sp_executesql @sql
    END
END
`
const sp_GetDocumentRowRecordValues = `CREATE PROCEDURE [dbo].sp_GetDocumentRowRecordValues @id BIGINT, @table_name NVARCHAR(100)
AS
BEGIN
    DECLARE @sql NVARCHAR(500) = CONCAT('SELECT * from ', @table_name, ' WHERE id = ', @id)

    EXEC sp_executesql @sql
END
`
const sp_UpdateDocumentRowRecord = `CREATE PROCEDURE [dbo].sp_UpdateDocumentRowRecord @id BIGINT, @document_id BIGINT, @fields NVARCHAR(MAX)
AS
BEGIN
    DECLARE @table_name NVARCHAR(100) = (SELECT [table] FROM documents WHERE id = @document_id)

    DECLARE @details TABLE
                     (
                         field     NVARCHAR(100),
                         typeField NVARCHAR(100)
                     )

    INSERT INTO @details (field, typeField)
    SELECT field, type_field
    FROM detail_documents
    WHERE document_id = @document_id

    DECLARE @sql NVARCHAR(MAX) = 'UPDATE t SET '
    SELECT @sql = @sql + CONCAT('t.', field, ' = fields.', field, ',') FROM @details

    SET @sql = dbo.fn_DropEndChar(@sql) + CONCAT(' FROM ', @table_name, ' t INNER JOIN (SELECT ', @id, ' AS id, ')

    SELECT @sql = @sql + CONCAT(field, ',') FROM @details

    SET @sql = dbo.fn_DropEndChar(@sql) + ' FROM OPENJSON(@fields) WITH ('

    SELECT @sql = @sql + CONCAT(field, ' ', typeField, ' ''', '$.', field, '''', ',') FROM @details

    SET @sql = dbo.fn_DropEndChar(@sql) + ')) AS fields ON fields.id = t.id'

    PRINT @sql
    EXEC sp_executesql @sql, N'@fields NVARCHAR(MAX)', @fields = @fields;
END
`
const sp_ExampleReport = `CREATE PROCEDURE [dbo].sp_example_report
AS
BEGIN
    DECLARE @report TABLE
                    (
                        field1 NVARCHAR(100),
                        field2 NVARCHAR(50),
                        date   DATETIME,
                        cost   DECIMAL(10, 2)
                    )

    DECLARE @columns table_columns;
    -- valid alignment 'start', 'center', 'end'

    INSERT INTO @report (field1, field2, date, cost) SELECT 'valor 1,1', 'valor 1,2', GETUTCDATE(), 10.40
    INSERT INTO @report (field1, field2, date, cost) SELECT 'valor 2,1', 'valor 2,2', GETUTCDATE(), 72.2
    INSERT INTO @report (field1, field2, date, cost) SELECT 'valor 3,1', 'valor 3,2', GETUTCDATE(), 9182.2
    INSERT INTO @report (field1, field2, date, cost) SELECT 'valor 4,1', 'valor 4,2', GETUTCDATE(), 92.4

    INSERT INTO @columns (uid, name) SELECT 'field1', 'Columna 1'
    INSERT INTO @columns (uid, name) SELECT 'field2', 'Columna 2'
    INSERT INTO @columns (uid, name) SELECT 'date', 'Fecha'
    INSERT INTO @columns (uid, name, align) SELECT 'cost', 'Costo por unidad de medida', 'end'

    SELECT * FROM @report
    SELECT * FROM @columns

END
`

// functions
const FuncDropEndChar = `CREATE FUNCTION fn_DropEndChar(@txt NVARCHAR(MAX))
    RETURNS NVARCHAR(MAX)
AS BEGIN
    RETURN LEFT(@txt, LEN(@txt) - 1)
END`
const FuncDropEndChars = `CREATE FUNCTION fn_DropEndChars(@txt NVARCHAR(MAX), @steps INT)
    RETURNS NVARCHAR(MAX)
AS
BEGIN
    RETURN LEFT(@txt, LEN(@txt) - @steps)
END
`

func ExistSP(db *gorm.DB, nombreSP string) bool {
	var existe int
	err := db.Raw("SELECT COUNT(*) FROM sys.procedures WHERE name = ?", nombreSP).Scan(&existe).Error

	if err != nil {
		fmt.Println("error when verifying the existence of the stored procedure: %w", err)
		return false
	}

	return existe > 0
}

func ExistFunc(db *gorm.DB, nombreFunc string) bool {
	var existe int
	err := db.Raw("SELECT COUNT(*) FROM sys.objects WHERE name = ? AND type = 'FN'", nombreFunc).Scan(&existe).Error

	if err != nil {
		fmt.Println("error when verifying the existence of the function: %w", err)
		return false
	}

	return existe > 0
}

func ExistTable(db *gorm.DB, name string) bool {
	var rowCount int
	err := db.Raw("SELECT COUNT(*) FROM sys.tables WHERE name = ?", name).Scan(&rowCount).Error

	fmt.Println(rowCount)
	if err != nil {
		fmt.Println("error when verifying the existence of the table: %w", err)
		return false
	}

	return rowCount > 0
}

func MigrateProcedures(db *gorm.DB) {

	// FUNCITONS
	exist := ExistFunc(db, "fn_DropEndChar")
	if !exist {
		db.Exec(FuncDropEndChar)
	}

	exist = ExistFunc(db, "fn_DropEndChars")
	if !exist {
		db.Exec(FuncDropEndChars)
	}

	// db.Exec("DROP PROCEDURE IF EXISTS sp_CreateTableToDocument")

	// STORED PROCEDURES
	exist = ExistSP(db, "sp_CreateTableToDocument")
	if !exist {
		db.Exec(sp_CreateTableToDocument)
	}

	exist = ExistSP(db, "sp_GetDocumentTableByID")
	if !exist {
		db.Exec(sp_GetDocumentTableByID)
	}

	exist = ExistSP(db, "sp_AppendDocumentData")
	if !exist {
		db.Exec(sp_AppendDocumentData)
	}

	exist = ExistSP(db, "sp_DropTableByDocument")
	if !exist {
		db.Exec(sp_DropTableByDocument)
	}

	exist = ExistSP(db, "sp_AddFieldToDocument")
	if !exist {
		db.Exec(sp_AddFieldToDocument)
	}

	exist = ExistSP(db, "sp_DeleteDocumentRowRecord")
	if !exist {
		db.Exec(sp_DeleteDocumentRowRecord)
	}

	exist = ExistSP(db, "sp_DropFielToDocument")
	if !exist {
		db.Exec(sp_DropFielToDocument)
	}

	exist = ExistSP(db, "sp_GetDocumentRowRecordValues")
	if !exist {
		db.Exec(sp_GetDocumentRowRecordValues)
	}

	exist = ExistSP(db, "sp_UpdateDocumentRowRecord")
	if !exist {
		db.Exec(sp_UpdateDocumentRowRecord)
	}

	// exist type
	db.Exec("CREATE TYPE table_columns AS TABLE (uid NVARCHAR(100), name NVARCHAR(100), align NVARCHAR(50));")

	exist = ExistSP(db, "sp_example_report")
	if !exist {
		db.Exec(sp_ExampleReport)
	}
}
