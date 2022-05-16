CREATE TABLE IF NOT EXISTS user (
    username TEXT PRIMARY KEY,
    password TEXT,
    email TEXT,
    oauth TEXT,
    picture TEXT
);

CREATE TABLE IF NOT EXISTS task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author TEXT,
    name TEXT,
    description TEXT,
    start TEXT, 
    end TEXT,
    allDay INTEGER, 
    done INTEGER,
    FOREIGN KEY(author) REFERENCES user(username)
);

CREATE TABLE IF NOT EXISTS task_label (
    taskId INTEGER NOT NULL ,
    label TEXT,
    FOREIGN KEY(taskId) REFERENCES task(id)
);

CREATE TABLE IF NOT EXISTS note(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    author INTEGER NOT NULL,
    title TEXT,
    content TEXT,
    FOREIGN KEY(author) REFERENCES user(username)
);

CREATE TABLE IF NOT EXISTS task_note (
    task INTEGER NOT NULL,
    note INTEGER NOT NULL,
    FOREIGN KEY(task) REFERENCES task(id),
    FOREIGN KEY(note) REFERENCES note(id)
);
