CREATE TRIGGER before_insert_todos
    BEFORE INSERT ON todos
    FOR EACH ROW
BEGIN
    SET new.slug = uuid();
END;