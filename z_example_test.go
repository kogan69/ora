// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

func ExampleDrvStmt_Exec_insert() {
	db, _ := sql.Open("ora", testConStr)
	defer db.Close()

	tableName := tableName()
	db.Exec(fmt.Sprintf("create table %v (c1 number)", tableName))

	// placeholder ':c1' is bound by position; ':c1' may be any name
	var value int64 = 9
	result, _ := db.Exec(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName), value)
	rowsAffected, _ := result.RowsAffected()
	fmt.Println(rowsAffected)
	// Output: 1
}

func ExampleDrvStmt_Exec_insert_return_identity() {
	db, _ := sql.Open("ora", testConStr)
	defer db.Close()

	tableName := tableName()
	db.Exec(fmt.Sprintf("create table %v (c1 number(19,0) generated always as identity (start with 1 increment by 1), c2 varchar2(48 char))", tableName))

	// use a 'returning into' SQL clause and specify a nil parameter to Exec
	// placeholder ':c1' is bound by position; ':c1' may be any name
	result, _ := db.Exec(fmt.Sprintf("insert into %v (c2) values ('go') returning c1 into :c1", tableName), nil)
	id, _ := result.LastInsertId()
	fmt.Println(id)
	// Output: 1
}

func ExampleDrvStmt_Exec_insert_bool() {
	db, _ := sql.Open("ora", testConStr)
	defer db.Close()

	tableName := tableName()
	db.Exec(fmt.Sprintf("create table %v (c1 char(1 byte))", tableName))

	// default false symbol is '0'
	// default true symbol is '1'
	// placeholder ':c1' is bound by position; ':c1' may be any name
	var value bool = true
	result, _ := db.Exec(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName), value)
	rowsAffected, _ := result.RowsAffected()
	fmt.Println(rowsAffected)
	// Output: 1
}

func ExampleDrvStmt_Exec_update() {
	db, _ := sql.Open("ora", testConStr)
	defer db.Close()

	tableName := tableName()
	db.Exec(fmt.Sprintf("create table %v (c1 number)", tableName))
	db.Exec(fmt.Sprintf("insert into %v (c1) values (9)", tableName))

	// placeholder ':three' and ':nine' are bound by position; ':three' and ':nine' may be any name
	var three int64 = 3
	var nine int64 = 9
	result, _ := db.Exec(fmt.Sprintf("update %v set c1 = :three where c1 = :nine", tableName), three, nine)
	rowsAffected, _ := result.RowsAffected()
	fmt.Println(rowsAffected)
	// Output: 1
}

func ExampleDrvStmt_Exec_delete() {
	db, _ := sql.Open("ora", testConStr)
	defer db.Close()

	tableName := tableName()
	db.Exec(fmt.Sprintf("create table %v (c1 number)", tableName))
	db.Exec(fmt.Sprintf("insert into %v (c1) values (9)", tableName))

	// placeholder ':1' is bound by position; ':1' may be any name
	var v int64 = 9
	result, _ := db.Exec(fmt.Sprintf("delete from %v where c1 = :1", tableName), v)
	rowsAffected, _ := result.RowsAffected()
	fmt.Println(rowsAffected)
	// Output: 1
}

func ExampleDrvStmt_Exec_Query() {
	db, _ := sql.Open("ora", testConStr)
	defer db.Close()

	tableName := tableName()
	db.Exec(fmt.Sprintf("create table %v (c1 number, c2 varchar2(48 char), c3 char(1 byte))", tableName))
	db.Exec(fmt.Sprintf("insert into %v (c1, c2, c3) values (3, 'slice', '0')", tableName))
	db.Exec(fmt.Sprintf("insert into %v (c1, c2, c3) values (7, 'map', '1')", tableName))
	db.Exec(fmt.Sprintf("insert into %v (c1, c2, c3) values (9, 'channel', '1')", tableName))

	// placeholder ':p' is bound by position; ':p' may be any name
	var value int64 = 8
	rows, _ := db.Query(fmt.Sprintf("select c1, c2, c3 from %v where c1 > :p", tableName), value)
	defer rows.Close()
	for rows.Next() {
		var c1 int64
		var c2 string
		var c3 bool
		rows.Scan(&c1, &c2, &c3)
		fmt.Printf("%v %v %v", c1, c2, c3)
	}
	// Output: 9 channel true
}

// TODO: Fix QueryRow
//func ExampleDrvStmt_Exec_QueryRow() {
//	db, _ := sql.Open("ora", testConStr)
//	defer db.Close()

//	tableName := tableName()
//	db.Exec(fmt.Sprintf("create table %v (c1 c1 number, c2 varchar2(48 char))", tableName))
//	db.Exec(fmt.Sprintf("insert into %v (c1) values (9, 'go')", tableName))

//	// placeholder ':p' is bound by position; ':p' may be any name
//	var c1 int64 = 9
//	var c2 string
//	db.QueryRow(fmt.Sprintf("select c2 from %v where c1 = :p", tableName), c1).Scan(&c2)
//	fmt.Println(c2)
//	// Output: go
//}

func ExampleStmt_Exec_insert() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert record
	var value int64 = 9
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	rowsAffected, _ := stmt.Exec(value)
	fmt.Println(rowsAffected)
	// Output: 1
}

func ExampleStmt_Exec_insert_return_identity() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(19,0) generated always as identity (start with 1 increment by 1), c2 varchar2(48 char))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert record
	var id int64
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c2) values ('go') returning c1 into :c1", tableName))
	defer stmt.Close()
	// pass a numeric pointer to rereive a database generated identity value
	stmt.Exec(&id)
	fmt.Println(id)
	// Output: 1
}

func ExampleStmt_Exec_insert_return_rowid() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert record
	var rowid string
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (9) returning rowid into :r", tableName))
	defer stmt.Close()
	// pass a string pointer to rereive a rowid
	stmt.Exec(&rowid)
	if rowid != "" {
		fmt.Println("Retrieved rowid")
	}
	// Output: Retrieved rowid
}

func ExampleStmt_Exec_insert_fetch_bool() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 char(1 byte))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert 'false' record
	var falseValue bool = false
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(falseValue)
	// insert 'true' record
	var trueValue bool = true
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(trueValue)

	// fetch inserted records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName))
	defer stmt.Close()
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v ", rset.Row[0])
	}
	// Output: false true
}

func ExampleStmt_Exec_insert_fetch_bool_alternate() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 char(1 byte))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// Update StmtConfig to change the FalseRune and TrueRune inserted into the database
	// insert 'false' record
	var falseValue bool = false
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Config.FalseRune = 'N'
	stmt.Exec(falseValue)
	// insert 'true' record
	var trueValue bool = true
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Config.TrueRune = 'Y'
	stmt.Exec(trueValue)

	// Update RsetConfig to change the TrueRune
	// used to translate an Oracle char to a Go bool
	// fetch inserted records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName))
	defer stmt.Close()
	stmt.Config.Rset.TrueRune = 'Y'
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v ", rset.Row[0])
	}
	// Output: false true
}

func ExampleStmt_Exec_update() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (9)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// update record
	var a int64 = 3
	var b int64 = 9
	stmt, _ = ses.Prep(fmt.Sprintf("update %v set c1 = :three where c1 = :nine", tableName))
	defer stmt.Close()
	rowsAffected, _ := stmt.Exec(a, b)
	fmt.Println(rowsAffected)
	// Output: 1
}

func ExampleStmt_Exec_delete() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (9)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// delete record
	var value int64 = 9
	stmt, _ = ses.Prep(fmt.Sprintf("delete from %v where c1 = :1", tableName))
	defer stmt.Close()
	rowsAffected, _ := stmt.Exec(value)
	fmt.Println(rowsAffected)
	// Output: 1
}

func ExampleStmt_Exec_insert_slice() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert one million rows with single round-trip to server
	values := make([]int64, 1000000)
	for n, _ := range values {
		values[n] = int64(n)
	}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	rowsAffected, _ := stmt.Exec(values)
	fmt.Println(rowsAffected)
	// Output: 1000000
}

func ExampleStmt_Exec_insert_nullable() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number, c2 varchar2(48 char), c3 char(1 byte))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// create nullable Go types for inserting null
	// insert record
	a := Int64{IsNull: true}
	b := String{IsNull: true}
	c := Bool{IsNull: true}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2, c3) values (:c1, :c2, :c3)", tableName))
	defer stmt.Close()
	rowsAffected, _ := stmt.Exec(a, b, c)
	fmt.Println(rowsAffected)
	// Output: 1
}

func ExampleStmt_Exec_insert_fetch_blob() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 blob)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// by default, byte slices are expected to be bound and retrieved
	// to/from a binary column such as a blob
	// insert record
	a := make([]byte, 10)
	for n, _ := range a {
		a[n] = byte(n)
	}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	rowsAffected, _ := stmt.Exec(a)
	fmt.Println(rowsAffected)

	// fetch record
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName))
	defer stmt.Close()
	rset, _ := stmt.Query()
	row := rset.NextRow()
	fmt.Println(row[0])

	// Output:
	// 1
	// [0 1 2 3 4 5 6 7 8 9]
}

func ExampleStmt_Exec_insert_fetch_byteSlice() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// note the NUMBER column
	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// Specify stmt.Config.SetByteSlice(U8)
	// Specify byte slice to be inserted into a NUMBER column
	// insert records
	a := make([]byte, 10)
	for n, _ := range a {
		a[n] = byte(n)
	}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Config.SetByteSlice(U8)
	rowsAffected, _ := stmt.Exec(a)
	fmt.Println(rowsAffected)

	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName))
	defer stmt.Close()
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v, ", rset.Row[0])
	}

	// Output:
	// 10
	// 0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
}

func ExampleStmt_Query() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number, c2 varchar2(48 char), c3 char(1 byte))", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2, c3) values (3, 'slice', '0')", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2, c3) values (7, 'map', '1')", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2, c3) values (9, 'channel', '1')", tableName))
	defer stmt.Close()
	stmt.Exec()

	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1, c2, c3 from %v", tableName))
	defer stmt.Close()
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v %v %v, ", rset.Row[0], rset.Row[1], rset.Row[2])
	}
	// Output: 3 slice false, 7 map true, 9 channel true,
}

func ExampleStmt_Query_nullable() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number, c2 varchar2(48 char), c3 char(1 byte))", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2, c3) values (null, 'slice', '0')", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2, c3) values (7, null, '1')", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2, c3) values (9, 'channel', null)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// Specify nullable return types to the Prep method
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1, c2, c3 from %v", tableName), OraI64, OraS, OraB)
	defer stmt.Close()
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v %v %v, ", rset.Row[0], rset.Row[1], rset.Row[2])
	}
	// Output: {true 0} {false slice} {false false}, {false 7} {true } {false true}, {false 9} {false channel} {true false},
}

func ExampleStmt_Query_numerics() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()
	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (9)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// Specify various numeric return types to the Prep method
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1, c1, c1, c1, c1, c1, c1, c1, c1, c1 from %v", tableName), I64, I32, I16, I8, U64, U32, U16, U8, F64, F32)
	defer stmt.Close()
	rset, _ := stmt.Query()
	row := rset.NextRow()
	fmt.Printf("%v %v %v %v %v %v %v %v %v %v",
		reflect.TypeOf(row[0]).Name(),
		reflect.TypeOf(row[1]).Name(),
		reflect.TypeOf(row[2]).Name(),
		reflect.TypeOf(row[3]).Name(),
		reflect.TypeOf(row[4]).Name(),
		reflect.TypeOf(row[5]).Name(),
		reflect.TypeOf(row[6]).Name(),
		reflect.TypeOf(row[7]).Name(),
		reflect.TypeOf(row[8]).Name(),
		reflect.TypeOf(row[9]).Name())
	// Output: int64 int32 int16 int8 uint64 uint32 uint16 uint8 float64 float32
}

func ExampleRset_Next() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert records
	a := make([]uint16, 5)
	for n, _ := range a {
		a[n] = uint16(n)
	}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	rowsAffected, _ := stmt.Exec(a)
	fmt.Println(rowsAffected)

	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), U16)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v, ", rset.Row[0])
	}
	// Output:
	// 5
	// 0, 1, 2, 3, 4,
}

func ExampleRset_NextRow() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number, c2 varchar2(48 char), c3 char(1 byte))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert record
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2, c3) values (7, 'go', '1')", tableName))
	defer stmt.Close()
	stmt.Exec()

	// fetch record
	stmt, _ = ses.Prep(fmt.Sprintf("select c1, c2, c3 from %v", tableName))
	rset, _ := stmt.Query()
	row := rset.NextRow()
	fmt.Printf("%v %v %v", row[0], row[1], row[2])
	// Output: 7 go true
}

func ExampleRset_cursor_single() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number, c2 varchar2(48 char))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert records
	a := make([]int64, 3)
	a[0] = 5
	a[1] = 7
	a[2] = 9
	b := make([]string, 3)
	b[0] = "Go is expressive, concise, clean, and efficient."
	b[1] = "Its concurrency mechanisms make it easy to"
	b[2] = "Go compiles quickly to machine code yet has"
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2) values (:1, :2)", tableName))
	stmt.Exec(a, b)

	// create proc
	stmt, _ = ses.Prep(fmt.Sprintf("create or replace procedure proc1(p1 out sys_refcursor) as begin open p1 for select c1, c2 from %v order by c1; end proc1;", tableName))
	defer stmt.Close()
	stmt.Exec()

	// pass *Rset to Exec for an out sys_refcursor
	// call proc
	stmt, _ = ses.Prep("call proc1(:1)")
	defer stmt.Close()
	rset := &Rset{}
	stmt.Exec(rset)
	if rset.IsOpen() {
		for rset.Next() {
			fmt.Println(rset.Row[0], rset.Row[1])
		}
	}
	// Output:
	// 5 Go is expressive, concise, clean, and efficient.
	// 7 Its concurrency mechanisms make it easy to
	// 9 Go compiles quickly to machine code yet has
}

func ExampleRset_cursor_multiple() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number, c2 varchar2(48 char))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert records
	a := make([]int64, 3)
	a[0] = 5
	a[1] = 7
	a[2] = 9
	b := make([]string, 3)
	b[0] = "Go is expressive, concise, clean, and efficient."
	b[1] = "Its concurrency mechanisms make it easy to"
	b[2] = "Go compiles quickly to machine code yet has"
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1, c2) values (:1, :2)", tableName))
	stmt.Exec(a, b)

	// create proc
	stmt, _ = ses.Prep(fmt.Sprintf("create or replace procedure proc1(p1 out sys_refcursor, p2 out sys_refcursor) as begin open p1 for select c1 from %v order by c1; open p2 for select c2 from %v order by c2; end proc1;", tableName, tableName))
	defer stmt.Close()
	stmt.Exec()

	// pass *Rset to Exec for an out sys_refcursor
	// call proc
	stmt, _ = ses.Prep("call proc1(:1, :2)")
	defer stmt.Close()
	rsetC1 := &Rset{}
	rsetC2 := &Rset{}
	stmt.Exec(rsetC1, rsetC2)
	fmt.Println("--- first result set ---")
	if rsetC1.IsOpen() {
		for rsetC1.Next() {
			fmt.Println(rsetC1.Row[0])
		}
	}
	fmt.Println("--- second result set ---")
	if rsetC2.IsOpen() {
		for rsetC2.Next() {
			fmt.Println(rsetC2.Row[0])
		}
	}
	// Output:
	// --- first result set ---
	// 5
	// 7
	// 9
	// --- second result set ---
	// Go compiles quickly to machine code yet has
	// Go is expressive, concise, clean, and efficient.
	// Its concurrency mechanisms make it easy to
}

func ExampleServer_Ping() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()

	// open a session before calling Ping
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	err := srv.Ping()
	if err == nil {
		fmt.Println("Ping sucessful")
	}
	// Output: Ping sucessful
}

func ExampleServer_Version() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()

	// open a session before calling Version
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	version, err := srv.Version()
	if version != "" && err == nil {
		fmt.Println("Received version from server")
	}
	// Output: Received version from server
}

func ExampleInt64() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(10,0))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Int64 slice
	a := make([]Int64, 5)
	a[0] = Int64{Value: -9}
	a[1] = Int64{Value: -1}
	a[2] = Int64{IsNull: true}
	a[3] = Int64{Value: 1}
	a[4] = Int64{Value: 9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraI64 to Prep method to return Int64 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraI64)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false -9}
	// {false -1}
	// {true 0}
	// {false 1}
	// {false 9}
}

func ExampleInt32() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(10,0))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Int32 slice
	a := make([]Int32, 5)
	a[0] = Int32{Value: -9}
	a[1] = Int32{Value: -1}
	a[2] = Int32{IsNull: true}
	a[3] = Int32{Value: 1}
	a[4] = Int32{Value: 9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraI32 to Prep method to return Int32 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraI32)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false -9}
	// {false -1}
	// {true 0}
	// {false 1}
	// {false 9}
}

func ExampleInt16() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(10,0))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Int16 slice
	a := make([]Int16, 5)
	a[0] = Int16{Value: -9}
	a[1] = Int16{Value: -1}
	a[2] = Int16{IsNull: true}
	a[3] = Int16{Value: 1}
	a[4] = Int16{Value: 9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraI16 to Prep method to return Int16 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraI16)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false -9}
	// {false -1}
	// {true 0}
	// {false 1}
	// {false 9}
}

func ExampleInt8() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(10,0))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Int8 slice
	a := make([]Int8, 5)
	a[0] = Int8{Value: -9}
	a[1] = Int8{Value: -1}
	a[2] = Int8{IsNull: true}
	a[3] = Int8{Value: 1}
	a[4] = Int8{Value: 9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraI8 to Prep method to return Int8 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraI8)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false -9}
	// {false -1}
	// {true 0}
	// {false 1}
	// {false 9}
}

func ExampleUint64() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(10,0))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Uint64 slice
	a := make([]Uint64, 5)
	a[0] = Uint64{Value: 0}
	a[1] = Uint64{Value: 3}
	a[2] = Uint64{IsNull: true}
	a[3] = Uint64{Value: 7}
	a[4] = Uint64{Value: 9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraU64 to Prep method to return Uint64 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraU64)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false 0}
	// {false 3}
	// {true 0}
	// {false 7}
	// {false 9}
}

func ExampleUint32() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(10,0))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Uint32 slice
	a := make([]Uint32, 5)
	a[0] = Uint32{Value: 0}
	a[1] = Uint32{Value: 3}
	a[2] = Uint32{IsNull: true}
	a[3] = Uint32{Value: 7}
	a[4] = Uint32{Value: 9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraU32 to Prep method to return Uint32 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraU32)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false 0}
	// {false 3}
	// {true 0}
	// {false 7}
	// {false 9}
}

func ExampleUint16() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(10,0))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Uint16 slice
	a := make([]Uint16, 5)
	a[0] = Uint16{Value: 0}
	a[1] = Uint16{Value: 3}
	a[2] = Uint16{IsNull: true}
	a[3] = Uint16{Value: 7}
	a[4] = Uint16{Value: 9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraU16 to Prep method to return Uint16 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraU16)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false 0}
	// {false 3}
	// {true 0}
	// {false 7}
	// {false 9}
}

func ExampleUint8() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(10,0))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Uint8 slice
	a := make([]Uint8, 5)
	a[0] = Uint8{Value: 0}
	a[1] = Uint8{Value: 3}
	a[2] = Uint8{IsNull: true}
	a[3] = Uint8{Value: 7}
	a[4] = Uint8{Value: 9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraU8 to Prep method to return Uint8 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraU8)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false 0}
	// {false 3}
	// {true 0}
	// {false 7}
	// {false 9}
}

func ExampleFloat64() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(16,15))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Float64 slice
	a := make([]Float64, 5)
	a[0] = Float64{Value: -float64(6.28318)}
	a[1] = Float64{Value: -float64(3.14159)}
	a[2] = Float64{IsNull: true}
	a[3] = Float64{Value: float64(3.14159)}
	a[4] = Float64{Value: float64(6.28318)}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraF64 to Prep method to return Float64 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraF64)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false -6.28318}
	// {false -3.14159}
	// {true 0}
	// {false 3.14159}
	// {false 6.28318}
}

func ExampleFloat32() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number(16,15))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Float32 slice
	a := make([]Float32, 5)
	a[0] = Float32{Value: -float32(6.28318)}
	a[1] = Float32{Value: -float32(3.14159)}
	a[2] = Float32{IsNull: true}
	a[3] = Float32{Value: float32(3.14159)}
	a[4] = Float32{Value: float32(6.28318)}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraF32 to Prep method to return Float32 values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraF32)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false -6.28318}
	// {false -3.14159}
	// {true 0}
	// {false 3.14159}
	// {false 6.28318}
}

func ExampleString() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 varchar2(48 char))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert String slice
	a := make([]String, 5)
	a[0] = String{Value: "Go is expressive, concise, clean, and efficient."}
	a[1] = String{Value: "Its concurrency mechanisms make it easy to"}
	a[2] = String{IsNull: true}
	a[3] = String{Value: "It's a fast, statically typed, compiled"}
	a[4] = String{Value: "One of Go's key design goals is code"}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraS to Prep method to return String values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraS)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false Go is expressive, concise, clean, and efficient.}
	// {false Its concurrency mechanisms make it easy to}
	// {true }
	// {false It's a fast, statically typed, compiled}
	// {false One of Go's key design goals is code}
}

func ExampleBool() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 char(1 byte))", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Bool slice
	a := make([]Bool, 5)
	a[0] = Bool{Value: true}
	a[1] = Bool{Value: false}
	a[2] = Bool{IsNull: true}
	a[3] = Bool{Value: false}
	a[4] = Bool{Value: true}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraB to Prep method to return Bool values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraB)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false true}
	// {false false}
	// {true false}
	// {false false}
	// {false true}
}

func ExampleTime() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 timestamp)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Time slice
	a := make([]Time, 5)
	a[0] = Time{Value: time.Date(2000, 1, 2, 3, 4, 5, 0, testDbsessiontimezone)}
	a[1] = Time{Value: time.Date(2001, 2, 3, 4, 5, 6, 0, testDbsessiontimezone)}
	a[2] = Time{IsNull: true}
	a[3] = Time{Value: time.Date(2003, 4, 5, 6, 7, 8, 0, testDbsessiontimezone)}
	a[4] = Time{Value: time.Date(2004, 5, 6, 7, 8, 9, 0, testDbsessiontimezone)}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraT to Prep method to return Time values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraT)
	rset, _ := stmt.Query()
	for rset.Next() {
		t := rset.Row[0].(Time)
		fmt.Printf("%v %v-%v-%v %v:%v:%v\n", t.IsNull, t.Value.Year(), t.Value.Month(), t.Value.Day(), t.Value.Hour(), t.Value.Minute(), t.Value.Second())
	}
	// Output:
	// false 2000-January-2 3:4:5
	// false 2001-February-3 4:5:6
	// true 1-January-1 0:0:0
	// false 2003-April-5 6:7:8
	// false 2004-May-6 7:8:9
}

func ExampleIntervalYM() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 interval year to month)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert IntervalYM slice
	a := make([]IntervalYM, 5)
	a[0] = IntervalYM{Year: 1, Month: 1}
	a[1] = IntervalYM{Year: 99, Month: 9}
	a[2] = IntervalYM{IsNull: true}
	a[3] = IntervalYM{Year: -1, Month: -1}
	a[4] = IntervalYM{Year: -99, Month: -9}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// fetch IntervalYM
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName))
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v, ", rset.Row[0])
	}
	// Output: {false 1 1}, {false 99 9}, {true 0 0}, {false -1 -1}, {false -99 -9},
}

func ExampleIntervalYM_ShiftTime() {
	interval := IntervalYM{Year: 1, Month: 1}
	actual := interval.ShiftTime(time.Date(2000, time.January, 0, 0, 0, 0, 0, time.Local))
	fmt.Println(actual.Year(), actual.Month(), actual.Day())
	// returns normalized date per time.AddDate
	// Output: 2001 January 31
}

func ExampleIntervalDS() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 interval day to second)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert IntervalDS slice
	a := make([]IntervalDS, 5)
	a[0] = IntervalDS{Day: 1, Hour: 1, Minute: 1, Second: 1, Nanosecond: 123456789}
	a[1] = IntervalDS{Day: 59, Hour: 59, Minute: 59, Second: 59, Nanosecond: 123456789}
	a[2] = IntervalDS{IsNull: true}
	a[3] = IntervalDS{Day: -1, Hour: -1, Minute: -1, Second: -1, Nanosecond: -123456789}
	a[4] = IntervalDS{Day: -59, Hour: -59, Minute: -59, Second: -59, Nanosecond: -123456789}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// fetch IntervalDS
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName))
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v, ", rset.Row[0])
	}
	// Output: {false 1 1 1 1 123457000}, {false 59 59 59 59 123457000}, {true 0 0 0 0 0}, {false -1 -1 -1 -1 -123457000}, {false -59 -59 -59 -59 -123457000},
}

func ExampleIntervalDS_ShiftTime() {
	interval := IntervalDS{Day: 1, Hour: 1, Minute: 1, Second: 1, Nanosecond: 123456789}
	actual := interval.ShiftTime(time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.Local))
	fmt.Println(actual.Day(), actual.Hour(), actual.Minute(), actual.Second(), actual.Nanosecond())
	// Output: 2 1 1 1 123456789
}

func ExampleBytes() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 blob)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Binary slice
	a := make([]Binary, 5)
	b := make([]byte, 10)
	for n, _ := range b {
		b[n] = byte(n)
	}
	a[0] = Binary{Value: b}
	b = make([]byte, 10)
	for n, _ := range b {
		b[n] = byte(n * 2)
	}
	a[1] = Binary{Value: b}
	a[2] = Binary{IsNull: true}
	b = make([]byte, 10)
	for n, _ := range b {
		b[n] = byte(n * 3)
	}
	a[3] = Binary{Value: b}
	b = make([]byte, 10)
	for n, _ := range b {
		b[n] = byte(n * 4)
	}
	a[4] = Binary{Value: b}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// Specify OraBin to Prep method to return ora.Binary values
	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName), OraBin)
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// {false [0 1 2 3 4 5 6 7 8 9]}
	// {false [0 2 4 6 8 10 12 14 16 18]}
	// {true []}
	// {false [0 3 6 9 12 15 18 21 24 27]}
	// {false [0 4 8 12 16 20 24 28 32 36]}
}

func ExampleBfile() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 bfile)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// insert Bfile
	a := Bfile{IsNull: false, DirectoryAlias: "TEMP_DIR", Filename: "test.txt"}
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (:c1)", tableName))
	defer stmt.Close()
	stmt.Exec(a)

	// fetch Bfile
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName))
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Printf("%v", rset.Row[0])
	}
	// Output: {false TEMP_DIR test.txt}
}

func ExampleTx() {
	// setup
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, _ := env.OpenSrv(testServerName)
	defer srv.Close()
	ses, _ := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()

	// create table
	tableName := tableName()
	stmt, _ := ses.Prep(fmt.Sprintf("create table %v (c1 number)", tableName))
	defer stmt.Close()
	stmt.Exec()

	// rollback
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (3)", tableName))
	tx, _ := ses.StartTx()
	stmt.Exec()
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (5)", tableName))
	stmt.Exec()
	tx.Rollback()

	// commit
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (7)", tableName))
	tx, _ = ses.StartTx()
	stmt.Exec()
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (9)", tableName))
	stmt.Exec()
	tx.Commit()

	// check that auto commit is reenabled
	stmt, _ = ses.Prep(fmt.Sprintf("insert into %v (c1) values (11)", tableName))
	stmt.Exec()

	// fetch records
	stmt, _ = ses.Prep(fmt.Sprintf("select c1 from %v", tableName))
	rset, _ := stmt.Query()
	for rset.Next() {
		fmt.Println(rset.Row[0])
	}
	// Output:
	// 7
	// 9
	// 11
}

func ExampleDriver_usage() {
	// example usage of the ora package driver
	// connect to a server and open a session
	env, _ := GetDrv().OpenEnv()
	defer env.Close()
	srv, err := env.OpenSrv(testServerName)
	defer srv.Close()
	if err != nil {
		panic(err)
	}
	ses, err := srv.OpenSes(testUsername, testPassword)
	defer ses.Close()
	if err != nil {
		panic(err)
	}

	// create table
	tableName := tableName()
	stmtTbl, err := ses.Prep(fmt.Sprintf("create table %v (c1 number(19,0) generated always as identity (start with 1 increment by 1), c2 varchar2(48 char))", tableName))
	defer stmtTbl.Close()
	if err != nil {
		panic(err)
	}
	rowsAffected, err := stmtTbl.Exec()
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected)

	// begin first transaction
	tx1, err := ses.StartTx()
	if err != nil {
		panic(err)
	}

	// insert record
	var id uint64
	str := "Go is expressive, concise, clean, and efficient."
	stmtIns, err := ses.Prep(fmt.Sprintf("insert into %v (c2) values (:c2) returning c1 into :c1", tableName))
	defer stmtIns.Close()
	rowsAffected, err = stmtIns.Exec(str, &id)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected)

	// insert nullable String slice
	a := make([]String, 4)
	a[0] = String{Value: "Its concurrency mechanisms make it easy to"}
	a[1] = String{IsNull: true}
	a[2] = String{Value: "It's a fast, statically typed, compiled"}
	a[3] = String{Value: "One of Go's key design goals is code"}
	stmtSliceIns, err := ses.Prep(fmt.Sprintf("insert into %v (c2) values (:c2)", tableName))
	defer stmtSliceIns.Close()
	if err != nil {
		panic(err)
	}
	rowsAffected, err = stmtSliceIns.Exec(a)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected)

	// fetch records
	stmtQuery, err := ses.Prep(fmt.Sprintf("select c1, c2 from %v", tableName))
	defer stmtQuery.Close()
	if err != nil {
		panic(err)
	}
	rset, err := stmtQuery.Query()
	if err != nil {
		panic(err)
	}
	for rset.Next() {
		fmt.Println(rset.Row[0], emptyString(rset.Row[1].(string)))
	}
	if rset.Err != nil {
		panic(rset.Err)
	}

	// commit first transaction
	err = tx1.Commit()
	if err != nil {
		panic(err)
	}

	// begin second transaction
	tx2, err := ses.StartTx()
	if err != nil {
		panic(err)
	}
	// insert null String
	nullableStr := String{IsNull: true}
	stmtTrans, err := ses.Prep(fmt.Sprintf("insert into %v (c2) values (:c2)", tableName))
	defer stmtTrans.Close()
	if err != nil {
		panic(err)
	}
	rowsAffected, err = stmtTrans.Exec(nullableStr)
	if err != nil {
		panic(err)
	}
	fmt.Println(rowsAffected)
	// rollback second transaction
	err = tx2.Rollback()
	if err != nil {
		panic(err)
	}

	// fetch and specify return type
	stmtCount, err := ses.Prep(fmt.Sprintf("select count(c1) from %v where c2 is null", tableName), U8)
	defer stmtCount.Close()
	if err != nil {
		panic(err)
	}
	rset, err = stmtCount.Query()
	if err != nil {
		panic(err)
	}
	row := rset.NextRow()
	if row != nil {
		fmt.Println(row[0])
	}
	if rset.Err != nil {
		panic(rset.Err)
	}

	// create stored procedure with sys_refcursor
	stmtProcCreate, err := ses.Prep(fmt.Sprintf("create or replace procedure proc1(p1 out sys_refcursor) as begin open p1 for select c1, c2 from %v where c1 > 2 order by c1; end proc1;", tableName))
	defer stmtProcCreate.Close()
	rowsAffected, err = stmtProcCreate.Exec()
	if err != nil {
		panic(err)
	}

	// call stored procedure
	// pass *Rset to Exec to receive the results of a sys_refcursor
	stmtProcCall, err := ses.Prep("call proc1(:1)")
	defer stmtProcCall.Close()
	if err != nil {
		panic(err)
	}
	procRset := &Rset{}
	rowsAffected, err = stmtProcCall.Exec(procRset)
	if err != nil {
		panic(err)
	}
	if procRset.IsOpen() {
		for procRset.Next() {
			fmt.Println(procRset.Row[0], emptyString(procRset.Row[1].(string)))
		}
		if procRset.Err != nil {
			panic(procRset.Err)
		}
		fmt.Println(procRset.Len())
	}

	// Output:
	// 0
	// 1
	// 4
	// 1 Go is expressive, concise, clean, and efficient.
	// 2 Its concurrency mechanisms make it easy to
	// 3 <empty>
	// 4 It's a fast, statically typed, compiled
	// 5 One of Go's key design goals is code
	// 1
	// 1
	// 3 <empty>
	// 4 It's a fast, statically typed, compiled
	// 5 One of Go's key design goals is code
	// 3
}

func emptyString(s string) string {
	if s == "" {
		return "<empty>"
	}
	return s
}
