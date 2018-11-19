package server

import(
    "fmt"
    "crypto/rand"
)

func (serv *Server) RegisterUser(creds Credentials) error {
    salt := make([]byte, SaltLength)
    rand.Read(salt)

    saltedHash := HashAndSaltPassword([]byte(creds.Password), salt)

    _, err := serv.db.Exec(
        `insert into users (username, salt, saltedhash) values (?, ?, ?);`,
        creds.Username, salt, saltedHash)

    return err
}

func (serv *Server) ResetDB() (err error) {
    var query string

    _, err = serv.db.Exec(`drop table if exists users;`)
    if err != nil {
        return
    }

    query = fmt.Sprintf(
        `create table users(
            uid int primary key auto_increment,
            username varchar(%d),
            salt binary(%d),
            saltedhash binary(%d));`,
        UsernameMaxLength,
        SaltLength,
        KeyHashLength)
    _, err = serv.db.Exec(query)

    return
}

