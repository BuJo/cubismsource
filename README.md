# cubismsource

## Building

I assume you have `go` installed.

    cd src
    go get code.google.com/p/gosqlite/sqlite
    go install jbossinfo jboss2sqlite cubismsource

As I haven't made the JBoss URLs configurable from outside, you should
probably adapt those:

    jboss2sqlite/pulljboss.go
    cubismsource/source_jboss.go

The sqlite3 database should have a table with the following schema:

    CREATE TABLE jvmmetrics (site text, date TIMESTAMP, free integer, max integer,total integer, threads numeric, xml blob, primary key(site, date));

## Using

    jboss2sqlite -interval=10 -mode=daemon -site=TST &
    cubismsource &

Both programs have a builtin help using `-h` showing the options and
the default values.

    http://localhost:8080/cubism

Wait a little, watch the output, of the programs - you should soon have
nice horizon graphs.

