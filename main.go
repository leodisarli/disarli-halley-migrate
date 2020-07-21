package main

import (
    "strings"
    "time"
)

var dbHost string
var dbPass string
var dbPort string
var dbUser string
var charsPerPage int
var bulkLimit int
var verbose int

type Collection struct {
    Field1 int `json:"field1"`
    Field2 int `json:"field1"`
    Field3 string `json:"field1"`
}

func main() {
    var startedAt time.Time = time.Now()
    
    initialize()

    var clientId string = getClientId()
    var totalPages int = 0
    var offset int = 0
    var mustGo bool = true
    
    db := dbConnect(clientId)
    defer db.Close()

    for ok := true; ok; ok = mustGo {
        mustGo = false

        results := searchData(
            db,
            offset,
        )
        defer results.Close()

        offset += bulkLimit
        for results.Next() {
            mustGo = true
            var collections Collection
            err := results.Scan(
                &collections.Field1,
                &collections.Field2,
                &collections.Field3,
            )
            panicAtDisco(err)

            debug("")
            debug(
                getDivider(),
            )
            debug("Deal with record: " + intToString(collections.Field1))
            progress()

            if len(collections.Field3) > 0 {
                
                debug("Saving data for record: " + intToString(collections.Field1))
                result := wordWrap(
                    collections.Field3,
                    charsPerPage,
                )
                result = stringFormatter(result)
                var pages []string
                pages = strings.Split(
                    result,
                    "\n",
                )
                
                for page, content := range pages {
                    totalPages ++
                    page ++
                    debug("Saving page: " + intToString(page))

                    insertData(
                        db,
                        collections.Field1,
                        collections.Field2,
                        content,
                        page+1,
                    )
                }
                pages = nil
            }
        }
        results = nil
    }
    log(
        getDivider(),
    )
    log("Total saved pages: " + intToString(totalPages))

    endedAt := time.Now()
    takeTime := subtractTime(
        endedAt,
        startedAt,
    )

    log("Takes: " + floatToString(takeTime))
}