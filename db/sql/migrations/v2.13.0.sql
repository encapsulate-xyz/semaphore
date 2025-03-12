alter table `runner` add column tag varchar(200);
alter table `project__template` add column runner_tag varchar(50);
alter table `project__template` add `runner_id` int null references runner(`id`) on delete set null;