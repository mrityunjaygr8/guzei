{ pkgs, config, ... }:

{
  # https://devenv.sh/basics/
  env = {
      GREET = "devenv";
      DB_NAME = "postgres123";
      DB_HOST = "localhost";
      DB_PORT = 5432;
      DB_USER = "postgres";
      DB_PASS = "postgres";
      DB_SSLMODE = "disable";
      DB_TEST_NAME = "tests";
  };

#  env = vars;

  # https://devenv.sh/packages/
  packages = with pkgs; [ git lazygit postgresql_15 fish go-migrate sqlc go-task air ];

  # https://devenv.sh/scripts/
  scripts.hello.exec = "echo hello from $GREET";

  enterShell = ''
    rm -rf .env
    cat > .env << EOF
DB_NAME=$DB_NAME
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_USER=$DB_USER
DB_PASS=$DB_PASS
DB_SSLMODE=$DB_SSLMODE
DB_TEST_NAME=$DB_TEST_NAME
EOF
  '';

#  dotenv.enable = true;
  dotenv.disableHint = true;

  # https://devenv.sh/languages/
  languages.go.enable = true;
  languages.go.package = pkgs.go_1_21;

  services.postgres = {
    enable = true;
    package = pkgs.postgresql_15;
    initialDatabases = [{ name = config.env.DB_NAME; } { name = config.env.DB_TEST_NAME;}];
    listen_addresses = "localhost";
    port = 5432;
    initialScript = ''
        CREATE USER ${config.env.DB_USER} SUPERUSER;
        ALTER USER ${config.env.DB_USER} WITH PASSWORD '${config.env.DB_PASS}';
    '';
  };

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  # processes.ping.exec = "ping example.com";

  # See full reference at https://devenv.sh/reference/options/
}
