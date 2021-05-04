ALTER TABLE "public"."credential" ADD COLUMN "tenant" text;
ALTER TABLE "public"."credential" ALTER COLUMN "tenant" DROP NOT NULL;
alter table "public"."credential" add constraint "credential_tenant_fkey"
             foreign key ("tenant")
             references "public"."tenant"
             ("name") on update cascade on delete cascade;

