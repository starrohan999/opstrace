alter table "public"."credential" drop constraint "credential_pkey";
alter table "public"."credential"
    add constraint "credential_pkey" 
    primary key ( "name", "tenant" );
