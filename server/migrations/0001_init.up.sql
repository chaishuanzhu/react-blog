CREATE TABLE users (
  id            BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  email         VARCHAR(128) NOT NULL,
  password_hash VARCHAR(100) NOT NULL,
  nickname      VARCHAR(32)  NOT NULL DEFAULT '',
  avatar        VARCHAR(255) NOT NULL DEFAULT '',
  website       VARCHAR(255) NOT NULL DEFAULT '',
  role          ENUM('admin') NOT NULL DEFAULT 'admin',
  created_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_users_email (email)
);

CREATE TABLE categories (
  id         BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name       VARCHAR(64) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_categories_name (name)
);

CREATE TABLE tags (
  id         BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name       VARCHAR(64) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_tags_name (name)
);

CREATE TABLE articles (
  id           BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  title        VARCHAR(200) NOT NULL,
  slug         VARCHAR(200) NOT NULL,
  content      MEDIUMTEXT   NOT NULL,
  category_id  BIGINT UNSIGNED NULL,
  status       ENUM('draft','published') NOT NULL DEFAULT 'draft',
  published_at DATETIME(3) NOT NULL,
  created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_articles_slug (slug),
  KEY idx_articles_status_published (status, published_at),
  KEY idx_articles_category (category_id),
  CONSTRAINT fk_articles_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL
);

CREATE TABLE article_tags (
  article_id BIGINT UNSIGNED NOT NULL,
  tag_id     BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (article_id, tag_id),
  KEY idx_article_tags_tag (tag_id),
  CONSTRAINT fk_article_tags_article FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE CASCADE,
  CONSTRAINT fk_article_tags_tag FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);

-- article_id NULL means the guestbook; parent_id NULL means a top-level comment.
CREATE TABLE comments (
  id         BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  article_id BIGINT UNSIGNED NULL,
  parent_id  BIGINT UNSIGNED NULL,
  nickname   VARCHAR(32)  NOT NULL,
  email      VARCHAR(128) NOT NULL,
  website    VARCHAR(255) NOT NULL DEFAULT '',
  avatar     VARCHAR(255) NOT NULL DEFAULT '',
  content    TEXT         NOT NULL,
  is_admin   TINYINT(1)   NOT NULL DEFAULT 0,
  ip         VARCHAR(45)  NOT NULL DEFAULT '',
  user_agent VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_comments_thread (article_id, parent_id, created_at),
  KEY idx_comments_parent (parent_id),
  KEY idx_comments_created (created_at),
  CONSTRAINT fk_comments_article FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE CASCADE,
  CONSTRAINT fk_comments_parent FOREIGN KEY (parent_id) REFERENCES comments (id) ON DELETE CASCADE
);

CREATE TABLE moments (
  id         BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  content    TEXT NOT NULL,
  images     JSON NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_moments_created (created_at)
);

CREATE TABLE friend_links (
  id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name        VARCHAR(64)  NOT NULL,
  url         VARCHAR(255) NOT NULL,
  avatar      VARCHAR(255) NOT NULL DEFAULT '',
  description VARCHAR(255) NOT NULL DEFAULT '',
  created_at  DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

CREATE TABLE changelogs (
  id         BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  items      JSON NOT NULL,
  logged_at  DATETIME(3) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_changelogs_logged (logged_at)
);

CREATE TABLE projects (
  id          BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name        VARCHAR(64)  NOT NULL,
  description VARCHAR(255) NOT NULL DEFAULT '',
  cover       VARCHAR(255) NOT NULL DEFAULT '',
  url         VARCHAR(255) NOT NULL DEFAULT '',
  sort_order  INT          NOT NULL DEFAULT 0,
  created_at  DATETIME(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_projects_sort (sort_order)
);

CREATE TABLE pages (
  page_key   ENUM('about_site','about_me') PRIMARY KEY,
  content    MEDIUMTEXT NOT NULL,
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
);

CREATE TABLE site_settings (
  id         TINYINT UNSIGNED PRIMARY KEY,
  notice     VARCHAR(255)    NOT NULL DEFAULT '',
  view_count BIGINT UNSIGNED NOT NULL DEFAULT 0,
  updated_at DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  CONSTRAINT chk_site_settings_singleton CHECK (id = 1)
);

INSERT INTO site_settings (id) VALUES (1);
INSERT INTO pages (page_key, content) VALUES ('about_site', ''), ('about_me', '');
