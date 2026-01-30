package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func Test_SelectClient_WhenOk(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "demo.db")

	require.NoError(t, err)

	defer db.Close()
	clientID := 1

	var client Client

	err = db.QueryRow("SELECT id, fio, login, birthday, email FROM clients WHERE id = :id",
		sql.Named("id", clientID),
	).Scan(&client.ID, &client.FIO, &client.Login, &client.Birthday, &client.Email)

	require.NoError(t, err)

	assert.Equal(t, clientID, client.ID)
	assert.NotEmpty(t, client.FIO)
	assert.NotEmpty(t, client.Login)
	assert.NotEmpty(t, client.Email)
	// напиши тест здесь
}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "demo.db")

	require.NoError(t, err)

	defer db.Close()
	clientID := -1

	var client Client
	// напиши тест здесь

	err = db.QueryRow("SELECT id, fio, login, birthday, email FROM clients WHERE id = :id",
		sql.Named("id", clientID),
	).Scan(&client.ID, &client.FIO, &client.Login, &client.Birthday, &client.Email)

	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)

}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	// настройте подключение к БД
	db, err := sql.Open("sqlite", "demo.db")

	require.NoError(t, err)

	defer db.Close()
	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// напиши тест здесь
	res, err := db.Exec("INSERT INTO clients (fio, login,birthday,email) VALUES (:fio, :login,:birthday,:email)",
		sql.Named("fio", cl.FIO),
		sql.Named("login", cl.Login),
		sql.Named("birthday", cl.Birthday),
		sql.Named("email", cl.Email))
	require.NoError(t, err)

	b, err := res.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), b)

	n, err := res.LastInsertId()
	require.NoError(t, err)
	require.Greater(t, n, int64(0))

	var got Client

	err = db.QueryRow(
		`SELECT id, fio, login, birthday, email
		 FROM clients
		 WHERE id = :id`,
		sql.Named("id", n),
	).Scan(&got.ID, &got.FIO, &got.Login, &got.Birthday, &got.Email)
	require.NoError(t, err)

	assert.Equal(t, int(n), got.ID)
	assert.Equal(t, cl.FIO, got.FIO)
	assert.Equal(t, cl.Login, got.Login)
	assert.Equal(t, cl.Birthday, got.Birthday)
	assert.Equal(t, cl.Email, got.Email)

}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()
	// настройте подключение к БД

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	// напиши тест здесь
	res, err := db.Exec("INSERT INTO clients (fio, login,birthday,email) VALUES (:fio, :login,:birthday,:email)",
		sql.Named("fio", cl.FIO),
		sql.Named("login", cl.Login),
		sql.Named("birthday", cl.Birthday),
		sql.Named("email", cl.Email))
	require.NoError(t, err)

	b, err := res.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), b)

	n, err := res.LastInsertId()
	require.NoError(t, err)
	require.Greater(t, n, int64(0))

	_, err = db.Exec("DELETE FROM clients WHERE id = :id", sql.Named("id", int(n)))
	require.NoError(t, err)

	var got Client

	err = db.QueryRow(
		`SELECT id, fio, login, birthday, email
		 FROM clients
		 WHERE id = :id`,
		sql.Named("id", n),
	).Scan(&got.ID, &got.FIO, &got.Login, &got.Birthday, &got.Email)
	require.ErrorIs(t, err, sql.ErrNoRows)
}
