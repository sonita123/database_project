migrate \
-path db/migrations \
-database "sqlserver://@127.0.0.1:1433?database=unibazar&trusted_connection=yes" \
down




#for docker file
#migrate \
#-path db/migrations \
#-database "sqlserver://sa:my_view_898@127.0.0.1:1433?database=unibazar&encrypt=disable" \
#down
