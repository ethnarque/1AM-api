{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    devenv.url = "github:cachix/devenv";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [
        inputs.devenv.flakeModule
      ];
      systems = ["x86_64-linux" "i686-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin"];

      perSystem = {
        lib,
        pkgs,
        ...
      }: {
        devenv.shells.default = {config, ...}: {
          name = "1am-api";

          packages = with pkgs; [
            ffmpeg
            air
            goose
            gopls
            nil
          ];

          dotenv.disableHint = true;
          dotenv.enable = true;

          env.DB_URL = "postgres://${config.env.DB_USER}@${config.env.DB_HOST}:${config.env.DB_PORT}/${config.env.DB_DATABASE}";
          env.GOOSE_DRIVER = "${config.env.DB_DRIVER}";
          env.GOOSE_DBSTRING = "postgres://${config.env.DB_USER}@${config.env.DB_HOST}:${config.env.DB_PORT}/${config.env.DB_DATABASE}";

          languages.go.enable = true;

          services.postgres.enable = true;
          services.postgres.initialDatabases = [
            {name = "${config.env.DB_DATABASE}";}
            {name = "${config.env.DB_DATABASE}-test";}
          ];
          services.postgres.initialScript = ''
            CREATE USER postgres SUPERUSER;
            CREATE ROLE ${config.env.DB_USER} WITH LOGIN PASSWORD '${config.env.DB_PASSWORD}';
            CREATE EXTENSION IF NOT EXISTS "citext"
            CREATE EXTENSION IF NOT EXISTS "uuid-ossp"
          '';
          services.postgres.listen_addresses = "127.0.0.1";

          processes.air.exec = "air";

          scripts.test-hello.exec = ''
            echo Hello!
          '';

          scripts."migrate-status".exec = ''
            goose -dir schemas/migrations status
          '';

          scripts."migrate-reset".exec = ''
            goose -dir schemas/migrations reset &> /dev/null
            goose -dir schemas/migrations up &> /dev/null
            echo "Database reinitialized!"
          '';

          scripts."migrate-up".exec = ''
            goose -dir schemas/migrations up
            echo "Database reinitialized!"
          '';

          scripts.test-app.exec = ''
            export GOOSE_DBSTRING=postgres://${config.env.DB_USER}@${config.env.DB_HOST}:${config.env.DB_PORT}/${config.env.DB_DATABASE}
            export DB_URL=postgres://pmlogist@localhost:5432/1am-test

            goose -dir schemas/migrations reset
            goose -dir schemas/migrations up
            goose -dir schemas/seeds -no-versioning up

            go test ./...
          '';
        };
      };
    };
}
