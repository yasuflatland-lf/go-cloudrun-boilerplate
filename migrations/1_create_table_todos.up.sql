create table todos(
   id BIGINT NOT NULL AUTO_INCREMENT,
   slug VARCHAR (50) NOT NULL,
   task MEDIUMTEXT NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   constraint todos_pk
       primary key (id),
       key (slug),
       key (task(512))
) comment 'Todo' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
