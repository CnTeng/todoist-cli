self: nixosTest:
let
  stateVersion = "25.05";
  home-manager = builtins.fetchTarball {
    url = "https://github.com/nix-community/home-manager/archive/release-${stateVersion}.tar.gz";
    sha256 = "12246mk1xf1bmak1n36yfnr4b0vpcwlp6q66dgvz8ip8p27pfcw2";
  };
in
nixosTest {
  name = "todoist-cli";
  nodes.machine = {
    imports = [ (import "${home-manager}/nixos") ];

    users.users.alice = {
      isNormalUser = true;
    };

    home-manager.sharedModules = [ self.homeModules.default ];
    home-manager.users.alice = {
      home.stateVersion = stateVersion;
      home.enableNixpkgsReleaseCheck = false;

      programs.todoist-cli.enable = true;
    };

    system.stateVersion = stateVersion;
  };

  testScript = ''
    machine.wait_for_unit("default.target")
    machine.succeed("su -- alice -c 'which todoist'")
  '';
}
