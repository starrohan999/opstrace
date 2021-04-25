alter table "public"."exporter" drop constraint "exporter_pkey";
alter table "public"."exporter"
    add constraint "exporter_pkey" 
    primary key ( "id" );
