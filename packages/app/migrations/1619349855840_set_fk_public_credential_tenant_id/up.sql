alter table "public"."credential" add constraint "credential_tenant_id_fkey"
             foreign key ("tenant_id")
             references "public"."tenant"
             ("id") on update cascade on delete cascade;
