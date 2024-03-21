package queries

var CreateDogTableQuery string = `create table if not exists dog(
    dog_id   INTEGER PRIMARY KEY,
    dog_name TEXT,
    breed TEXT,
    location TEXT,
    image_url TEXT,
    contact_number INTEGER,
    owner_id INTEGER,
    is_active INTEGER,
    created_at TEXT,
    lastmodified_at TEXT,
    FOREIGN KEY (owner_id) REFERENCES user(user_id)
)`

var CreateUserTableQuery string = `create table if not exists user(
    user_id INTEGER PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    mail_id TEXT,
    user_name TEXT,
    encrypted_password TEXT,
    is_active INTEGER,
    created_at TEXT,
    lastmodified_at TEXT

)`

var CreateBreedTableQuery string = `create table if not exists breed(
breed_id INTEGER PRIMARY KEY,
breed_name TEXT
)`

var CreateFavouriteTableQuery string = `CREATE TABLE IF NOT EXISTS Favorite (
    favorite_id INTEGER PRIMARY KEY,
    user_id INTEGER,
    dog_id INTEGER,
    FOREIGN KEY(User_ID) REFERENCES user(user_id),
    FOREIGN KEY(dog_id) REFERENCES dog(dog_id)
);`
