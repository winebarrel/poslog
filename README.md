# poslog

Parser to extract SQL from postgresql.log

## Usage

```
$ cat postgresql.log
2022-05-30 04:59:41 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: select now();
2022-05-30 04:59:46 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: begin;
2022-05-30 04:59:48 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: insert into hello values (1);
2022-05-30 04:59:50 UTC:10.0.3.147(57382):postgres@postgres:[12768]:LOG:  statement: commit;
...

$ poslog postgresql.log
{"Timestamp":"2022-05-30 04:59:41 UTC","Host":"10.0.3.147","Port":"57382","User":"postgres","Database":"postgres","Pid":"[12768]","MessageType":"LOG","Duration":"","Statement":" select now();"}
{"Timestamp":"2022-05-30 04:59:46 UTC","Host":"10.0.3.147","Port":"57382","User":"postgres","Database":"postgres","Pid":"[12768]","MessageType":"LOG","Duration":"","Statement":" begin;"}
{"Timestamp":"2022-05-30 04:59:48 UTC","Host":"10.0.3.147","Port":"57382","User":"postgres","Database":"postgres","Pid":"[12768]","MessageType":"LOG","Duration":"","Statement":" insert into hello values (1);"}
{"Timestamp":"2022-05-30 04:59:50 UTC","Host":"10.0.3.147","Port":"57382","User":"postgres","Database":"postgres","Pid":"[12768]","MessageType":"LOG","Duration":"","Statement":" commit;"}
...
```
