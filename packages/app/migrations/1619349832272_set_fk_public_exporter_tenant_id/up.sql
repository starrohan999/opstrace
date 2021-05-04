alter table "public"."exporter" add constraint "exporter_tenant_id_fkey"
             foreign key ("tenant_id")
             references "public"."tenant"
             ("id") on update cascade on delete cascade;
