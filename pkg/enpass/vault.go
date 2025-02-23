package enpass

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// sqlcipher is necessary for sqlite crypto support
	"github.com/gdanko/enpass/globals"
	"github.com/gdanko/enpass/pkg/unlock"
	"github.com/gdanko/enpass/util"
	sqlcipher "github.com/gdanko/gorm-sqlcipher"
	"github.com/miquella/ask"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	rows      []RawCard
	tableName = "item"
)

const (
	pinDefaultKdfIterCount = 100000
	pinMinLength           = 8
	vaultFileName          = "vault.enpassdb"
	vaultInfoFileName      = "vault.json"
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
	flagKeyFilePath string
	Password        string
	DBKey           []byte
}

func prompt(logger *logrus.Logger, flagNonInteractive bool, msg string) string {
	if !flagNonInteractive {
		if response, err := ask.HiddenAsk("Enter " + msg + ": "); err != nil {
			logger.WithError(err).Fatal("could not prompt for " + msg)
		} else {
			return response
		}
	}
	return ""
}

func (credentials *VaultCredentials) IsComplete() bool {
	return credentials.Password != "" || credentials.DBKey != nil
}

// QuoteList - quote a list for use in a database IN clause
func QuoteList(items []string) string {
	quoted := []string{}
	for _, item := range items {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", item))
	}
	return strings.Join(quoted, ",")
}

// DetermineVaultPath : Try to programatically determine the vault path based on the default path value
func DetermineVaultPath(logger *logrus.Logger, flagVaultPath string) (vaultPath string) {
	vaultPathFromConfig := globals.GetConfig().VaultPath

	if flagVaultPath != "" {
		vaultPath = util.ExpandPath(flagVaultPath)
	} else if vaultPathFromConfig != "" {
		vaultPath = util.ExpandPath(vaultPathFromConfig)
	} else {
		homeDir := globals.GetHomeDirectory()
		vaultPath = filepath.Join(homeDir, "Documents/Enpass/Vaults/primary")
	}
	logger.Debugf("vault path found at %s", vaultPath)

	return vaultPath
}

func OpenVault(logger *logrus.Logger, flagEnablePin bool, flagNonInteractive bool, vaultPath string, flagKeyFilePath string, logLevel logrus.Level, flagNoColor bool) (vault *Vault, credentials *VaultCredentials, err error) {
	vault, err = NewVault(vaultPath, logLevel, flagNoColor)
	if err != nil {
		panic(err)
	}

	var store *unlock.SecureStore
	if !flagEnablePin {
		logger.Debug("PIN disabled")
	} else {
		logger.Debug("PIN enabled, using store")
		store = InitializeStore(logger, vaultPath, flagNonInteractive)
		logger.Debug("initialized store")
	}
	credentials = AssembleVaultCredentials(logger, vaultPath, flagKeyFilePath, flagNonInteractive, store)

	return vault, credentials, nil
}

func InitializeStore(logger *logrus.Logger, vaultPath string, flagNonInteractive bool) *unlock.SecureStore {
	vaultPath, _ = filepath.EvalSymlinks(vaultPath)
	store, err := unlock.NewSecureStore(filepath.Base(vaultPath), logger.Level)
	if err != nil {
		logger.WithError(err).Fatal("could not create store")
	}

	pin := os.Getenv("ENP_PIN")
	if pin == "" {
		pin = prompt(logger, flagNonInteractive, "PIN")
	}
	if len(pin) < pinMinLength {
		logger.Fatal("PIN too short")
	}

	pepper := os.Getenv("ENP_PIN_PEPPER")

	pinKdfIterCount, err := strconv.ParseInt(os.Getenv("ENP_PIN_ITER_COUNT"), 10, 32)
	if err != nil {
		pinKdfIterCount = pinDefaultKdfIterCount
	}

	if err := store.GeneratePassphrase(pin, pepper, int(pinKdfIterCount)); err != nil {
		logger.WithError(err).Fatal("could not initialize store")
	}

	return store
}

func AssembleVaultCredentials(logger *logrus.Logger, vaultPath string, flagKeyFilePath string, flagNonInteractive bool, store *unlock.SecureStore) *VaultCredentials {
	var (
		vaultPassword           string
		vaultPasswordFromEnv    = os.Getenv("MASTERPW")
		vaultPasswordFromConfig = globals.GetConfig().VaultPassword
	)

	if vaultPasswordFromConfig != "" {
		logger.Debug("found a vault password in the configuration file")
		vaultPassword = vaultPasswordFromConfig
	} else if vaultPasswordFromEnv != "" {
		logger.Debug("found a vault password in the environment")
		vaultPassword = vaultPasswordFromEnv
	} else {
		logger.Debug("no vault password found - will prompt")
	}

	credentials := &VaultCredentials{
		Password:        vaultPassword,
		flagKeyFilePath: flagKeyFilePath,
	}

	if !credentials.IsComplete() && store != nil {
		var err error
		if credentials.DBKey, err = store.Read(); err != nil {
			logger.WithError(err).Fatal("could not read credentials from store")
		}
		logger.Debug("read credentials from store")
	}

	if !credentials.IsComplete() {
		credentials.Password = prompt(logger, flagNonInteractive, "vault password")
	}

	return credentials
}

// NewVault : Create new instance of vault and load vault info
func NewVault(vaultPath string, logLevel logrus.Level, flagNoColor bool) (*Vault, error) {
	v := Vault{
		logger:       *util.ConfigureLogger(logLevel, flagNoColor),
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

func (v *Vault) openEncryptedDatabase(path string, dbKey []byte, logLevel logrus.Level, flagNoColor bool) (err error) {
	colorful := true
	if flagNoColor {
		colorful = false
	}
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
			Colorful:                  colorful,           // Disable color
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

	if credentials.flagKeyFilePath == "" && v.vaultInfo.HasKeyfile == 1 {
		return errors.New("you should specify a keyfile")
	} else if credentials.flagKeyFilePath != "" && v.vaultInfo.HasKeyfile == 0 {
		return errors.New("you are specifying an unnecessary keyfile")
	}

	v.logger.Debug("generating master password")
	masterPassword, err := v.generateMasterPassword([]byte(credentials.Password), credentials.flagKeyFilePath)
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
func (v *Vault) Open(credentials *VaultCredentials, logLevel logrus.Level, flagNoColor bool) error {
	v.logger.Debug("generating database key")
	if err := v.generateAndSetDBKey(credentials); err != nil {
		return errors.Wrap(err, "could not generate database key")
	}

	v.logger.Debug("opening encrypted database")
	if err := v.openEncryptedDatabase(v.databaseFilename, credentials.DBKey, logLevel, flagNoColor); err != nil {
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
	if len(results) <= 0 {
		return errors.New("could not connect to database, please check the database credentials")
	}

	for _, result = range results {
		if result.Name != "item" {
			return errors.New("could not connect to database, please check the database credentials")
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

// GetEntries : return the flagCardType entries in the Enpass database filtered by option flags.
func (v *Vault) GetEntries(flagCardType string, flagRecordCategory, flagRecordTitle, flagRecordLogin, flagRecordUuid, flagLabel []string, flagCaseSensitive bool, flagOrderBy []string, validOrderBy []string) ([]Card, error) {
	if v.db == nil || v.vaultInfo.VaultName == "" {
		return nil, errors.New("vault is not initialized")
	}

	rows, err := v.executeEntryQuery(flagCardType, flagRecordCategory, flagRecordTitle, flagRecordLogin, flagRecordUuid, flagLabel, flagCaseSensitive, flagOrderBy, validOrderBy)
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
			Created:        card.Created,
			Type:           card.Type,
			Updated:        card.Updated,
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

func (v *Vault) GetEntry(flagCardType string, flagRecordCategory, flagRecordTitle, flagRecordLogin, flagRecordUuid, flagLabel []string, flagCaseSensitive bool, flagOrderBy []string, validOrderBy []string, unique bool) (*Card, error) {
	cards, err := v.GetEntries(flagCardType, flagRecordCategory, flagRecordTitle, flagRecordLogin, flagRecordUuid, flagLabel, flagCaseSensitive, flagOrderBy, validOrderBy)
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

func (v *Vault) processFilters(filterList []string, columnName string, flagCaseSensitive bool) (tx *gorm.DB) {
	var keyword string
	tx = v.db.Session(&gorm.Session{NewDB: true})

	if len(filterList) > 0 {
		for _, item := range filterList {
			keyword = "LIKE"
			if flagCaseSensitive {
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

func (v *Vault) executeEntryQuery(flagCardType string, flagRecordCategory, flagRecordTitle, flagRecordLogin, flagRecordUuid, flagLabel []string, flagCaseSensitive bool, flagOrderBy []string, validOrderBy []string) (cards []Card, err error) {
	var (
		configDefaultLabels []string
		configOrderByFields []string
		labels              []string = []string{}
		orderByFields       []string = []string{}
	)

	if len(flagOrderBy) > 0 {
		orderByFields = flagOrderBy
	} else if len(flagOrderBy) <= 0 {
		configOrderByFields = globals.GetConfig().OrderBy
		if len(configOrderByFields) > 0 {
			orderByFields = configOrderByFields
		}
	}

	if len(flagLabel) > 0 {
		labels = flagLabel
	} else if len(flagLabel) <= 0 {
		configDefaultLabels = globals.GetConfig().DefaultLabels
		if len(configDefaultLabels) > 0 {
			labels = configDefaultLabels
		}
	}

	query := v.db.Select("item.uuid", "itemField.type", "item.created_at AS created", "item.updated_at AS updated", "item.title", "item.subtitle", "item.note", "item.trashed", "item.deleted", "item.category", "itemfield.label", "itemfield.value AS raw_value", "item.key", "item.last_used", "itemfield.sensitive", "item.icon").Table("item").Joins("INNER JOIN itemfield ON uuid = item_uuid")

	query.Where("item.deleted = ?", 0)
	query.Where("type = ?", flagCardType)

	query.Where(v.processFilters(flagRecordCategory, "category", flagCaseSensitive))
	query.Where(v.processFilters(flagRecordTitle, "title", flagCaseSensitive))
	query.Where(v.processFilters(flagRecordLogin, "subtitle", flagCaseSensitive))
	query.Where(v.processFilters(flagRecordUuid, "uuid", flagCaseSensitive))
	query.Where(v.processFilters(labels, "label", flagCaseSensitive))

	if len(orderByFields) > 0 {
		badFields := funk.SubtractString(orderByFields, validOrderBy)
		goodFields := funk.IntersectString(orderByFields, validOrderBy)
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

	for i := range rows {
		cards = append(cards, Card{
			UUID:           rows[i].UUID,
			Created:        util.ToHuman(rows[i].Created),
			Type:           rows[i].Type,
			Updated:        util.ToHuman(rows[i].Updated),
			Title:          rows[i].Title,
			Subtitle:       rows[i].Subtitle,
			Note:           rows[i].Note,
			Trashed:        rows[i].Trashed,
			Deleted:        rows[i].Deleted,
			Category:       rows[i].Category,
			Label:          rows[i].Label,
			LastUsed:       util.ToHuman(rows[i].LastUsed),
			Icon:           rows[i].Icon,
			DecryptedValue: rows[i].DecryptedValue,
			RawValue:       rows[i].RawValue,
			Key:            rows[i].Key,
		})
	}

	return cards, nil
}
