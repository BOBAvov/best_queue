-- Таблица для хранения информации о группах
CREATE TABLE IF NOT EXISTS groups (
    id serial PRIMARY KEY, -- Уникальный идентификатор группы
    code varchar(32) UNIQUE NOT NULL, -- Уникальный код группы
    comment varchar(255) -- Комментарий к группе
);

-- Таблица для хранения информации о пользователях
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY, -- Уникальный идентификатор пользователя
    username varchar(255) NOT NULL, -- Имя пользователя
    tg_nick varchar(255) NOT NULL UNIQUE, -- Уникальный Telegram-ник пользователя
    group_id integer NOT NULL, -- Идентификатор группы, к которой принадлежит пользователь
    password_hash varchar(255) NOT NULL, -- Хэш пароля пользователя
    is_admin boolean NOT NULL DEFAULT FALSE -- Флаг, указывающий, является ли пользователь администратором
);

-- Таблица для хранения информации об очередях
CREATE TABLE IF NOT EXISTS queues (
    id serial NOT NULL UNIQUE, -- Уникальный идентификатор очереди
    title varchar(255), -- Название очереди
    time_start time without time zone NOT NULL, -- Время начала очереди
    time_end time without time zone NOT NULL, -- Время окончания очереди
    PRIMARY KEY (id) -- Первичный ключ
);

-- Таблица для хранения участников очередей
CREATE TABLE IF NOT EXISTS queue_participants (
    id serial PRIMARY KEY, -- Уникальный идентификатор участника очереди
    queue_id integer NOT NULL, -- Идентификатор очереди
    user_id integer NOT NULL, -- Идентификатор пользователя
    position integer NOT NULL, -- Позиция участника в очереди
    joined_at timestamp without time zone NOT NULL DEFAULT NOW(), -- Время присоединения к очереди
    is_active boolean NOT NULL DEFAULT TRUE, -- Флаг активности участника
    UNIQUE(queue_id, user_id) -- Уникальность участника в рамках одной очереди
);

-- Внешний ключ для связи пользователей с группами
ALTER TABLE users
    ADD CONSTRAINT Users_in_Groups_fk FOREIGN KEY (group_id) REFERENCES groups(id);

-- Внешний ключ для связи участников очереди с очередями
ALTER TABLE queue_participants
    ADD CONSTRAINT Queue_participants_queue_fk FOREIGN KEY (queue_id) REFERENCES queues(id);

-- Внешний ключ для связи участников очереди с пользователями
ALTER TABLE queue_participants
    ADD CONSTRAINT Queue_participants_user_fk FOREIGN KEY (user_id) REFERENCES users(id);

-- Индекс для ускорения поиска по Telegram-нику пользователей
CREATE INDEX IF NOT EXISTS users_tg_nick_index ON users (tg_nick);

-- Индекс для ускорения поиска по названию очередей
CREATE INDEX IF NOT EXISTS queues_title_index ON queues (title);

-- Индекс для ускорения поиска участников очереди по идентификатору очереди и позиции
CREATE INDEX IF NOT EXISTS queue_participants_queue_position_index ON queue_participants (queue_id, position);

-- Индекс для ускорения поиска участников очереди по идентификатору пользователя
CREATE INDEX IF NOT EXISTS queue_participants_user_index ON queue_participants (user_id);

-- Индекс для ускорения поиска пользователей по идентификатору группы
CREATE INDEX IF NOT EXISTS users_group_id_index ON users (group_id);

-- Индекс для ускорения поиска участников очереди по активности
CREATE INDEX IF NOT EXISTS queue_participants_is_active_index ON queue_participants (is_active);

-- Индекс для ускорения поиска участников очереди по позиции
CREATE INDEX IF NOT EXISTS queue_participants_position_index ON queue_participants (position);