package gsql

import (
	"database/sql"
	"sort"
	"strings"
)

var Conn *sql.DB

func SetConn(driverName, dataSourceName string) (*sql.DB, error) {
	var err error
	Conn, err = sql.Open(driverName, dataSourceName)
	return Conn, err
}

// this allows us to the the go ` and replaces the MySql ` with ""
// so `"ID" int` becomes `ID` int
var quotes = strings.NewReplacer(`"`, "`")

type Migration struct {
	Up   string
	Down string
}

var Migrations = map[string]Migration{
	"000000000000_create_migrations": Migration{
		Up: quotes.Replace(`
			CREATE TABLE IF NOT EXISTS "migrations"(
			  "id" int unsigned not null auto_increment,
			  PRIMARY KEY("id"),
			  "Name" varchar(255)
			);`),
		Down: "DROP TABLE `migrations`;",
	},
}

func Migrate() error {
	hasRun, err := getMigrationsRan()
	if err != nil {
		return err
	}
	var run []string
	for name, _ := range Migrations {
		if !hasRun[name] {
			run = append(run, name)
		}
	}
	sort.Strings(run)

	for _, name := range run {
		if err := runMigration(name); err != nil {
			return err
		}
	}
	return nil
}

func Describe() (string, error) {
	hasRan, err := getMigrationsRan()
	if err != nil {
		return "", err
	}
	var all []string
	for name, _ := range Migrations {
		all = append(all, name)
	}
	sort.Strings(all)
	for i, s := range all {
		if hasRan[s] {
			all[i] = "* " + s
		} else {
			all[i] = "  " + s
		}
	}
	return strings.Join(all, "\n"), nil
}

func AddMigration(name, up, down string) Migration {
	m := Migration{
		Up:   quotes.Replace(up),
		Down: quotes.Replace(down),
	}
	Migrations[name] = m
	return m
}
func getMigrationsRan() (map[string]bool, error) {
	hasRan := make(map[string]bool)

	rows, err := Conn.Query("SHOW TABLES LIKE 'migrations';")
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		// not even migration migration has run yet
		return hasRan, nil
	}

	rows, err = Conn.Query("SELECT `name` FROM `migrations`")
	if err != nil {
		return nil, err
	}
	var name string

	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		hasRan[name] = true
	}
	return hasRan, nil
}

func runMigration(migration string) error {
	_, err := Conn.Exec(Migrations[migration].Up)
	if err != nil {
		return err
	}
	_, err = Conn.Exec("INSERT INTO migrations (name) VALUES (?)", migration)
	return err
}

func Rollback() (string, error) {
	rows, err := Conn.Query("SELECT `id`,`Name` FROM `migrations` ORDER BY `Name` DESC LIMIT 1;")
	if err != nil || !rows.Next() {
		return "", err
	}
	var id int
	var name string
	rows.Scan(&id, &name)
	_, err = Conn.Exec(Migrations[name].Down)
	if err != nil {
		return "", err
	}
	_, err = Conn.Exec("DELETE FROM `migrations` WHERE `id`=?;", id)
	if err != nil {
		return "", err
	}
	return name, nil
}
