package repository

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/RomanshkVolkov/ws-beta-autopartes/internal/core/domain"
	"golang.org/x/crypto/blake2b"
)

func ValidateError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func MaskString(s string) string {
	lenght := len(s)
	return strings.Repeat("*", lenght)
}

func TxtToHash(s string) string {
	hash := blake2b.Sum512([]byte(s))
	hexHash := hex.EncodeToString(hash[:])
	return hexHash
}

func TxtToRandomNumbers(s string) string {
	hexHash := TxtToHash(s)
	fmt.Println(hexHash)

	hashBytes, err := hex.DecodeString(hexHash)
	if err != nil {
		panic(err)
	}

	// Crear un lector de bytes a partir del hash
	reader := bytes.NewReader(hashBytes)

	// Usar el hash como semilla para un PRNG
	var seed int64
	err = binary.Read(reader, binary.LittleEndian, &seed)
	if err != nil {
		panic(err)
	}

	// validar que la semilla sea positiva
	if seed < 0 {
		seed = -seed
	}

	rng := rand.Reader
	_, err = rand.Int(rng, big.NewInt(seed))
	if err != nil {
		panic(err)
	}

	randomNumber, err := rand.Int(rng, big.NewInt(1e18))
	if err != nil {
		panic(err)
	}

	return randomNumber.String()
}

func CurrentTime() string {
	return time.Now().UTC().Format("2024-01-01 13:01:01")
}

func Capitalize(s string) string {
	firstLetter := strings.ToUpper(string(s[0]))
	restLetters := strings.ToLower(s[1:])
	return firstLetter + restLetters
}

func CapitalizeAll(s string) string {
	words := strings.Split(s, " ")
	var capitalizedWords []string
	for _, word := range words {
		capitalizedWords = append(capitalizedWords, Capitalize(word))
	}
	return strings.Join(capitalizedWords, " ")
}

func RemoveAccents(str string) string {
	charReplacer := map[rune]rune{
		'Á': 'A',
		'É': 'E',
		'Í': 'I',
		'Ó': 'O',
		'Ú': 'U',
		'á': 'a',
		'é': 'e',
		'í': 'i',
		'ó': 'o',
		'ú': 'u',
	}

	var result strings.Builder

	for _, char := range str {
		if newChar, ok := charReplacer[char]; ok {
			result.WriteRune(newChar)
		} else {
			result.WriteRune(char)
		}
	}

	return result.String()
}

func ReplaceSpacesWithUnderscores(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = RemoveAccents(s)
	s = ReplaceSpacesWithUnderscores(s)
	return s
}

func Stringify(v any) (string, error) {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func SerializedRowsProcedure(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		rowMap := make(map[string]any)
		for i, col := range columns {
			val := values[i]

			switch v := val.(type) {
			case []byte:
				rowMap[col] = string(v)
			default:
				rowMap[col] = v
			}
		}

		result = append(result, rowMap)
	}

	return result, nil
}

func SerializedTableAndColumns(rows *sql.Rows) ([]map[string]any, []domain.TableViewDefinition, error) {
	table, err := SerializedRowsProcedure(rows)
	if err != nil {
		return nil, nil, err
	}

	columnsResult := []domain.TableViewDefinition{}
	if rows.NextResultSet() {
		for rows.Next() {
			var uid, name string
			var align domain.Align

			err := rows.Scan(&uid, &name, &align)
			if err != nil {
				fmt.Println("error al escanear las filas de usuarios: %w", err)
				return nil, nil, err
			}

			columnsResult = append(columnsResult, domain.TableViewDefinition{
				UID:   uid,
				Name:  name,
				Align: domain.Align(align),
			})
		}
	}

	return table, columnsResult, nil
}
