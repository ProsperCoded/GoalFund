AWS EC2 Crashed severally because of memory and cpu cap, had to upgrade to a larger instance (t3.micro -> t3.small), run docker build in sequentially instead of parallel to avoid memory and cpu cap issues.

Relized that Most GoLang ORMs don't support schema generation and migration support, had to use third party tools like **atlas + gomigrate + gorm** just to manage the database schema and migrations.
