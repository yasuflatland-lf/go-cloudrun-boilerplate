create table todos(
   id BIGINT NOT NULL AUTO_INCREMENT,
   slug VARCHAR (50) NOT NULL,
   task MEDIUMTEXT NOT NULL,
   status BOOLEAN DEFAULT false,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   constraint todos_pk
       primary key (id),
       key (slug),
       key (task(512)),
       key (status)
) comment 'Todo' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
