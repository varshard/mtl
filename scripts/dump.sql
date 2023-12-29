create table user
(
    password text        null,
    name     varchar(64) null,
    id       int auto_increment
        primary key
);

create table vote_item
(
    id          int auto_increment
        primary key,
    name        varchar(200) null,
    description text         null
);

create table user_vote
(
    user_id      int null,
    vote_item_id int null,
    constraint use_fk
        foreign key (user_id) references user (id),
    constraint vote_item_fk
        foreign key (vote_item_id) references vote_item (id)
);

