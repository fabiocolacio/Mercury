package server

import(
    "fmt"
    "database/sql"
    "errors"
    "bytes"
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
    if err != nil {
        return
    }

    _, err = serv.db.Exec(`create table devices(
        did int primary key auto_increment,
        owner int,
        public_key blob,
        foreign key (owner) references users(uid));`)
    if err != nil {
        return
    }

    return
}

func (serv *Server) ResetDB() (err error) {
    _, err = serv.db.Exec(`drop table if exists devices;`)
    if err != nil {
        return
    }

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

func (serv *Server) MsgFetch(yourName string, myUid int, since string) ([]byte, error) {
    var(
        data []byte
        myName string
        yourUid int
    )

    row := serv.db.QueryRow(`select uid from users where username = ?`, yourName)
    err := row.Scan(&yourUid)
    if err == sql.ErrNoRows {
        err = ErrNoSuchUser
        return data, err
    }

    row = serv.db.QueryRow(`select username from users where uid = ?`, myUid)
    err = row.Scan(&myName)
    if err == sql.ErrNoRows {
        err = ErrNoSuchUser
        return data, err
    }

    rows, err := serv.db.Query(
        `SELECT users.username, messages.timesent, messages.message
        FROM messages
        INNER JOIN users
        ON messages.sender = users.uid
        WHERE
        (sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)`,
        myUid, yourUid, yourUid, myUid)
    if err != nil {
        return data, err
    }

    buffer := new(bytes.Buffer)

    fmt.Fprint(buffer, "[")
    firstMsg := true
    for rows.Next() {
        if firstMsg {
            firstMsg = false
        } else {
            fmt.Fprint(buffer, ",")
        }

        var(
            username string
            timestamp string
            message []byte
        )

        if err = rows.Scan(&username, &timestamp, &message); err != nil {
            return data, err
        }

        fmt.Fprint(buffer, "{")
        fmt.Fprintf(buffer, `"Username": "%s",`, username)
        fmt.Fprintf(buffer, `"Timestamp": "%s",`, timestamp)
        fmt.Fprintf(buffer, `"Message": %s`, string(message))
        fmt.Fprint(buffer, "}")
    }
    fmt.Fprint(buffer, "]")
    data = buffer.Bytes()

    return data, err
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
