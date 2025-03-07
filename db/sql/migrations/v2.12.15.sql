alter table `project__schedule`
    alter column last_commit_hash type varchar(64);

alter table `task`
    alter column commit_hash type varchar(64);
