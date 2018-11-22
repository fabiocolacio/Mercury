package server

import(
    "fmt"
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
    return
}

func (serv *Server) ResetDB() (err error) {
    _, err = serv.db.Exec(`drop table if exists users;`)
    if err != nil {
        return
    }

    err = serv.InitDB()
    
    return
}
