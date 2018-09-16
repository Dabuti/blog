CREATE TABLE post (
    id int auto_increment primary key,
    title varchar(50),
    date DATETIME,
    body varchar(5000)
);

CREATE TABLE tag (
    id int auto_increment primary key,
    title varchar(50)
);

CREATE TABLE post_tag (
    id_post int,
    id_tag int,
    primary key (id_post, id_tag)
);