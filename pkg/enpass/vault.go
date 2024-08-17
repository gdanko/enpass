package enpass

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	sqliteEncrypt "github.com/BoredTape/gorm-sqlcipher"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const (
	// filename of the sqlite vault file
	vaultFileName = "vault.enpassdb"
	// contains info about your vault
	vaultInfoFileName = "vault.json"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

var (
	gormConfig = &gorm.Config{
		PrepareStmt:            true,
		QueryFields:            true,
		SkipDefaultTransaction: true,
		Logger:                 gormLogger.Default.LogMode(gormLogger.Info),
	}
	validFields = []string{
		"category",
		"login",
		"title",
	}
)

// Vault : vault is the container object for vault-related operations
type Vault struct {
	// Logger : the logger instance
	logger logrus.Logger

	// settings for filtering entries
	FilterFields []string
	FilterAnd    bool

	// vault.enpassdb : SQLCipher database
	databaseFilename string

	// vault.json
	vaultInfoFilename string

	// <uuid>.enpassattach : SQLCipher database files for attachments >1KB
	//attachments []string

	// pointer to our opened database
	db *gorm.DB

	// vault.json : contains info about your vault for synchronizing
	vaultInfo VaultInfo
}

type VaultCredentials struct {
	KeyfilePath string
	Password    string
	DBKey       []byte
}

func (credentials *VaultCredentials) IsComplete() bool {
	return credentials.Password != "" || credentials.DBKey != nil
}

// func processSqlIn(items []string, caseSensitive bool, columnName string) (whereItem string, outputValues []string) {
// 	var orSlice []string
// 	for _, value := range items {
// 		if caseSensitive {
// 			if strings.Contains(value, "%") {
// 				value = strings.Replace(value, "%", "*", -1)
// 				orSlice = append(orSlice, fmt.Sprintf("%s GLOB ?", columnName))
// 			} else {
// 				orSlice = append(orSlice, fmt.Sprintf("%s GLOB ?", columnName))
// 			}
// 		} else {
// 			orSlice = append(orSlice, fmt.Sprintf("%s LIKE ?", columnName))
// 		}
// 		outputValues = append(outputValues, value)
// 	}
// 	return fmt.Sprintf("(%s)", strings.Join(orSlice, " OR ")), outputValues
// }

// FileOrDirectoryExists : Determine if a file or directory exists
func FileOrDirectoryExists(path string) (exists bool, err error) {
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, err
	}
	return false, err
}

// QuoteList - quote a list for use in a database IN clause
func QuoteList(items []string) string {
	quoted := []string{}
	for _, item := range items {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", item))
	}
	return strings.Join(quoted, ",")
}

// FindDefaultVaultPath : Try to programatically determine the vault path based on the default path value
func FindDefaultVaultPath() (vaultPath string, err error) {
	userObj, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to determine the path of your home directory: %s", err)
	}
	vaultPath = filepath.Join(userObj.HomeDir, "Documents/Enpass/Vaults/primary")
	return vaultPath, nil
}

// ValidateVaultPath : Try to validate the specified vault path
func ValidateVaultPath(vaultPath string) (err error) {
	var exists bool

	vaultFile1 := filepath.Join(vaultPath, vaultFileName)
	vaultFile2 := filepath.Join(vaultPath, vaultInfoFileName)

	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		exists, err = FileOrDirectoryExists(vaultPath)
		if !exists && err != nil {
			return fmt.Errorf("the vault path \"%s\" does not exist - please use the --vault flag", vaultPath)
		}

		exists, err = FileOrDirectoryExists(vaultFile1)
		if !exists && err != nil {
			return fmt.Errorf("the vault file \"%s\" does not exist - please use the --vault flag", vaultFile1)
		}

		exists, err = FileOrDirectoryExists(vaultFile2)
		if !exists && err != nil {
			return fmt.Errorf("the vault file \"%s\" does not exist - please use the --vault flag", vaultFile2)
		}
	}
	return nil
}

// NewVault : Create new instance of vault and load vault info
func NewVault(vaultPath string, logLevel logrus.Level) (*Vault, error) {
	v := Vault{
		logger:       *logrus.New(),
		FilterFields: []string{"title", "subtitle"},
	}
	v.logger.SetLevel(logLevel)

	vaultPath, _ = filepath.EvalSymlinks(vaultPath)
	v.databaseFilename = filepath.Join(vaultPath, vaultFileName)
	v.vaultInfoFilename = filepath.Join(vaultPath, vaultInfoFileName)
	v.logger.Debug("checking provided vault paths")
	if err := v.checkPaths(); err != nil {
		return nil, err
	}

	v.logger.Debug("loading vault info")
	var err error
	v.vaultInfo, err = v.loadVaultInfo()
	if err != nil {
		return nil, errors.Wrap(err, "could not load vault info")
	}

	v.logger.
		WithField("db_path", vaultFileName).
		WithField("info_path", vaultInfoFileName).
		Debug("initialized paths")

	return &v, nil
}

func (v *Vault) openEncryptedDatabase(path string, dbKey []byte) (err error) {
	key := hex.EncodeToString(dbKey)[:masterKeyLength]
	var dbName string
	// dbName = fmt.Sprintf("%s?_pragma_key=x'%s'", path, key)
	// dbName = fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_compatibility=3", path, key)
	// dbName = fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_page_size=4096", path, key)
	dbName = fmt.Sprintf("%s?_pragma_key=x'%s'&_pragma_cipher_compatibility=3&_pragma_cipher_page_size=4096", path, key)

	db, err := gorm.Open(sqliteEncrypt.Open(dbName), &gorm.Config{})
	fmt.Println(111)
	fmt.Println(path)
	fmt.Println(key)
	fmt.Println(dbName)
	fmt.Println()
	fmt.Println(db)
	fmt.Println(err)
	fmt.Println(222)

	// os.Exit(0)

	if err != nil {
		fmt.Println("oh no!")
		panic(err)
	}
	fmt.Println(db)
	// The raw key for the sqlcipher database is given
	// by the first 64 characters of the hex-encoded key
	// hexKey := hex.EncodeToString(dbKey)[:masterKeyLength]
	// dbName := fmt.Sprintf(
	// 	"%s?_pragma_key=x'%s'&_pragma_cipher_compatibility=3",
	// 	path,
	// 	hexKey,
	// )
	// fmt.Println(111)
	// fmt.Println(path)
	// fmt.Println(hexKey)
	// fmt.Println(dbName)
	// fmt.Println(222)
	// v.db, err = gorm.Open(sqlcipher.Open(dbName), gormConfig)
	// if err != nil {
	// 	return errors.Wrap(err, "could not open database")
	// }

	// defer func() {
	// 	dbInstance, _ := v.db.DB()
	// 	_ = dbInstance.Close()
	// }()

	return nil
}

func (v *Vault) checkPaths() error {
	if _, err := os.Stat(v.databaseFilename); os.IsNotExist(err) {
		return errors.New("vault does not exist: " + v.databaseFilename)
	}

	if _, err := os.Stat(v.vaultInfoFilename); os.IsNotExist(err) {
		return errors.New("vault info file does not exist: " + v.vaultInfoFilename)
	}

	return nil
}

func (v *Vault) generateAndSetDBKey(credentials *VaultCredentials) error {
	if credentials.DBKey != nil {
		v.logger.Debug("skipping database key generation, already set")
		return nil
	}

	if credentials.Password == "" {
		return errors.New("empty vault password provided")
	}

	if credentials.KeyfilePath == "" && v.vaultInfo.HasKeyfile == 1 {
		return errors.New("you should specify a keyfile")
	} else if credentials.KeyfilePath != "" && v.vaultInfo.HasKeyfile == 0 {
		return errors.New("you are specifying an unnecessary keyfile")
	}

	v.logger.Debug("generating master password")
	masterPassword, err := v.generateMasterPassword([]byte(credentials.Password), credentials.KeyfilePath)
	if err != nil {
		return errors.Wrap(err, "could not generate vault unlock key")
	}

	v.logger.Debug("extracting salt from database")
	keySalt, err := v.extractSalt()
	if err != nil {
		return errors.Wrap(err, "could not get master password salt")
	}

	v.logger.Debug("deriving decryption key")
	credentials.DBKey, err = v.deriveKey(masterPassword, keySalt)
	if err != nil {
		return errors.Wrap(err, "could not derive database key from master password")
	}

	return nil
}

// Open : setup a connection to the Enpass database. Call this before doing anything.
func (v *Vault) Open(credentials *VaultCredentials) error {
	v.logger.Debug("generating database key")
	if err := v.generateAndSetDBKey(credentials); err != nil {
		return errors.Wrap(err, "could not generate database key")
	}

	v.logger.Debug("opening encrypted database")
	if err := v.openEncryptedDatabase(v.databaseFilename, credentials.DBKey); err != nil {
		return errors.Wrap(err, "could not open encrypted database")
	}
	fmt.Println("yay!!!")

	type Result struct {
		Name string `db:"name"`
	}

	var (
	// result  Result
	// results []Result
	)
	// v.db.Select("name").Table("sqlite_master")

	// v.db.Select("name").Table("sqlite_master").Where("type = ?", "table").Where("name = ?", "item").Find(&results)
	// for _, result = range results {
	// 	if result.Name == "item" {
	// 		return errors.New("could not connect to database")
	// 	}
	// }
	v.db.Select("*").Table("item")
	fmt.Println("yay!!!!!!")
	os.Exit(0)

	// var tableName string
	// err := v.db.QueryRow(`
	// 	SELECT name
	// 	FROM sqlite_master
	// 	WHERE type='table' AND name='item'
	// `).Scan(&tableName)
	// if err != nil {
	// 	return errors.Wrap(err, "could not connect to database")
	// } else if tableName != "item" {
	// 	return errors.New("could not connect to database")
	// }

	return nil
}

// Close : close the connection to the underlying database. Always call this in the end.
func (v *Vault) Close() {
	// if v.db != nil {
	// 	err := v.db.Close()
	// 	v.logger.WithError(err).Debug("closed vault")
	// }
}

// GetEntries : return the cardType entries in the Enpass database filtered by option flags.
func (v *Vault) GetEntries(cardType string, cardCategory, cardTitle, cardLogin []string, caseSensitive bool, orderbyFlag []string) ([]Card, error) {
	if v.db == nil || v.vaultInfo.VaultName == "" {
		return nil, errors.New("vault is not initialized")
	}

	rows, err := v.executeEntryQuery(cardType, cardCategory, cardTitle, cardLogin, caseSensitive, orderbyFlag)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve cards from database")
	}

	var cards []Card

	for rows.Next() {
		var card Card

		// read the database columns into Card object
		if err := rows.Scan(
			&card.UUID, &card.Type, &card.CreatedAt, &card.UpdatedAt, &card.Title,
			&card.Subtitle, &card.Note, &card.Trashed, &card.Deleted, &card.Category,
			&card.Label, &card.value, &card.itemKey, &card.LastUsed, &card.Sensitive, &card.Icon,
		); err != nil {
			return nil, errors.Wrap(err, "could not read card from database")
		}

		card.RawValue = card.value

		err = card.Decrypt()
		if err != nil {
			return nil, errors.Wrap(err, "could not decrypt card value")
		}

		cards = append(cards, card)
	}

	return cards, nil
}

func (v *Vault) GetEntry(cardType string, cardCategory, cardTitle, cardLogin []string, caseSensitive bool, orderbyFlag []string, unique bool) (*Card, error) {
	cards, err := v.GetEntries(cardType, cardCategory, cardTitle, cardLogin, caseSensitive, orderbyFlag)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve cards")
	}

	var ret *Card
	for _, card := range cards {
		if card.IsTrashed() || card.IsDeleted() {
			continue
		} else if ret == nil {
			ret = &card
		} else if unique {
			return nil, errors.New("multiple cards match that title")
		} else {
			break
		}
	}

	if ret == nil {
		return nil, errors.New("card not found")
	}

	return ret, nil
}

func (v *Vault) executeEntryQuery(cardType string, cardCategory, cardTitle, cardLogin []string, caseSensitive bool, orderbyFlag []string) (*sql.Rows, error) {

	// if len(cardCategory) > 0 {
	// 	var orSlice []string
	// 	for _, value := range cardCategory {
	// 		if caseSensitive {
	// 			if strings.Contains(value, "%") {
	// 				value = strings.Replace(value, "%", "*", -1)
	// 				orSlice = append(orSlice, fmt.Sprintf("%s GLOB ?", columnName))
	// 		}
	// }

	query := v.db.Select("uuid", "type", "created_at", "field_updated_at", "title", "subtitle", "note", "trashed", "item.deleted", "category", "label", "value", "key", "last_used", "sensitive", "item.icon").Table("item").Joins("INNER JOIN itemfield ON uuid = item_uuid")

	query = query.Where("deleted = ?", 0)
	query = query.Where("type = ?", "password")

	if len(cardCategory) > 0 {
		for _, categoryName := range cardCategory {
			if caseSensitive {
				categoryName = strings.Replace(categoryName, "%", "*", -1)
				query = query.Or("category GLOB ?", categoryName)
			} else {
				query = query.Or("category LIKE ?", categoryName)
			}
		}
	}

	if len(cardTitle) > 0 {
		for _, titleName := range cardTitle {
			if caseSensitive {
				titleName = strings.Replace(titleName, "%", "*", -1)
				query = query.Or("title GLOB ?", titleName)
			} else {
				query = query.Or("title LIKE ?", titleName)
			}
		}
	}

	if len(cardLogin) > 0 {
		for _, loginName := range cardLogin {
			if caseSensitive {
				loginName = strings.Replace(loginName, "%", "*", -1)
				query = query.Or("subtitle GLOB ?", loginName)
			} else {
				query = query.Or("subtitle LIKE ?", loginName)
			}
		}
	}

	if 1 == 2 {
		fmt.Println(cardType)
		fmt.Println(query)
	}

	return nil, nil
}
