package enpass

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	// sqlcipher is necessary for sqlite crypto support
	sqlcipher "github.com/gdanko/gorm-sqlcipher"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	rows        []Card
	tableName   = "item"
	validFields = []string{
		"category",
		"login",
		"title",
		"uuid",
	}
)

const (
	// filename of the sqlite vault file
	vaultFileName = "vault.enpassdb"
	// contains info about your vault
	vaultInfoFileName = "vault.json"
)

type Tabler interface {
	TableName() string
}

func (Card) TableName() string {
	return tableName
}

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

func (v *Vault) openEncryptedDatabase(path string, dbKey []byte, logLevel logrus.Level) (err error) {
	var levelMap = map[logrus.Level]logger.LogLevel{
		logrus.PanicLevel: logger.Silent,
		logrus.FatalLevel: logger.Silent,
		logrus.ErrorLevel: logger.Silent,
		logrus.WarnLevel:  logger.Silent,
		logrus.InfoLevel:  logger.Silent,
		logrus.DebugLevel: logger.Info,
		logrus.TraceLevel: logger.Info,
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,        // Slow SQL threshold
			LogLevel:                  levelMap[logLevel], // Log level
			IgnoreRecordNotFoundError: true,               // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,               // Don't include params in the SQL log
			Colorful:                  true,               // Disable color
		},
	)

	gormConfig := &gorm.Config{
		PrepareStmt:            true,
		QueryFields:            true,
		SkipDefaultTransaction: true,
		Logger:                 newLogger,
	}

	// The raw key for the sqlcipher database is given
	// by the first 64 characters of the hex-encoded key
	dbName := fmt.Sprintf(
		"%s?_pragma_key=x'%s'&_pragma_cipher_compatibility=3",
		path,
		hex.EncodeToString(dbKey)[:masterKeyLength],
	)

	v.db, err = gorm.Open(sqlcipher.Open(dbName), gormConfig)
	if err != nil {
		return errors.Wrap(err, "could not open database")
	}

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
func (v *Vault) Open(credentials *VaultCredentials, logLevel logrus.Level) error {
	v.logger.Debug("generating database key")
	if err := v.generateAndSetDBKey(credentials); err != nil {
		return errors.Wrap(err, "could not generate database key")
	}

	v.logger.Debug("opening encrypted database")
	if err := v.openEncryptedDatabase(v.databaseFilename, credentials.DBKey, logLevel); err != nil {
		return errors.Wrap(err, "could not open encrypted database")
	}

	type Result struct {
		Name string `db:"name"`
	}

	var (
		result  Result
		results []Result
	)

	v.db.Select("name").Table("sqlite_master").Where("type = ?", "table").Where("name = ?", "item").Find(&results)
	for _, result = range results {
		if result.Name != "item" {
			return errors.New("could not connect to database")
		}
	}
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
func (v *Vault) GetEntries(cardType string, recordCategory, recordTitle, recordLogin, recordUuid []string, caseSensitive bool, orderbyFlag []string) ([]Card, error) {
	if v.db == nil || v.vaultInfo.VaultName == "" {
		return nil, errors.New("vault is not initialized")
	}

	rows, err := v.executeEntryQuery(cardType, recordCategory, recordTitle, recordLogin, recordUuid, caseSensitive, orderbyFlag)
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve cards from database")
	}

	var cards []Card
	for _, card := range rows {
		err = card.Decrypt()
		if err != nil {
			panic(err)
		}
		cards = append(cards, Card{
			UUID:           card.UUID,
			CreatedAt:      card.CreatedAt,
			Type:           card.Type,
			UpdatedAt:      card.UpdatedAt,
			Title:          card.Title,
			Subtitle:       card.Subtitle,
			Note:           card.Note,
			Trashed:        card.Trashed,
			Deleted:        card.Deleted,
			Category:       card.Category,
			Label:          card.Label,
			LastUsed:       card.LastUsed,
			Sensitive:      card.Sensitive,
			Icon:           card.Icon,
			DecryptedValue: card.DecryptedValue,
		})
	}

	return cards, nil
}

func (v *Vault) GetEntry(cardType string, recordCategory, recordTitle, recordLogin, recordUuid []string, caseSensitive bool, orderbyFlag []string, unique bool) (*Card, error) {
	cards, err := v.GetEntries(cardType, recordCategory, recordTitle, recordLogin, recordUuid, caseSensitive, orderbyFlag)
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

func (v *Vault) processFilters(filterList []string, columnName string, caseSensitive bool) (tx *gorm.DB) {
	var keyword string
	tx = v.db.Session(&gorm.Session{NewDB: true})

	if len(filterList) > 0 {
		for _, item := range filterList {
			keyword = "LIKE"
			if caseSensitive {
				keyword = "GLOB"
				item = strings.Replace(item, "%", "*", -1)
			}
			tx = tx.Or(
				fmt.Sprintf("%s %s ?", columnName, keyword),
				item,
			)
		}
	}
	return tx
}

func (v *Vault) executeEntryQuery(cardType string, recordCategory, recordTitle, recordLogin, recordUuid []string, caseSensitive bool, orderbyFlag []string) ([]Card, error) {
	query := v.db.Select("item.uuid", "itemField.type", "item.created_at", "item.field_updated_at", "item.title", "item.subtitle", "item.note", "item.trashed", "item.deleted", "item.category", "itemfield.label", "itemfield.value AS raw_value", "item.key", "item.last_used", "itemfield.sensitive", "item.icon").Table("item").Joins("INNER JOIN itemfield ON uuid = item_uuid")

	query.Where("item.deleted = ?", 0)
	query.Where("type = ?", cardType)

	query.Where(v.processFilters(recordCategory, "category", caseSensitive))
	query.Where(v.processFilters(recordTitle, "title", caseSensitive))
	query.Where(v.processFilters(recordLogin, "subtitle", caseSensitive))
	query.Where(v.processFilters(recordUuid, "uuid", caseSensitive))

	if len(orderbyFlag) > 0 {
		badFields := funk.SubtractString(orderbyFlag, validFields)
		goodFields := funk.IntersectString(orderbyFlag, validFields)
		if len(badFields) > 0 {
			v.logger.Warningf("the following fields cannot be used by --orderby: %s\n", strings.Join(badFields, ", "))
			if len(goodFields) <= 0 {
				v.logger.Warningf("after removing invalid --orderby fields, there are no fields remaining")
			}
		}

		if len(goodFields) > 0 {
			query.Order(strings.Join(goodFields, ","))
		}
	}

	query.Find(&rows)

	return rows, nil
}
