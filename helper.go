package main

import (
    "database/sql"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
    "github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql"
)

func wordWrap(
    content string,
    totalChars int,
) string {
	words := strings.Fields(
        strings.TrimSpace(content),
    )
	if len(words) == 0 {
		return content
	}
	result := words[0]
	charsLeft := totalChars - len(result)
	for _, word := range words[1:] {
		if len(word)+1 > charsLeft {
			result += "\n" + word
			charsLeft = totalChars - len(word)
		} else {
			result += " " + word
			charsLeft -= 1 + len(word)
		}
	}
	return result
}

func getNow() string {
    myTime := time.Now()
    now := fmt.Sprintf(
        "%d-%02d-%02d %02d:%02d:%02d",
        myTime.Year(),
        myTime.Month(),
        myTime.Day(),
        myTime.Hour(),
        myTime.Minute(),
        myTime.Second(),
    )
    return now
}

func subtractTime(
    startedTime time.Time,
    endedTime time.Time,
) float64 {
    diff := startedTime.Sub(
        endedTime,
    ).Seconds()
    return diff
}

func dbConnect(clientId string) *sql.DB {
    db, err := sql.Open(
        "mysql",
        dbUser +
            ":" +
            dbPass +
            "@tcp(" +
            dbHost +
            ":" +
            dbPort +
            ")/database_" +
            clientId,
    )
    panicAtDisco(err)
    return db
}

func searchData(
    db *sql.DB,
    offset int,
) *sql.Rows {
    sqlQuery := 
        `SELECT
            field1,
            field2,
            field3
        FROM table
        WHERE
            condition
        ORDER BY order
        DESC
        LIMIT ?
        OFFSET ?`
    results, err := db.Query(
        sqlQuery,
        bulkLimit,
        offset,
    )
    panicAtDisco(err)
    return results
}

func insertData(
    db *sql.DB,
    field1 int,
    field2 int,
    field3 string,
    page int,
) bool {
    sqlStatement := 
        `INSERT INTO
            table
            (field1, field2, field3, page)
        VALUES
            (?, ?, ?, ?)`
    _, err := db.Exec(
        sqlStatement,
        field1,
        field2,
        field3,
        page,
    )
    panicAtDisco(err)
    return true
}

func stringFormatter(content string) string {
    replacer := strings.NewReplacer(
        "<p>",
        "\n<p>",
        "<!--",
        "\n<!--",
        "-->",
        "-->\n",
        "<br",
        "\n<br",
    )
    result := replacer.Replace(content)
    return result
}

func panicAtDisco(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func initialize() {
    err := godotenv.Load()
    panicAtDisco(err)

    dbHost = os.Getenv("DB_HOST")
    dbUser = os.Getenv("DB_USER")
    dbPass = os.Getenv("DB_PASS")
    dbPort = os.Getenv("DB_PORT")
    charsPerPage = stringToInt(os.Getenv("CHARS_PER_PAGE"))
    bulkLimit = stringToInt(os.Getenv("BULK_LIMIT"))
    verbose = getVerbose()
}

func stringToInt(value string) int {
	result, err := strconv.Atoi(value)
    panicAtDisco(err)

	return result
}

func intToString(value int) string {
	result := strconv.Itoa(value)
	return result
}

func floatToString(value float64) string {
    result := fmt.Sprintf("%f", value)
	return result
}

func getClientId() string {
    if len(os.Args) > 1 {
        return os.Args[1]
    }
    fmt.Println("Missing mandatory parameter clientId")
    os.Exit(2)
    return ""
}

func getVerbose() int {
    if len(os.Args) > 2 {
        return stringToInt(os.Args[2])
    }
    return 0
}

func log(content string) {
    fmt.Println(content)
}

func debug(content string) {
    if verbose == 1 {
        fmt.Println(content)
    }
}

func getDivider() string {
    return "================================="
}

func progress() {
    if verbose == 0 {
        fmt.Print(".")
    }
}
