
-- ALTER TABLE "public"."tenant" DROP CONSTRAINT "tenant_name_key";

-- ALTER TABLE "public"."credential" ADD COLUMN "tenant" text;
-- ALTER TABLE "public"."credential" ALTER COLUMN "tenant" DROP NOT NULL;

-- ALTER TABLE "public"."exporter" ADD COLUMN "tenant" text;
-- ALTER TABLE "public"."exporter" ALTER COLUMN "tenant" DROP NOT NULL;

-- alter table "public"."tenant" drop constraint "tenant_pkey";
-- alter table "public"."tenant"
--     add constraint "tenant_pkey"
--     primary key ( "name" );

-- alter table "public"."credential" drop constraint "credential_tenant_id_fkey",
--           add constraint "credential_tenant_fkey"
--           foreign key ("tenant")
--           references "public"."tenant"
--           ("name")
--           on update cascade
--           on delete cascade;

-- alter table "public"."exporter" drop constraint "exporter_tenant_id_credential_fkey";

-- alter table "public"."credential" drop constraint "credential_pkey";
-- alter table "public"."credential"
--     add constraint "credential_pkey"
--     primary key ( "name", "tenant" );

-- alter table "public"."exporter" add foreign key ("credential", "tenant") references "public"."credential"("name", "tenant") on update cascade on delete restrict;

-- alter table "public"."exporter" drop constraint "exporter_tenant_id_fkey",
--           add constraint "exporter_tenant_fkey"
--           foreign key ("tenant")
--           references "public"."tenant"
--           ("name")
--           on update cascade
--           on delete cascade;

-- alter table "public"."exporter" drop constraint "exporter_pkey";
-- alter table "public"."exporter"
--     add constraint "exporter_pkey"
--     primary key ( "tenant", "name" );

