alter table "public"."exporter" drop constraint "exporter_pkey";
alter table "public"."exporter"
    add constraint "exporter_pkey"
    primary key ( "tenant_id", "name" );

alter table "public"."exporter" drop constraint "exporter_tenant_fkey",
             add constraint "exporter_tenant_id_fkey"
             foreign key ("tenant_id")
             references "public"."tenant"
             ("id") on update cascade on delete cascade;

alter table "public"."exporter" drop constraint "exporter_tenant_credential_fkey";

alter table "public"."credential" drop constraint "credential_pkey";
alter table "public"."credential"
    add constraint "credential_pkey"
    primary key ( "name", "tenant_id" );

alter table "public"."exporter"
           add constraint "exporter_tenant_id_credential_fkey"
           foreign key ("tenant_id", "credential")
           references "public"."credential"
           ("tenant_id", "name") on update cascade on delete restrict;

alter table "public"."credential" drop constraint "credential_tenant_fkey",
             add constraint "credential_tenant_id_fkey"
             foreign key ("tenant_id")
             references "public"."tenant"
             ("id") on update cascade on delete cascade;

alter table "public"."tenant" drop constraint "tenant_pkey";
alter table "public"."tenant"
    add constraint "tenant_pkey"
    primary key ( "id" );

ALTER TABLE "public"."exporter" DROP COLUMN "tenant" CASCADE;
ALTER TABLE "public"."credential" DROP COLUMN "tenant" CASCADE;

ALTER TABLE "public"."tenant" ADD CONSTRAINT "tenant_name_key" UNIQUE ("name");

