ALTER TABLE "public"."exporter" DROP constraint "exporter_pkey";
ALTER TABLE "public"."exporter" ADD constraint "exporter_pkey"
    primary key ( "id" );

alter table "public"."exporter"
           add constraint "exporter_tenant_id_fkey"
           foreign key ("tenant_id")
           references "public"."tenant"
           ("id") on update no action on delete restrict;
ALTER TABLE "public"."exporter" DROP constraint "exporter_tenant_fkey"
-- ALTER TABLE "public"."exporter" DROP COLUMN "tenant";

alter table "public"."exporter"
           add constraint "exporter_credential_id_fkey"
           foreign key ("credential_id")
           references "public"."credential"
           ("id") on update no action on delete restrict;
ALTER TABLE "public"."exporter" DROP constraint "exporter_tenant_credential_fkey"
-- ALTER TABLE "public"."exporter" DROP COLUMN "credential";


ALTER TABLE "public"."credential" DROP constraint "credential_pkey";
ALTER TABLE "public"."credential" ADD constraint "credential_pkey"
    primary key ( "id" );

alter table "public"."credential"
           add constraint "credential_tenant_id_fkey"
           foreign key ("tenant_id")
           references "public"."tenant"
           ("id") on update no action on delete cascade;
ALTER TABLE "public"."credential" DROP constraint "credential_tenant_fkey"
-- ALTER TABLE "public"."credential" DROP COLUMN "tenant";

ALTER TABLE "public"."tenant" DROP constraint "tenant_pkey";
ALTER TABLE "public"."tenant" ADD constraint "tenant_pkey"
    primary key ( "id" );
