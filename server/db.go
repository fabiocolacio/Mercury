package server

import(
    "fmt"
)

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

