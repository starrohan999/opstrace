CREATE TRIGGER t_tenant_insert
  BEFORE INSERT ON tenant
  FOR EACH ROW
  WHEN (NEW.name IS NOT NULL AND NEW.url_slug IS NULL)
  EXECUTE PROCEDURE set_url_slug_from_name();