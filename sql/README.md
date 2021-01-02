# GA4GH Search SQL Parser

[SqlBase.g4](SqlBase.g4) (SQL grammar) based
on [SqlBase.g4](https://github.com/trinodb/trino/blob/master/presto-parser/src/main/antlr4/io/prestosql/sql/parser/SqlBase.g4)
in [Trino DB](https://github.com/trinodb/trino) project

to compile grammar

```
antlr4 -Dlanguage=Go -package parser -o parser -Werror -visitor SqlBase.g4
```


```
java \
  -cp /home/mlinhard/dev/projects/ga4gh-search-go/workspace/trino/presto-parser/target/presto-parser-351-SNAPSHOT.jar \
  -jar /home/mlinhard/dev/projects/ga4gh-search-go/workspace/java2go/build/libs/java2go-1.0-SNAPSHOT-all.jar \
  --src-package io.prestosql.sql.tree \
  --dest-package tree \
  -o tree.go
```
