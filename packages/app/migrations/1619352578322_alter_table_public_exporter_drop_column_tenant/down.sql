ALTER TABLE "public"."exporter" ADD COLUMN "tenant" text;
ALTER TABLE "public"."exporter" ALTER COLUMN "tenant" DROP NOT NULL;
alter table "public"."exporter" add constraint "exporter_tenant_fkey"
             foreign key ("tenant")
             references "public"."tenant"
             ("name") on update cascade on delete cascade;
