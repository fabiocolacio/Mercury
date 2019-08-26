package main

import(
    "database/sql"
    "fmt"
)

func initDB(db *sql.DB) error {
    _, err := db.Exec(fmt.Sprintf(`
        create table if not exists users (
            uid        int          not null  auto_increment,
            salt       binary(%d)   not null,
            shash      binary(%d)   not null,
            challenge  binary(%d),
            primary key (uid)
        );
    `, SALT_SIZE, SHASH_SIZE, CHALLENGE_SIZE))
    if err != nil {
        return err
    }

    _, err = db.Exec(`
        create table if not exists rooms (
            rid  int  not null  auto_increment,
            primary key (rid)
        );
    `)
    if err != nil {
        return err
    }

    _, err = db.Exec(`
        create table if not exists participants (
            user int  not null,
            room int  not null,
            primary key (user, room)
        );
    `)
    if err !=  nil {
        return err
    }

    _, err = db.Exec(`
        create table if not exists messages (
            mid  int       not null  auto_increment,
            data tinyblob  not null,
            room int       not null,
            primary key (mid)
        );
    `)
    return err
}

