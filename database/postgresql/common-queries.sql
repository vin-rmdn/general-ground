-- Query to list columns in all tables in the public schema
select * from information_schema.columns where table_schema = 'public';