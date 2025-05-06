self:
{
  config,
  lib,
  pkgs,
  ...
}:
let
  cfg = config.programs.todoist-cli;
  format = pkgs.formats.toml { };
in
{
  options.programs.todoist-cli = {
    enable = lib.mkEnableOption "Enable todoist-cli, a command line interface for Todoist";

    package = lib.mkPackageOption self.packages.${pkgs.system} "todoist-cli" { };

    settings = lib.mkOption {
      type = format.type;
      default = { };
      example = {
        daemon = {
          address = "@todo.sock";
          api_token_file = "/path/to/api_token";
        };
      };
      description = ''
        Configuration written to
        {file}`$XDG_CONFIG_HOME/todoist/config.toml` on Linux or
        {file}`$HOME/Library/Application Support/todoist/config.toml` on Darwin.
      '';
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [ cfg.package ];

    xdg.configFile = lib.mkIf (cfg.settings != { }) {
      "todoist/config.toml".source = format.generate "config.toml" cfg.settings;
    };

    systemd.user.services.todoist-cli = {
      Unit.Description = "Todoist CLI Daemon";

      Install.WantedBy = [ "default.target" ];

      Service = {
        ExecStart = "${lib.getExe cfg.package} daemon";
        Restart = "on-failure";
      };
    };
  };
}
