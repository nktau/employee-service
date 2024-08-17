create table if not exists employees (
    id serial primary key,
    name varchar(64),
    surname varchar(64),
    phone varchar(32),
    companyId int,
    passportType varchar(64),
    passportNumber varchar(32),
    departmentName varchar(64),
    departmentPhone varchar(32)
);