package server

import(
    "fmt"
    "database/sql"
    "errors"
)

var(
    ErrNoSuchUser error = errors.New("No such user")
    ErrMsgUnsent error = errors.New("Failed to send message")
)

func (serv *Server) InitDB() (err error) {
    query := fmt.Sprintf(
        `create table users(
            uid int primary key auto_increment,
            username varchar(%d),
            salt binary(%d),
            saltedhash binary(%d));`,
        UsernameMaxLength,
        SaltLength,
        KeyHashLength)
    _, err = serv.db.Exec(query)
    if err != nil {
        return
    }

    _, err = serv.db.Exec(`create table messages(
        sender int,
        recipient int,
        timesent timestamp,
        message blob);`)

    return
}

func (serv *Server) ResetDB() (err error) {
    _, err = serv.db.Exec(`drop table if exists users;`)
    if err != nil {
        return
    }

    _, err = serv.db.Exec(`drop table if exists messages;`)
    if err != nil {
        return
    }

    err = serv.InitDB()

    return
}

func (serv *Server) SendMsg(message []byte, receiver string, sender int) (err error) {
    row := serv.db.QueryRow(`select uid from users where username = ?`, receiver)

    var recipient int
    err = row.Scan(&recipient)
    if err == sql.ErrNoRows {
        err = ErrNoSuchUser
        return
    }

    _, err = serv.db.Exec(`insert into messages
            (sender, recipient, message, timesent)
            values (?, ?, ?, NOW())`,
            sender, recipient, message)
    if err != nil {
        err = ErrMsgUnsent
    }

    return
}
