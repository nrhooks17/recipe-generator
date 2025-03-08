/*create user newUser with password 'password';
  create database newDatabase;
  grant connect on database newDatabase to newUser
  */

/*grant user full permissions across database*/
grant usage on schema public to nrhooks17;

grant all privileges on all tables in schema public to nrhooks17;
grant all privileges on all sequences in schema public to nrhooks17;

alter default privileges in schema public grant all privileges on tables to nrhooks17;
alter default privileges in schema public grant all privileges on sequences to nrhooks17;

