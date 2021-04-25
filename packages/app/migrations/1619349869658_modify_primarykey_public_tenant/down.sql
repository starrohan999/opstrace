alter table "public"."tenant" drop constraint "tenant_pkey";
alter table "public"."tenant"
    add constraint "tenant_pkey" 
    primary key ( "name" );
