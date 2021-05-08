package connconv

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"net/url"
	"strings"
)

// ParesClickhouseConnectionURL returns the mysql connection url for
// github.com/ClickHouse/clickhouse-go. It accepts a url in
// `clickhouse://user:pass@host:port/db?query` format.
// Any other value returns an error.
func ParesClickhouseConnectionURL(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	if u.User != nil {
		q.Add("username", u.User.Username())
		if pass, ok := u.User.Password(); ok {
			q.Add("password", pass)
		}
		u.User = nil
	}
	q.Add("database", strings.TrimPrefix(u.Path, "/"))
	u.Path = ""
	u.RawQuery = q.Encode()
	u.Scheme = "tcp"
	return u.String(), nil
}

// ParesMySQLConnectionURL returns the mysql connection url for
// github.com/go-sql-driver/mysql. It accepts a url in
// `mysql://user:password@host:port/database` format.
// Any other value returns an error.
func ParesMySQLConnectionURL(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	if u.Scheme != "mysql" {
		return "", fmt.Errorf("database url scheme must be mysql")
	}
	p, _ := u.User.Password()
	params := map[string]string{}
	q := u.Query()
	for k := range q {
		params[k] = q.Get(k)
	}
	params["parseTime"] = "true"
	c := mysql.NewConfig()
	c.User = u.User.Username()
	c.Passwd = p
	c.Net = "tcp"
	c.Addr = u.Host
	c.DBName = strings.TrimLeft(u.Path, "/")
	c.Params = params
	return c.FormatDSN(), nil
}
